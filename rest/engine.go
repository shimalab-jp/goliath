package rest

import (
    "encoding/json"
    "fmt"
    "github.com/pkg/errors"
    "github.com/shimalab-jp/goliath/database"
    "io/ioutil"
    "net/http"
    "path"
    "strconv"
    "strings"
    "time"

    "github.com/shimalab-jp/goliath/config"
    "github.com/shimalab-jp/goliath/log"
    "github.com/shimalab-jp/goliath/message"
)

type ExecutionHooks interface {
    PreExecute(engine *Engine, request *Request, response *Response) (error)
    PostExecute(engine *Engine, request *Request, response *Response) (error)
}

type Engine struct {
    initialized     bool
    hooks           *ExecutionHooks
    resourceManager *resourceManager
    sessionManager  *sessionManager
}

var instance *Engine = nil

func InitializeEngine() {
    // エンジンを初期化
    if instance == nil {
        // エンジンのインスタンスを作成
        instance = &Engine{
            resourceManager: createResourceManager(),
            sessionManager:  createSessionManager(),
            initialized:     true}
    }
}

func GetEngine() (*Engine) {
    return instance
}

func (e *Engine) getAcceptLanguages(httpRequest *http.Request) ([]message.AcceptLanguage, error) {
    parseTargetString := httpRequest.Header.Get("Accept-Language")
    langList := strings.Split(parseTargetString, ",")

    var returnValue []message.AcceptLanguage
    for _, languageString := range langList {
        languageString = strings.TrimSpace(languageString)

        token := strings.Split(languageString, ";")
        if len(token) == 1 {
            langCode := strings.Split(token[0], "-")
            if len(langCode) > 0 {
                lq := message.AcceptLanguage{Lang: strings.ToLower(langCode[0]), Q: 1}
                returnValue = append(returnValue, lq)
            }
        } else {
            qp := strings.Split(token[1], "=")
            q, err := strconv.ParseFloat(qp[1], 64)
            if err != nil {
                return returnValue, errors.WithStack(err)
            }

            langCode := strings.Split(token[0], "-")
            if len(langCode) > 0 {
                lq := message.AcceptLanguage{Lang: strings.ToLower(langCode[0]), Q: q}
                returnValue = append(returnValue, lq)
            }
        }
    }
    return returnValue, nil
}

func (e *Engine) parseRequest(httpRequest *http.Request) (*Request, error) {
    startTime := time.Now().UnixNano()

    // 戻り値のインスタンスを作成
    returnValue := &Request{
        Languages:    []message.AcceptLanguage{},
        Headers:      map[string]string{},
        GetData:      map[string]string{},
        PostData:     map[string]interface{}{},
        Resource:     nil,
        ResourceInfo: nil,
        Session:      nil}

    // 営業日を取得
    returnValue.BusinessDay = GetBusinessDay(time.Now().Unix())

    // リモートアドレスを取得
    returnValue.RemoteAddress = httpRequest.RemoteAddr

    // ユーザーエージェントを取得
    returnValue.UserAgent = httpRequest.UserAgent()

    // 要求言語を取得
    languages, err := e.getAcceptLanguages(httpRequest)
    if err == nil {
        returnValue.Languages = languages
    }

    // クライアント向けメッセージマネージャのインスタンスを作成
    returnValue.MessageManager = message.CreateMessageManager(&returnValue.Languages)

    // Methodを取得
    returnValue.Method = strings.ToUpper(httpRequest.Method)

    // 出力フォーマット(既定値)
    returnValue.OutputFormat = "json"

    /*
    現在はjsonだけ
    // 出力フォーマットを取得
    if strings.ToLower(path.Ext(httpRequest.URL.Path)) == ".json" {
        returnValue.OutputFormat = "json"
    } else if strings.ToLower(path.Ext(httpRequest.URL.Path)) == ".xml" {
        returnValue.OutputFormat = "xml"
    }
    */

    // リソース名を取得
    removeUrl := strings.TrimRight(config.Values.Server.BaseUrl, "/") + "/"
    apiPath := strings.TrimPrefix(httpRequest.URL.Path, removeUrl)
    returnValue.Name = strings.ToLower("/" + strings.TrimLeft(strings.TrimSuffix(apiPath, path.Ext(httpRequest.URL.Path)), "/"))

    // ヘッダーを取得
    for key := range httpRequest.Header {
        returnValue.Headers[key] = httpRequest.Header.Get(key)
    }

    // QueryStringを取得
    returnValue.QueryString = httpRequest.URL.RawQuery

    // QueryStringをパース
    getData := httpRequest.URL.Query()
    if getData != nil {
        for key, value := range getData {
            returnValue.GetData[key] = value[0]
        }
    }

    // Content-Typeを取得
    returnValue.ContentType = strings.ToLower(httpRequest.Header.Get("Content-Type"))

    // POSTデータをパース
    if strings.ToLower(returnValue.Method) == "post" {
        if returnValue.ContentType == "application/json" {
            // 全データを読込
            buffer, _ := ioutil.ReadAll(httpRequest.Body)

            // jsonを解析
            json.Unmarshal(buffer, &returnValue.PostData)
        } else {
            for key, value := range httpRequest.PostForm {
                returnValue.PostData[key] = value[0]
            }
        }
    }

    // RESTリソースのインスタンスを作成
    returnValue.Resource = e.resourceManager.FindResource(returnValue.Name)
    if returnValue.Resource == nil || (*returnValue.Resource) == nil {
        return nil, nil
    }

    // RESTリソースの定義情報を取得
    returnValue.ResourceInfo = (*returnValue.Resource).Define()

    // 解析時間を記録
    returnValue.ParseTime = time.Now().UnixNano() - startTime

    return returnValue, nil
}

func (e *Engine) checkMethod(request *Request, response *Response) (bool) {
    if _, ok := request.ResourceInfo.Methods[request.Method]; !ok {
        response.StatusCode = http.StatusMethodNotAllowed
        response.ResultCode = ResultNotImplemented
        return false
    }
    return true
}

func (e *Engine) checkSession(request *Request, response *Response) (bool) {
    // ヘッダからセッションIDを取得
    var sessionID string
    if val, ok := request.Headers["X-Goliath-SSID"]; ok {
        sessionID = val
    }

    // アクセスメソッドの定義を取得
    var def *ResourceDefine
    if val, ok := request.ResourceInfo.Methods[request.Method]; ok {
        def = &val
    }

    // アクセスメソッドの定義がされていない場合
    if def == nil {
        response.StatusCode = http.StatusMethodNotAllowed
        response.ResultCode = ResultNotImplemented
        return false
    }

    if len(sessionID) != 36 {
        // 認証を要求されている場合で、セッションIDが不正な場合はエラーとする
        if def.RequireAuthentication {
            response.SetSystemErrorMessage("ERR_ACC_121", []interface{}{}, "SER_RST_221", request.Name, sessionID)
            return false
        }
    } else {
        // 記録されているセッション情報を取得
        request.Session = e.sessionManager.Get(sessionID)

        // 認証を要求されている場合で、セッションを取得できなかった場合はエラーとする
        if request.Session == nil {
            if def.RequireAuthentication {
                response.SetSystemErrorMessage("ERR_ACC_121", []interface{}{}, "SER_RST_222", request.Name, sessionID)
                return false
            }
        }
    }

    return true
}

func (e *Engine) checkBan(request *Request, response *Response) (bool) {
    if request.Session != nil && request.Session.Account.IsBan {
        response.SetSystemErrorMessage("ERR_ACC_102", []interface{}{}, "SER_RST_231", request.Name, request.Session.Account.UserID)
        return false
    }
    return true
}

func (e *Engine) checkApiSwitch(request *Request, response *Response) (bool) {
    con, err := database.Connect("goliath")
    if err != nil {
        response.SetSystemErrorMessage("ERR_CMN_101", []interface{}{}, "SER_RDB_101", err)
        return false
    }

    result, err := con.Query("SELECT `enable` FROM `goliath_mst_api_switch` WHERE `api_name` = ?", request.Name)
    if err != nil {
        response.SetSystemErrorMessage("ERR_CMN_101", []interface{}{}, "SER_RDB_102", err)
        return false
    }

    enable := -1
    for result.Rows.Next() {
        result.Rows.Scan(&enable)
        break
    }

    result.Rows.Close()

    if enable == 1 {
        return true
    } else if enable == 0 {
        response.SetErrorMessage("ERR_CMN_102")
        return false
    } else {
        _, err := con.Execute("INSERT INTO `goliath_mst_api_switch` (`api_name`, `enable`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `enable` = ?;", request.Name, 1, 1)
        if err != nil {
            response.SetSystemErrorMessage("ERR_CMN_101", []interface{}{}, "SER_RDB_103", err)
            return false
        }
        return true
    }
}

func (e *Engine) checkMaintenance(request *Request, response *Response) (bool) {
    // 実行ユーザーが管理者の場合
    if request.Session != nil && request.Session.Account.isAdmin {
        return true
    }

    // メンテナンス中での実行を許可されているAPIの場合
    if v, ok := request.ResourceInfo.Methods[request.Method]; ok {
        if v.RunInMaintenance {
            return true
        }
    }

    con, err := database.Connect("goliath")
    if err != nil {
        response.SetSystemErrorMessage("ERR_CMN_101", []interface{}{}, "SER_RDB_101", err)
        return false
    }

    now := time.Now().Unix()
    result, err := con.Query("SELECT `start_time`, `end_time`, `subject`, `body` FROM `goliath_mst_maintenance` WHERE (`start_time` <= ? AND ? <= `end_time`) OR `start_time` >= ? LIMIT 1;", now, now, now)
    if err != nil {
        response.SetSystemErrorMessage("ERR_CMN_101", []interface{}{}, "SER_RDB_102", err)
        return false
    }

    for result.Rows.Next() {
        var startTime, endTime int64
        var subject, body string
        err := result.Rows.Scan(&startTime, &endTime, &subject, &body)
        if err == nil {
            response.MaintenanceInfo = MaintenanceInfo{
                StartTime: startTime,
                EndTime: endTime,
                Subject: subject,
                Body: body}
        }
        break
    }

    if response.MaintenanceInfo.StartTime <= now && now <= response.MaintenanceInfo.EndTime {
        response.SetErrorMessage("ERR_CMN_103")
        return false
    }
    return true
}

func (e *Engine) checkDebugOnly(request *Request, response *Response) (bool) {
    if v, ok := request.ResourceInfo.Methods[request.Method]; ok {
        if v.IsDebugModeOnly && !config.Values.Server.Debug {
            response.SetErrorMessage("ERR_CMN_102")
            return false
        }
    }
    return true
}

func (e *Engine) checkAdminOnly(request *Request, response *Response) (bool) {
    if v, ok := request.ResourceInfo.Methods[request.Method]; ok {
        if v.IsAdminModeOnly && request.Session != nil && !request.Session.Account.isAdmin {
            response.SetErrorMessage("ERR_CMN_102")
            return false
        }
    }
    return true
}

func (e *Engine) updateLastAccess(request *Request, response *Response) {
    if request.Session == nil {
        return
    }

    var err error = nil
    var con *database.Connection = nil

    tz, _ := time.LoadLocation(config.Values.Server.TimeZone)
    hour := time.Now().In(tz).Hour()

    if err == nil {
        con, err = database.Connect("goliath")
        if err != nil {
            response.SetSystemErrorMessage("ERR_CMN_101", []interface{}{}, "SER_RDB_101", err)
        }
    }

    if err == nil {
        err = con.BeginTransaction()
        if err != nil {
            response.SetSystemErrorMessage("ERR_CMN_101", []interface{}{}, "SER_RDB_104", err)
        }
    }

    if err == nil {
        _, err = con.Execute("REPLACE INTO `goliath_dat_account` (`last_access`) VALUES (?);", time.Now().UnixNano())
        if err != nil {
            response.SetSystemErrorMessage("ERR_CMN_101", []interface{}{}, "SER_RDB_103", err)
        }
    }

    if err == nil {
        _, err = con.Execute(
            "REPLACE INTO `goliath_log_hau` (`access_date`, `access_hour`, `user_id`, `platform`) VALUES (?, ?, ?, ?);",
            request.BusinessDay.BusinessDay,
            hour,
            request.Session.Account.UserID,
            request.Session.Account.Platform)
        if err != nil {
            response.SetSystemErrorMessage("ERR_CMN_101", []interface{}{}, "SER_RDB_103", err)
        }
    }

    if err == nil {
        err = con.Commit()
        if err != nil {
            response.SetSystemErrorMessage("ERR_CMN_101", []interface{}{}, "SER_RDB_105", err)
        }
    } else {
        err = con.Rollback()
        if err != nil {
            response.SetSystemErrorMessage("ERR_CMN_101", []interface{}{}, "SER_RDB_106", err)
        }
    }

    if con != nil {
        con.Disconnect()
    }
}

func (e *Engine) getStatusMessage(statusCode int) (string) {
    switch statusCode {
    case http.StatusOK:
        return "OK"
        break
    case http.StatusBadRequest:
        return "Bad Request"
        break
    case http.StatusUnauthorized:
        return "Unauthorized"
        break
    case http.StatusPaymentRequired:
        return "Payment Required"
        break
    case http.StatusForbidden:
        return "Forbidden"
        break
    case http.StatusNotFound:
        return "Not Found"
        break
    case http.StatusConflict:
        return "Conflict"
        break
    case http.StatusInternalServerError:
        return "Internal Server Error"
        break
    case http.StatusNotImplemented:
        return "Not Implemented"
        break
    case http.StatusServiceUnavailable:
        return "Service Unavailable"
        break
    }
    return "Not Implemented"
}

func (e *Engine) AppendResource(resource *IRestResource) (error) {
    // 初期化チェック
    if !e.initialized {
        return errors.New("not initialized.")
    }

    return e.resourceManager.Append(resource)
}

func (e *Engine) SetHooks(hooks *ExecutionHooks) {
    if e.initialized {
        e.hooks = hooks
    }
}

func (e *Engine) Execute(httpRequest *http.Request, writer http.ResponseWriter) (error) {
    startTime := time.Now().UnixNano()

    // 初期化チェック
    if !e.initialized {
        return errors.New("not initialized.")
    }

    var err error
    var request *Request
    var response *Response

    // レスポンスを作成
    response = &Response{
        Method:          strings.ToUpper(httpRequest.Method),
        StatusCode:      http.StatusOK,
        ResultCode:      ResultOK,
        MaintenanceInfo: MaintenanceInfo{},
        Times:           map[string]int64{},
        MessageManager:  nil}

    // リクエストデータをパース
    if response.ResultCode == ResultOK {
        request, err = e.parseRequest(httpRequest)
        if err != nil || request == nil || request.Resource == nil || (*request.Resource) == nil {
            response.StatusCode = http.StatusNotFound
            response.ResultCode = ResultNotFound
        } else {
            response.MessageManager = request.MessageManager
            response.Times["T120_PARSE_REQUEST"] = request.ParseTime
        }
    }

    // RequestからResultへコピー可能な値をコピーする
    if response.ResultCode == ResultOK {
        st := time.Now().UnixNano()
        response.Name = request.Name
        response.RemoteAddress = request.RemoteAddress
        response.OutputFormat = request.OutputFormat
        response.Times["T201_COPY_REQUEST_TO_RESPONSE"] = time.Now().UnixNano() - st
    }

    // メソッドサポートチェック
    if response.ResultCode == ResultOK {
        st := time.Now().UnixNano()
        e.checkMethod(request, response)
        response.Times["T211_CHECK_METHOD"] = time.Now().UnixNano() - st
    }

    // TODO: クライアントバージョンチェック
    if response.ResultCode == ResultOK {
        st := time.Now().UnixNano()
        response.Times["T212_CLIENT_VERSION"] = time.Now().UnixNano() - st
    }

    // 認証
    if response.ResultCode == ResultOK {
        st := time.Now().UnixNano()
        e.checkSession(request, response)
        response.Times["T213_AUTHENTICATION"] = time.Now().UnixNano() - st
    }

    // 垢バンチェック
    if response.ResultCode == ResultOK {
        st := time.Now().UnixNano()
        e.checkBan(request, response)
        response.Times["T214_CHECK_BAN"] = time.Now().UnixNano() - st
    }

    // Switchチェック
    if response.ResultCode == ResultOK {
        st := time.Now().UnixNano()
        e.checkApiSwitch(request, response)
        response.Times["T220_CHECK_SWITCH"] = time.Now().UnixNano() - st
    }

    // メンテ中チェック
    if response.ResultCode == ResultOK {
        st := time.Now().UnixNano()
        e.checkMaintenance(request, response)
        response.Times["T230_CHECK_MAINTENANCE"] = time.Now().UnixNano() - st
    }

    // デバッグ環境専用チェック
    if response.ResultCode == ResultOK {
        st := time.Now().UnixNano()
        e.checkDebugOnly(request, response)
        response.Times["T241_CHECK_DEBUG_ONLY"] = time.Now().UnixNano() - st
    }

    // 管理者専用チェック
    if response.ResultCode == ResultOK {
        st := time.Now().UnixNano()
        e.checkAdminOnly(request, response)
        response.Times["T242_CHECK_ADMIN_ONLY"] = time.Now().UnixNano() - st
    }

    // TODO: URLパラメータチェック
    if response.ResultCode == ResultOK {
        st := time.Now().UnixNano()
        response.Times["T251_CHECK_URL_PARAMETERS"] = time.Now().UnixNano() - st
    }

    // TODO: POSTパラメータチェック
    if response.ResultCode == ResultOK {
        st := time.Now().UnixNano()
        response.Times["T252_CHECK_POST_PARAMETERS"] = time.Now().UnixNano() - st
    }

    // 最終アクセス日時更新
    if response.ResultCode == ResultOK {
        if request.Session != nil {
            st := time.Now().UnixNano()
            e.updateLastAccess(request, response)
            response.Times["T280_UPDATE_LAST_ACCESS_TIME"] = time.Now().UnixNano() - st
        }
    }

    // 実行
    if response.ResultCode == ResultOK {
        // ユーザー定義の実行前処理を実行
        if e.hooks != nil && (*e.hooks) != nil {
            st := time.Now().UnixNano()
            err = (*e.hooks).PreExecute(e, request, response)
            if err != nil {
                log.Ee(err)
            }
            response.Times["T500_PRE_EXEC"] = time.Now().UnixNano() - st
        }

        // API実行
        st := time.Now().UnixNano()
        switch request.Method {
        case "GET":
            err = (*request.Resource).Get(request, response)
            if err != nil {
                log.Ee(err)
            }
            break
        case "POST":
            err = (*request.Resource).Post(request, response)
            if err != nil {
                log.Ee(err)
            }
            break
        case "PUT":
            err = (*request.Resource).Put(request, response)
            if err != nil {
                log.Ee(err)
            }
            break
        case "DELETE":
            err = (*request.Resource).Delete(request, response)
            if err != nil {
                log.Ee(err)
            }
            break
        default:
            return errors.New(fmt.Sprintf("%s is unsupported method", request.Method))
        }
        response.Times["T510_EXEC"] = time.Now().UnixNano() - st

        // ユーザー定義の実行後処理を実行
        if e.hooks != nil {
            st := time.Now().UnixNano()
            err = (*e.hooks).PostExecute(e, request, response)
            if err != nil {
                log.Ee(err)
            }
            response.Times["T520_POST_EXEC"] = time.Now().UnixNano() - st
        }
    }

    // 処理時間を記録
    response.ProcessTime = time.Now().UnixNano() - startTime

    if response.StatusCode == http.StatusOK {
        // 出力形式に変換
        outputData := response.CreateOutputData()

        // ResponseDataをjsonに変換
        jsonData := outputData.ToJson()
        if jsonData != nil {
            // データサイズを格納
            response.ContentLength = int64(len(*jsonData))

            // 出力
            writer.Header().Set("Content-Type", response.ContentType)
            writer.Header().Set("Content-Length", strconv.FormatInt(response.ContentLength, 10))
            writer.WriteHeader(response.StatusCode)
            writer.Write(*jsonData)
        } else {
            // jsonへの変換に失敗した場合は 500 Internal Server Error
            writer.WriteHeader(http.StatusInternalServerError)
        }

        log.D("%s %s, ProcessTime: %dms", request.Method, request.Name, response.ProcessTime/1000/1000)
    } else {
        writer.WriteHeader(response.StatusCode)
    }

    return err
}

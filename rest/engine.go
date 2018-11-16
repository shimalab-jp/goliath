package rest

import (
    "encoding/json"
    "fmt"
    "github.com/pkg/errors"
    "github.com/shimalab-jp/goliath/config"
    "github.com/shimalab-jp/goliath/database"
    "github.com/shimalab-jp/goliath/log"
    "github.com/shimalab-jp/goliath/message"
    "io/ioutil"
    "net/http"
    "path"
    "reflect"
    "regexp"
    "strconv"
    "strings"
    "time"
)

type ExecutionHooks interface {
    PreExecute(engine *Engine, request *Request, response *Response) (error)
    PostExecute(engine *Engine, request *Request, response *Response) (error)
}

type Engine struct {
    initialized      bool
    hooks            *ExecutionHooks
    resourceManager  *resourceManager
    userAgentPattern *regexp.Regexp
}

var instance *Engine = nil

func InitializeEngine() {
    // エンジンを初期化
    if instance == nil {
        // エンジンのインスタンスを作成
        instance = &Engine{
            resourceManager:  createResourceManager(),
            userAgentPattern: regexp.MustCompile(config.Values.Client.UserAgentPattern),
            initialized:      true}
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
        GetData:      map[string]interface{}{},
        PostData:     map[string]interface{}{},
        Resource:     nil,
        MethodInfo:   nil,
        Account:      nil}

    // バージョン判定
    for _, v := range config.Values.Server.Versions {
        leftVal := strings.ToLower(strings.TrimRight(v.Url, "/") + "/")
        if len(leftVal) > len(httpRequest.URL.Path) {
            continue
        }

        checkVal := strings.ToLower(httpRequest.URL.Path[0:len(leftVal)])
        if checkVal == leftVal {
            returnValue.Version = v.Version
            returnValue.BaseUrl = v.Url
            break
        }
    }

    // 営業日を取得
    returnValue.BusinessDay = GetBusinessDay(time.Now().Unix())

    // リモートアドレスを取得
    returnValue.RemoteAddress = httpRequest.RemoteAddr

    // ヘッダーを取得
    for key := range httpRequest.Header {
        returnValue.Headers[key] = httpRequest.Header.Get(key)
    }

    // ユーザーエージェントを取得
    returnValue.UserAgent = httpRequest.UserAgent()

    // リファレンス用のユーザーエージェントを取得
    if ua, ok := returnValue.Headers["X-Goliath-User-Agent"]; ok {
        returnValue.UserAgent = ua
    }

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
    removeUrl := strings.TrimRight(returnValue.BaseUrl, "/") + "/"
    apiPath := strings.TrimPrefix(httpRequest.URL.Path, removeUrl)
    returnValue.Name = strings.ToLower("/" + strings.TrimLeft(strings.TrimSuffix(apiPath, path.Ext(httpRequest.URL.Path)), "/"))


    // QueryStringを取得
    returnValue.QueryString = httpRequest.URL.RawQuery

    // QueryStringをパース
    getData := httpRequest.URL.Query()
    if getData != nil {
        for key, value := range getData {
            returnValue.GetData[key] = value[0]
        }
    }

    // POSTデータをパース
    if strings.ToLower(returnValue.Method) == "post" {
        if returnValue.OutputFormat == "json" {
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
    returnValue.Resource = e.resourceManager.FindResource(returnValue.Version, returnValue.Name)
    if returnValue.Resource == nil || (*returnValue.Resource) == nil {
        return nil, nil
    }

    // RESTリソースの定義情報を取得
    returnValue.ResourceDefine = (*returnValue.Resource).Define()

    // 解析時間を記録
    returnValue.ParseTime = time.Now().UnixNano() - startTime

    return returnValue, nil
}

func (e *Engine) parseUserAgent(userAgent string) (*ClientVersion) {
    result := e.userAgentPattern.FindAllStringSubmatch(userAgent, -1)
    if len(result) > 0 && len(result[0]) == 6 {
        major, err := strconv.ParseUint(result[0][1], 10, 32)
        if err != nil {
            return nil
        }

        minor, err := strconv.ParseUint(result[0][2], 10, 32)
        if err != nil {
            return nil
        }

        revision, err := strconv.ParseUint(result[0][3], 10, 32)
        if err != nil {
            return nil
        }

        platform := result[0][4]
        environment := result[0][5]

        return &ClientVersion{
            major:       uint32(major),
            minor:       uint32(minor),
            revision:    uint32(revision),
            platform:    platform,
            environment: environment}
    }

    return nil
}

func (e *Engine) checkClientVersion(request *Request, response *Response) (bool) {
    // ユーザーエージェントからバージョン情報を取得
    cv := e.parseUserAgent(request.UserAgent)
    if cv == nil && config.Values.Client.MismatchAction != 0 {
        // MismatchActionが0以外の場合、バージョン情報を取得できなかったら不正クライアントとしてエラーとする
        response.SetErrorMessage("ERR_RES_301")
        return false
    }

    // ユーザーエージェントからバージョン情報を取得できなかった場合は暫定値を設定
    if cv == nil {
        cv = &ClientVersion{major: 1, minor: 0, revision: 0, resourceVersion: "latest", platform: "PC"}
    }

    if config.Values.Client.MismatchAction != 0 {
        // アップデート要求バージョンを取得
        // アップデートが要求されていない場合はnilが返ります
        rv, err := cv.GetUpdateRequireVersion()
        if err != nil {
            // 稼働中のバージョン番号一覧の取得に失敗した場合
            response.SetSystemErrorMessage("ERR_RES_100", []interface{}{}, "SER_RST_201", err)
        }

        if rv != nil {
            // アップデート要求バージョンがあるので、アップデート要求のステータスコード(600)を返す
            response.SetErrorMessage("ERR_RES_302", rv.GetVersion())

            // アップデート用のストアURLを設定
            switch cv.GetPlatform() {
            case PlatformApple:
                response.LinkUrl = config.Values.StoreUrl.Apple
                break
            case PlatformGoogle:
                response.LinkUrl = config.Values.StoreUrl.Google
                break
            }
            return false
        }
    }

    return true
}

func (e *Engine) checkMethod(request *Request, response *Response) (bool) {
    if v, ok := request.ResourceDefine.Methods[request.Method]; !ok {
        response.StatusCode = http.StatusMethodNotAllowed
        response.ResultCode = ResultNotImplemented
        return false
    } else {
        request.MethodInfo = &v
        return true
    }
}

func (e *Engine) checkToken(request *Request, response *Response) (bool) {
    // ヘッダからトークンを取得
    var token string
    if val, ok := request.Headers["X-Goliath-Token"]; ok {
        token = strings.ToLower(val)
    }

    if len(token) != 128 {
        // 認証を要求されている場合で、トークンが不正な場合はエラーとする
        if request.MethodInfo.RequireAuthentication {
            response.SetSystemErrorMessage("ERR_RES_122", []interface{}{}, "SER_RST_221", request.Name, token)
            return false
        }
    } else {
        // アカウントを取得
        account, err := GetAccountManager().GetAccountByToken(token)
        if err != nil {
            log.Ee(err)
            response.SetSystemErrorMessage("ERR_RES_100", []interface{}{}, "SER_RST_222", request.Name, token, err)
            return false
        }

        // アカウント情報を格納
        request.Account = account

        // 認証を要求されている場合で、アカウントを取得できなかった場合はエラーとする
        if request.Account == nil {
            if request.MethodInfo.RequireAuthentication {
                response.SetSystemErrorMessage("ERR_RES_122", []interface{}{}, "SER_RST_221", request.Name, token)
                return false
            }
        } else {
            request.Token = account.Token
        }
    }

    return true
}

func (e *Engine) checkBan(request *Request, response *Response) (bool) {
    if request.Account != nil && request.Account.IsBan {
        response.SetErrorMessage("ERR_RES_121")
        return false
    }
    return true
}

func (e *Engine) checkApiSwitch(request *Request, response *Response) (bool) {
    con, err := database.Connect("goliath")
    if err != nil {
        response.SetSystemErrorMessage("ERR_RES_100", []interface{}{}, "SER_RDB_101", err)
        return false
    }

    result, err := con.Query("SELECT `enable` FROM `goliath_mst_api_switch` WHERE `api_name` = ?", (*request.Resource).GetPath())
    if err != nil {
        response.SetSystemErrorMessage("ERR_RES_100", []interface{}{}, "SER_RDB_102", err)
        return false
    }

    enable := -1
    if result.MoveFirst() {
        enable = result.GetInt("enable", enable)
    }

    if enable == 1 {
        return true
    } else if enable == 0 {
        response.SetErrorMessage("ERR_RES_111")
        return false
    } else {
        _, err := con.Execute("INSERT INTO `goliath_mst_api_switch` (`api_name`, `enable`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `enable` = ?;", (*request.Resource).GetPath(), 1, 1)
        if err != nil {
            response.SetSystemErrorMessage("ERR_RES_100", []interface{}{}, "SER_RDB_103", err)
            return false
        }
        return true
    }
}

func (e *Engine) checkMaintenance(request *Request, response *Response) (bool) {
    // 実行ユーザーが管理者の場合
    if request.Account != nil && request.Account.IsAdmin {
        return true
    }

    // メンテナンス中での実行を許可されているAPIの場合
    if v, ok := request.ResourceDefine.Methods[request.Method]; ok {
        if v.RunInMaintenance {
            return true
        }
    }

    con, err := database.Connect("goliath")
    if err != nil {
        response.SetSystemErrorMessage("ERR_RES_100", []interface{}{}, "SER_RDB_101", err)
        return false
    }

    now := time.Now().Unix()
    result, err := con.Query("SELECT `start_time`, `end_time`, `subject`, `body` FROM `goliath_mst_maintenance` WHERE (`start_time` <= ? AND ? <= `end_time`) OR `start_time` >= ? LIMIT 1;", now, now, now)
    if err != nil {
        response.SetSystemErrorMessage("ERR_RES_100", []interface{}{}, "SER_RDB_102", err)
        return false
    }

    if result.MoveFirst() {
        startTime := result.GetInt64("start_time", 0)
        endTime := result.GetInt64("end_time", 0)
        subject := result.GetString("subject", "")
        body := result.GetString("body", "")

        response.MaintenanceInfo = MaintenanceInfo{
            StartTime: startTime,
            EndTime:   endTime,
            Subject:   subject,
            Body:      body}
    }

    if response.MaintenanceInfo.StartTime <= now && now <= response.MaintenanceInfo.EndTime {
        response.SetErrorMessage("ERR_RES_112")
        return false
    }
    return true
}

func (e *Engine) checkDebugOnly(request *Request, response *Response) (bool) {
    if request.MethodInfo.IsDebugModeOnly && !config.Values.Server.Debug.Enable {
        response.SetErrorMessage("ERR_RES_111")
        return false
    }
    return true
}

func (e *Engine) checkAdminOnly(request *Request, response *Response) (bool) {
    if request.MethodInfo.IsAdminModeOnly && request.Account != nil && !request.Account.IsAdmin {
        response.SetErrorMessage("ERR_RES_111")
        return false
    }
    return true
}

func (e *Engine) checkPostParameters(request *Request, response *Response) (bool) {
    if len(request.MethodInfo.PostParameters) <= 0 {
        return true
    }

    var result = true
    for name, def := range request.MethodInfo.PostParameters {
        // 無視パラメータチェック
        if name == "RequestID" {
            continue
        }

        // 必須パラメータチェック
        if result {
            if def.Require {
                if _, ok := request.PostData[name]; !ok {
                    response.SetErrorMessage("ERR_RES_101", name)
                    result = false
                }
            } else {
                if _, ok := request.PostData[name]; !ok {
                    continue
                }
            }
        }

        // 値を文字列に変換
        strVal := fmt.Sprintf("%+v", request.PostData[name])

        // 型チェック(全て)
        if result {
            switch def.Type {
            case reflect.Bool:
                if _, err := strconv.ParseBool(strVal); err != nil {
                    result = false
                    response.SetErrorMessage("ERR_RES_102", name, "Bool")
                }
                break
            case reflect.Int8:
                if _, err := strconv.ParseInt(strVal, 10, 8); err != nil {
                    result = false
                    response.SetErrorMessage("ERR_RES_102", name, "Int8")
                }
                break
            case reflect.Int16:
                if _, err := strconv.ParseInt(strVal, 10, 16); err != nil {
                    result = false
                    response.SetErrorMessage("ERR_RES_102", name, "Int16")
                }
                break
            case reflect.Int32:
                if _, err := strconv.ParseInt(strVal, 10, 32); err != nil {
                    result = false
                    response.SetErrorMessage("ERR_RES_102", name, "Int32")
                }
                break
            case reflect.Int64:
            case reflect.Int:
                if _, err := strconv.ParseInt(strVal, 10, 64); err != nil {
                    result = false
                    response.SetErrorMessage("ERR_RES_102", name, "Int64")
                }
                break
            case reflect.Uint8:
                if _, err := strconv.ParseUint(strVal, 10, 8); err != nil {
                    result = false
                    response.SetErrorMessage("ERR_RES_102", name, "Uint8")
                }
                break
            case reflect.Uint16:
                if _, err := strconv.ParseUint(strVal, 10, 16); err != nil {
                    result = false
                    response.SetErrorMessage("ERR_RES_102", name, "Uint16")
                }
                break
            case reflect.Uint32:
                if _, err := strconv.ParseUint(strVal, 10, 32); err != nil {
                    result = false
                    response.SetErrorMessage("ERR_RES_102", name, "Uint32")
                }
                break
            case reflect.Uint64:
            case reflect.Uint:
                if _, err := strconv.ParseUint(strVal, 10, 64); err != nil {
                    result = false
                    response.SetErrorMessage("ERR_RES_102", name, "Uint64")
                }
                break
            case reflect.Float32:
                if _, err := strconv.ParseFloat(strVal, 32); err != nil {
                    result = false
                    response.SetErrorMessage("ERR_RES_102", name, "Float32")
                }
                break
            case reflect.Float64:
                if _, err := strconv.ParseFloat(strVal, 64); err != nil {
                    result = false
                    response.SetErrorMessage("ERR_RES_102", name, "Float64")
                }
                break
            default:
                break
            }
        }

        // 範囲チェック(数値のみ)
        if result && len(def.Range) >= 2 {
            switch def.Type {
            case reflect.Int8:
            case reflect.Int16:
            case reflect.Int32:
            case reflect.Int64:
            case reflect.Int:
                if v, err := strconv.ParseInt(strVal, 10, 64); err == nil {
                    if !(int64(def.Range[0]) <= v && v <= int64(def.Range[1])) {
                        result = false
                        response.SetErrorMessage("ERR_RES_103", name, def.Range[0], def.Range[1])
                    }
                }
                break
            case reflect.Uint8:
            case reflect.Uint16:
            case reflect.Uint32:
            case reflect.Uint64:
            case reflect.Uint:
                if v, err := strconv.ParseUint(strVal, 10, 64); err == nil {
                    if !(uint64(def.Range[0]) <= v && v <= uint64(def.Range[1])) {
                        result = false
                        response.SetErrorMessage("ERR_RES_103", name, def.Range[0], def.Range[1])
                    }
                }
                break
            case reflect.Float32:
            case reflect.Float64:
                if v, err := strconv.ParseFloat(strVal, 64); err == nil {
                    if !(def.Range[0] <= v && v <= def.Range[1]) {
                        result = false
                        response.SetErrorMessage("ERR_RES_103", name, def.Range[0], def.Range[1])
                    }
                }
                break
            }
        }

        // 選択チェック(数値、文字列)
        if result && len(def.Select) >= 1 {
            switch def.Type {
            case reflect.Int8:
            case reflect.Int16:
            case reflect.Int32:
            case reflect.Int64:
            case reflect.Int:
            case reflect.Uint8:
            case reflect.Uint16:
            case reflect.Uint32:
            case reflect.Uint64:
            case reflect.Uint:
            case reflect.Float32:
            case reflect.Float64:
            case reflect.String:
                hit := false
                search := ""
                for _, val := range def.Select {
                    selVal := fmt.Sprintf("%+v", val)
                    if len(search) > 0 {
                        search = fmt.Sprintf("%s, %s", search, selVal)
                    } else {
                        search = selVal
                    }
                    if strVal == selVal {
                        hit = true
                        break
                    }
                }
                if !hit {
                    result = false
                    response.SetErrorMessage("ERR_RES_104", name, search)
                }
                break
            }
        }

        // 正規表現チェック(文字列のみ)
        if result && len(def.Regex) > 0 && def.Type == reflect.String {
            r := regexp.MustCompile(def.Regex)
            result = r.MatchString(strVal)
            if !result {
                response.SetErrorMessage("ERR_RES_105", name)
            }
        }
    }
    return result
}

func (e *Engine) updateLastAccess(request *Request, response *Response) {
    if request.Account == nil {
        return
    }

    var err error = nil
    var con *database.Connection = nil

    tz, _ := time.LoadLocation(config.Values.Server.TimeZone)
    hour := time.Now().In(tz).Hour()

    if err == nil {
        con, err = database.Connect("goliath")
        if err != nil {
            response.SetSystemErrorMessage("ERR_RES_100", []interface{}{}, "SER_RDB_101", err)
        }
    }

    if err == nil {
        err = con.BeginTransaction()
        if err != nil {
            response.SetSystemErrorMessage("ERR_RES_100", []interface{}{}, "SER_RDB_104", err)
        }
    }

    if err == nil {
        _, err = con.Execute("REPLACE INTO `goliath_dat_account` (`last_access`) VALUES (?);", time.Now().UnixNano())
        if err != nil {
            response.SetSystemErrorMessage("ERR_RES_100", []interface{}{}, "SER_RDB_103", err)
        }
    }

    if err == nil {
        _, err = con.Execute(
            "REPLACE INTO `goliath_log_hau` (`access_date`, `access_hour`, `user_id`, `platform`) VALUES (?, ?, ?, ?);",
            request.BusinessDay.BusinessDay,
            hour,
            request.Account.UserID,
            request.Account.Platform)
        if err != nil {
            response.SetSystemErrorMessage("ERR_RES_100", []interface{}{}, "SER_RDB_103", err)
        }
    }

    if err == nil {
        err = con.Commit()
        if err != nil {
            response.SetSystemErrorMessage("ERR_RES_100", []interface{}{}, "SER_RDB_105", err)
        }
    } else {
        err = con.Rollback()
        if err != nil {
            response.SetSystemErrorMessage("ERR_RES_100", []interface{}{}, "SER_RDB_106", err)
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

func (e *Engine) AppendResource(version uint32, path string, resource *IRestResource) (error) {
    // 初期化チェック
    if !e.initialized {
        return errors.New("not initialized.")
    }

    return e.resourceManager.Append(version, path, resource)
}

func (e *Engine) SetHooks(hooks *ExecutionHooks) {
    if e.initialized {
        e.hooks = hooks
    }
}

func (e *Engine) GetResourceManager() (*resourceManager) {
    return e.resourceManager
}

func (e *Engine) Execute(httpRequest *http.Request, writer http.ResponseWriter) (error) {
    // 処理開始時間を記録
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
            response.Times["T120_PARSE_REQUEST"] = request.ParseTime
        }
    }

    // RequestからResultへコピー可能な値をコピーする
    if response.ResultCode == ResultOK {
        st := time.Now().UnixNano()
        response.Name = request.Name
        response.ContentType = request.ContentType
        response.MessageManager = request.MessageManager
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

    // クライアントバージョンチェック
    if response.ResultCode == ResultOK {
        st := time.Now().UnixNano()
        e.checkClientVersion(request, response)
        response.Times["T212_CLIENT_VERSION"] = time.Now().UnixNano() - st
    }

    // 認証
    if response.ResultCode == ResultOK {
        st := time.Now().UnixNano()
        e.checkToken(request, response)
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
        //e.checkParameters(request, response, &request.MethodInfo.UrlParameters, &request.GetData)
        response.Times["T251_CHECK_URL_PARAMETERS"] = time.Now().UnixNano() - st
    }

    // POSTパラメータチェック
    if response.ResultCode == ResultOK {
        st := time.Now().UnixNano()
        e.checkPostParameters(request, response)
        response.Times["T252_CHECK_POST_PARAMETERS"] = time.Now().UnixNano() - st
    }

    // 最終アクセス日時更新
    if response.ResultCode == ResultOK {
        if request.Account != nil {
            st := time.Now().UnixNano()
            e.updateLastAccess(request, response)
            response.Times["T280_UPDATE_LAST_ACCESS_TIME"] = time.Now().UnixNano() - st
        }
    }

    // 実行
    if err == nil && response.ResultCode == ResultOK {
        // ユーザー定義の実行前処理を実行
        if err == nil && e.hooks != nil && (*e.hooks) != nil {
            st := time.Now().UnixNano()
            err = (*e.hooks).PreExecute(e, request, response)
            response.Times["T500_PRE_EXEC"] = time.Now().UnixNano() - st
        }

        // API実行
        if err == nil {
            st := time.Now().UnixNano()
            switch request.Method {
            case "GET":
                err = (*request.Resource).Get(request, response)
                break
            case "POST":
                err = (*request.Resource).Post(request, response)
                break
            case "PUT":
                err = (*request.Resource).Put(request, response)
                break
            case "DELETE":
                err = (*request.Resource).Delete(request, response)
                break
            default:
                response.StatusCode = http.StatusMethodNotAllowed
                response.ResultCode = ResultNotImplemented
            }
            response.Times["T510_EXEC"] = time.Now().UnixNano() - st
        }

        // ユーザー定義の実行後処理を実行
        if err == nil && e.hooks != nil {
            st := time.Now().UnixNano()
            err = (*e.hooks).PostExecute(e, request, response)
            response.Times["T520_POST_EXEC"] = time.Now().UnixNano() - st
        }

        // エラーの場合
        if err != nil {
            response.SetSystemErrorMessage("ERR_RES_100", []interface{}{}, "SER_RST_231", err)
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
            if response.OutputFormat == "json" {
                writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
            }
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

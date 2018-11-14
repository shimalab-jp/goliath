package goliath

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "mime"
    "net/http"
    "os"
    "path"
    "strings"

    "github.com/shimalab-jp/goliath/config"
    "github.com/shimalab-jp/goliath/log"
    "github.com/shimalab-jp/goliath/rest"
    "github.com/shimalab-jp/goliath/util"
)

const (
    DataVersion    string = "1.0.0"
    DataVersionNum uint32 = 1
)

func requestHandler(w http.ResponseWriter, r *http.Request) {
    err := rest.GetEngine().Execute(r, w)
    if err != nil {
        log.Ee(err)
    }
}

func referenceHandler(w http.ResponseWriter, r *http.Request) {
    accessPath := r.URL.Path[len(config.Values.Server.Reference.Url):]
    if len(accessPath) <= 0 {
        accessPath = "index.html"
    }
    vPath := fmt.Sprintf("${GOPATH}/src/github.com/shimalab-jp/goliath/reference/%s", accessPath)
    realPath := os.ExpandEnv(vPath)

    if accessPath == "config.json" {
        // リファレンス情報を取得
        envName := config.Values.Server.Reference.Environment
        envCode := util.ToEnvironmentCode(envName)
        className := ""
        switch envCode {
        case util.EnvironmentDemo:
            envName = "DEMO"
            className = "env_name_demo"
            break
        case util.EnvironmentDevelop1:
            envName = "DEVELOP1"
            className = "env_name_dev1"
            break
        case util.EnvironmentDevelop2:
            envName = "DEVELOP2"
            className = "env_name_dev2"
            break
        case util.EnvironmentTest:
            envName = "TEST"
            className = "env_name_test"
            break
        case util.EnvironmentAppleReview:
            envName = "APPLE REVIEW"
            className = "env_name_apl"
            break
        case util.EnvironmentStaging:
            envName = "STAGING"
            className = "env_name_stg"
            break
        case util.EnvironmentProduction:
            envName = "PRODUCTION"
            className = "env_name_prd"
            break
        default:
            envName = "LOCAL"
            className = "env_name_loc"
            break
        }

        // リファレンスデータを作成
        ret := map[string]map[string]rest.ResourceDefine{}

        for i, v := range *rest.GetEngine().GetResourceManager().GetAllResources() {
            if v != nil && *v != nil {
                dir, _ := path.Split((*v).Define().Path)

                if group, ok := ret[dir]; ok {
                    group[i] = *(*v).Define()
                } else {
                    ret[dir] = map[string]rest.ResourceDefine{i:*(*v).Define()}
                }
            }
        }

        // jsonデータを作成
        data := map[string]interface{}{
            "Name": config.Values.Server.Reference.Name,
            "EnvName": config.Values.Server.Reference.Environment,
            "EnvCode": envCode,
            "EnvClass": className,
            "Logo": config.Values.Server.Reference.Logo,
            "UserAgent": config.Values.Server.Reference.UserAgent,
            "Versions": config.Values.Server.Versions,
            "Resources": ret}

        // jsonにエンコード
        buffer, err := json.Marshal(data)

        if err != nil {
            w.WriteHeader(http.StatusForbidden)
        } else {
            w.Header().Set("Content-Type", "application/json")
            w.Write(buffer)
        }
    } else if util.FileExists(realPath) {
        buffer, err := ioutil.ReadFile(realPath)
        if err != nil {
            w.WriteHeader(http.StatusForbidden)
        } else {
            mimeType := mime.TypeByExtension(path.Ext(realPath))
            w.Header().Set("Content-Type", mimeType)
            w.Write(buffer)
        }
    } else {
        w.WriteHeader(http.StatusNotFound)
    }
}

func AppendResource(resource rest.IRestResource) (error) {
    return rest.GetEngine().AppendResource(&resource)
}

func SetHooks(hooks rest.ExecutionHooks) {
    rest.GetEngine().SetHooks(&hooks)
}

func Listen() (error) {
    // API用のハンドラを追加
    for _, v := range config.Values.Server.Versions {
        apiUrl := strings.TrimRight(v.Url, "/") + "/"
        http.HandleFunc(apiUrl, requestHandler)
    }

    // リファレンス用
    if config.Values.Server.Reference.Enable {
        // リファレンス用のハンドラを追加
        referenceUrl := strings.TrimRight(config.Values.Server.Reference.Url, "/") + "/"
        http.HandleFunc(referenceUrl, referenceHandler)
    }

    // httpサーバーを起動
    return http.ListenAndServe(fmt.Sprintf(":%d", config.Values.Server.Port), nil)
}

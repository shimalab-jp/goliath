package goliath

import (
    "fmt"
    "github.com/shimalab-jp/goliath/rest/resources"
    "net/http"
    "strings"

    "github.com/shimalab-jp/goliath/config"
    "github.com/shimalab-jp/goliath/log"
    "github.com/shimalab-jp/goliath/rest"
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
    fmt.Fprintf(w, "hello")
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
        // リファレンス用のリソースを追加
        ref := resources.Reference{}
        ref.Resources = rest.GetEngine().GetResourceManager().GetAllResources()
        AppendResource(&ref)

        // リファレンス用のハンドラを追加
        referenceUrl := strings.TrimRight(config.Values.Server.Reference.Url, "/") + "/"
        http.HandleFunc(referenceUrl, referenceHandler)
    }

    // httpサーバーを起動
    return http.ListenAndServe(fmt.Sprintf(":%d", config.Values.Server.Port), nil)
}

package goliath

import (
    "fmt"
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

func AppendResource(resource rest.IRestResource) (error) {
    return rest.GetEngine().AppendResource(&resource)
}

func SetHooks(hooks rest.ExecutionHooks) {
    rest.GetEngine().SetHooks(&hooks)
}

func Listen() (error) {
    listenUrl := strings.TrimRight(config.Values.Server.BaseUrl, "/") + "/"
    http.HandleFunc(listenUrl, requestHandler)
    return http.ListenAndServe(fmt.Sprintf(":%d", config.Values.Server.Port), nil)
}

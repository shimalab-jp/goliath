package account

import (
    "github.com/shimalab-jp/goliath/rest"
    "reflect"
)

type Password struct {
    rest.ResourceBase
}

func (res Password) Define() (*rest.ResourceDefine) {
    return &rest.ResourceDefine{
        Path:    "/account/password",
        Methods: map[string]rest.ResourceMethodDefine{
            "POST": {
                Summary:       "パスワードリセット",
                Description:   "データ移行用のパスワードをリセットします。",
                UrlParameters: map[string]rest.Parameter{},
                PostParameters: map[string]rest.Parameter{},
                Returns: map[string]rest.Return{
                    "AccountInfo": {
                        Type:        reflect.Map,
                        Description: "アカウント情報"}},
                RequireAuthentication: false,
                IsDebugModeOnly:       false,
                RunInMaintenance:      false}}}
}

func (res Password) Post(request *rest.Request, response *rest.Response) (error) {
    return nil
}

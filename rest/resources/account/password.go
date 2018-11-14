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
                UrlParameters: []rest.UrlParameter{},
                PostParameters: map[string]rest.PostParameter{},
                Returns: map[string]rest.Return{
                    "AccountInfo": {
                        Type:        reflect.Map,
                        Description: "アカウント情報"}},
                RequireAuthentication: true,
                IsDebugModeOnly:       false,
                RunInMaintenance:      false}}}
}

func (res Password) Post(request *rest.Request, response *rest.Response) (error) {
    return nil
}

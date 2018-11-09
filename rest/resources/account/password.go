package account

import (
    "github.com/shimalab-jp/goliath/rest"
)

type Password struct {
    rest.ResourceBase
}

func (res Password) Define() (*rest.ResourceInfo) {
    return &rest.ResourceInfo{
        Path:    "/account/password",
        Methods: map[string]rest.ResourceDefine{
            "POST": {
                Summary:       "パスワードリセット",
                Description:   "データ移行用のパスワードをリセットします。",
                UrlParameters: map[string]rest.Parameter{},
                PostParameters: map[string]rest.Parameter{},
                Returns: map[string]rest.Return{
                    "account_auth_info": {
                        Type:        "array",
                        Description: "アカウント認証情報"}},
                RequireAuthentication: false,
                IsDebugModeOnly:       false,
                RunInMaintenance:      false}}}
}

func (res Password) Post(request *rest.Request, response *rest.Response) (error) {
    return nil
}

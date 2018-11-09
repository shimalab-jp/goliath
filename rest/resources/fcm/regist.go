package fcm

import (
    "net/http"

    "github.com/shimalab-jp/goliath/rest"
)

type Regist struct {
    rest.ResourceBase
}

func (res Regist) Define() (*rest.ResourceInfo) {
    return &rest.ResourceInfo{
        Path:    "/fcm/regist",
        Methods: map[string]rest.ResourceDefine{
            "POST": {
                Summary:       "FCMトークン登録／更新",
                Description:   "プッシュ通知用のデバイストークンを登録または更新します。",
                UrlParameters: map[string]rest.Parameter{},
                PostParameters: map[string]rest.Parameter{
                    "token": {
                        Type:        "string",
                        Default:     "",
                        Require:     false,
                        Description: "FCMのトークン。空文字を渡すと削除されます。"}},
                Returns: map[string]rest.Return{
                    "account_auth_info": {
                        Type:        "array",
                        Description: "アカウント認証情報"}},
                RequireAuthentication: false,
                IsDebugModeOnly:       false,
                RunInMaintenance:      false}}}
}

func (res Regist) Get(request *rest.Request, response *rest.Response) (error) {
    response.StatusCode = http.StatusMethodNotAllowed
    response.ResultCode = rest.ResultNotImplemented
    return nil
}

func (res Regist) Post(request *rest.Request, response *rest.Response) (error) {
    return nil
}

func (res Regist) Delete(request *rest.Request, response *rest.Response) (error) {
    response.StatusCode = http.StatusMethodNotAllowed
    response.ResultCode = rest.ResultNotImplemented
    return nil
}

func (res Regist) Put(request *rest.Request, response *rest.Response) (error) {
    response.StatusCode = http.StatusMethodNotAllowed
    response.ResultCode = rest.ResultNotImplemented
    return nil
}

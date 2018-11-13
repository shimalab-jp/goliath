package fcm

import (
    "reflect"

    "github.com/shimalab-jp/goliath/rest"
)

type Regist struct {
    rest.ResourceBase
}

func (res Regist) Define() (*rest.ResourceDefine) {
    return &rest.ResourceDefine{
        Path:    "/fcm/regist",
        Methods: map[string]rest.ResourceMethodDefine{
            "POST": {
                Summary:       "FCMトークン登録／更新",
                Description:   "プッシュ通知用のデバイストークンを登録または更新します。",
                UrlParameters: map[string]rest.Parameter{},
                PostParameters: map[string]rest.Parameter{
                    "Token": {
                        Type:        reflect.String,
                        Default:     "",
                        Require:     false,
                        Description: "FCMのトークン。空文字を渡すと削除されます。"}},
                Returns: map[string]rest.Return{
                    "AccountInfo": {
                        Type:        reflect.Map,
                        Description: "アカウント情報"}},
                RequireAuthentication: false,
                IsDebugModeOnly:       false,
                RunInMaintenance:      false}}}
}

func (res Regist) Post(request *rest.Request, response *rest.Response) (error) {
    return nil
}

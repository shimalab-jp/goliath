package v1

import (
    "github.com/shimalab-jp/goliath/rest"
    "reflect"
)

type FcmRegist struct {
    rest.ResourceBase
}

func (res FcmRegist) Define() *rest.ResourceDefine {
    return &rest.ResourceDefine{
        Methods: map[string]rest.ResourceMethodDefine{
            "POST": {
                Summary:         "FCMトークン登録／更新",
                Description:     "プッシュ通知用のデバイストークンを登録または更新します。",
                UrlParameters:   []rest.UrlParameter{},
                QueryParameters: map[string]rest.QueryParameter{},
                PostParameters: map[string]rest.PostParameter{
                    "Token": {
                        Type:        reflect.String,
                        Default:     "",
                        Require:     false,
                        Description: "FCMのトークン。空文字を渡すと削除されます。"}},
                Returns: map[string]rest.Return{
                    "AccountInfo": {
                        Type:        reflect.Map,
                        Description: "アカウント情報"}},
                RequireAuthentication: true,
                IsDebugModeOnly:       false,
                RunInMaintenance:      false}}}
}

func (res FcmRegist) Post(request *rest.Request, response *rest.Response) error {
    return nil
}

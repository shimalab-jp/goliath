package v1

import (
    "github.com/shimalab-jp/goliath/rest"
    "reflect"
    "strconv"
)

type AccountTrans struct {
    rest.ResourceBase
}

func (res AccountTrans) Define() (*rest.ResourceDefine) {
    return &rest.ResourceDefine{
        Methods: map[string]rest.ResourceMethodDefine{
            "POST": {
                Summary:         "アカウント移譲",
                Description:     "新しい端末にプレイでデータを移譲します。",
                UrlParameters:   []rest.UrlParameter{},
                QueryParameters: map[string]rest.QueryParameter{},
                PostParameters: map[string]rest.PostParameter{
                    "PlayerID": {
                        Type:        reflect.String,
                        Default:     rest.PlatformNone,
                        Regex:       "/[0-9]{4,4}-[0-9]{4,4}/",
                        Require:     true,
                        Description: "プレイヤーID"},
                    "Password": {
                        Type:        reflect.String,
                        Default:     rest.PlatformNone,
                        Regex:       "/[0-9A-F]{40,40}/",
                        Require:     true,
                        Description: "パスワード。SHA1でハッシュ化した値を指定してください。"},
                    "Platform": {
                        Type:        reflect.Uint8,
                        Default:     rest.PlatformNone,
                        Select:      []interface{}{rest.PlatformNone, rest.PlatformApple, rest.PlatformGoogle},
                        Require:     true,
                        Description: "プラットフォーム。" + strconv.Itoa(rest.PlatformNone) + ":None, " + strconv.Itoa(rest.PlatformApple) + ":Apple, " + strconv.Itoa(rest.PlatformGoogle) + ":Google"}},
                Returns: map[string]rest.Return{
                    "AccountInfo": {
                        Type:        reflect.Map,
                        Description: "アカウント情報"}},
                RequireAuthentication: false,
                IsDebugModeOnly:       false,
                RunInMaintenance:      false}}}
}

func (res AccountTrans) Post(request *rest.Request, response *rest.Response) (error) {
    return nil
}

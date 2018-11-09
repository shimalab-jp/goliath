package account

import (
    "github.com/shimalab-jp/goliath/rest"
)

type Trans struct {
    rest.ResourceBase
}

func (res Trans) Define() (*rest.ResourceInfo) {
    return &rest.ResourceInfo{
        Path:    "/account/trans",
        Methods: map[string]rest.ResourceDefine{
            "POST": {
                Summary:       "アカウント移譲",
                Description:   "新しい端末にプレイでデータを移譲します。",
                UrlParameters: map[string]rest.Parameter{},
                PostParameters: map[string]rest.Parameter{
                    "player_id": {
                        Type:        "string",
                        Default:     rest.PlatformNone,
                        Regex:       "/[0-9]{3,3}-[0-9]{3,3}-[0-9]{3,3}/",
                        Require:     true,
                        Description: "プレイヤーID"},
                    "password": {
                        Type:        "string",
                        Default:     rest.PlatformNone,
                        Regex:       "/[0-9A-F]{40,40}/",
                        Require:     true,
                        Description: "パスワード。SHA1でハッシュ化した値を指定してください。"},
                    "platform": {
                        Type:        "int",
                        Default:     rest.PlatformNone,
                        Select:      []interface{}{rest.PlatformNone, rest.PlatformApple, rest.PlatformGoogle},
                        Require:     true,
                        Description: "プラットフォーム。" + string(rest.PlatformNone) + ":None, " + string(rest.PlatformApple) + ":Apple, " + string(rest.PlatformGoogle) + ":Google"}},
                Returns: map[string]rest.Return{
                    "account_auth_info": {
                        Type:        "array",
                        Description: "アカウント認証情報"}},
                RequireAuthentication: false,
                IsDebugModeOnly:       false,
                RunInMaintenance:      false}}}
}


func (res Trans) Post(request *rest.Request, response *rest.Response) (error) {
    return nil
}

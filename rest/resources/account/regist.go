package account

import (
    "github.com/shimalab-jp/goliath/rest"
    "reflect"
)

type Regist struct {
    rest.ResourceBase
}

func (res Regist) Define() (*rest.ResourceInfo) {
    return &rest.ResourceInfo{
        Path:    "/account/regist",
        Methods: map[string]rest.ResourceDefine{
            "POST": {
                Summary:       "アカウント登録",
                Description:   "新規アカウントを登録します。",
                UrlParameters: map[string]rest.Parameter{},
                PostParameters: map[string]rest.Parameter{
                    "Platform": {
                        Type:        reflect.Uint8,
                        Default:     rest.PlatformNone,
                        Select:      []interface{}{rest.PlatformNone, rest.PlatformApple, rest.PlatformGoogle},
                        Require:     true,
                        Description: "プラットフォーム。" + string(rest.PlatformNone) + ":None, " + string(rest.PlatformApple) + ":Apple, " + string(rest.PlatformGoogle) + ":Google"}},
                Returns: map[string]rest.Return{
                    "AccountInfo": {
                        Type:        reflect.Map,
                        Description: "アカウント情報"}},
                RequireAuthentication: false,
                IsDebugModeOnly:       false,
                RunInMaintenance:      false}}}
}

func (res Regist) Post(request *rest.Request, response *rest.Response) (error) {
    // パラメータを取得
    platform := request.GetInt8("Platform", rest.PlatformNone)

    // アカウントを作成
    am := rest.GetAccountManager()
    account, err := am.Create(request, platform)
    if err != nil {
        return err
    }

    // 戻り値に値をセット
    response.Result = map[string]interface{}{ "AccountInfo": account }

    return nil
}

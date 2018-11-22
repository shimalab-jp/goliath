package v1

import (
    "github.com/shimalab-jp/goliath/rest"
    "reflect"
    "strconv"
)

type AccountRegist struct {
    rest.ResourceBase
}

func (res AccountRegist) Define() *rest.ResourceDefine {
    return &rest.ResourceDefine{
        Methods: map[string]rest.ResourceMethodDefine{
            "POST": {
                Summary:         "アカウント登録",
                Description:     "新規アカウントを登録します。",
                UrlParameters:   []rest.UrlParameter{},
                QueryParameters: map[string]rest.QueryParameter{},
                PostParameters: map[string]rest.PostParameter{
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

func (res AccountRegist) Post(request *rest.Request, response *rest.Response) error {
    // パラメータを取得
    platform, _ := request.GetParamInt8(rest.PostParam, "Platform", rest.PlatformNone)

    // アカウントを作成
    am := rest.GetAccountManager()
    account, err := am.Create(request, platform)
    if err != nil {
        return err
    }

    // 戻り値に値をセット
    response.Result = map[string]interface{}{"AccountInfo": account.Output()}

    return nil
}

package account

import (
    "github.com/shimalab-jp/goliath/rest"
    "reflect"
)

type Auth struct {
    rest.ResourceBase
}

func (res Auth) Define() (*rest.ResourceDefine) {
    return &rest.ResourceDefine{
        Path:    "/account/auth",
        Methods: map[string]rest.ResourceMethodDefine{
            "POST": {
                Summary:       "アカウント認証",
                Description:   "アカウントトークンでアカウントを認証します。",
                UrlParameters: []rest.UrlParameter{},
                PostParameters: map[string]rest.PostParameter{
                    "Token": {
                        Type:        reflect.String,
                        Description: "アカウントトークン",
                        Regex:       "/[0-9A-F]{53,53}/",
                        Require:     true}},
                Returns: map[string]rest.Return{
                    "AccountInfo": {
                        Type:        reflect.Map,
                        Description: "アカウント情報"}},
                RequireAuthentication: false,
                IsDebugModeOnly:       false,
                RunInMaintenance:      false}}}
}

func (res Auth) Post(request *rest.Request, response *rest.Response) (error) {
    // パラメータを取得
    token := request.GetPostString("Token", "")

    // アカウントを取得
    am := rest.GetAccountManager()
    account, err := am.GetAccountByToken(token)
    if err != nil {
        return err
    }

    // アカウント取得チェック
    if account == nil {
        response.SetErrorMessage("ERR_RES_122")
    } else if account.IsBan {
        response.SetErrorMessage("ERR_RES_121")
    }

    // セッションに登録
    if response.ResultCode == rest.ResultOK {

    }

    // 戻り値に値をセット
    response.Result = map[string]interface{}{"AccountInfo": account}




    return nil
}

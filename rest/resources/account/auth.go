package account

import (
    "github.com/shimalab-jp/goliath/rest"
)

type Auth struct {
    rest.ResourceBase
}

func (res Auth) Define() (*rest.ResourceInfo) {
    return &rest.ResourceInfo{
        Path:    "/account/auth",
        Methods: map[string]rest.ResourceDefine{
            "POST": {
                Summary:       "アカウント認証",
                Description:   "アカウントトークンでアカウントを認証します。",
                UrlParameters: map[string]rest.Parameter{},
                PostParameters: map[string]rest.Parameter{
                    "Token": {
                        Type:        "string",
                        Description: "アカウントトークン",
                        Regex:       "/[0-9A-F]{53,53}/"}},
                Returns: map[string]rest.Return{
                    "accountInfo": {
                        Type:        "array",
                        Description: "アカウント認証情報"}},
                RequireAuthentication: false,
                IsDebugModeOnly:       false,
                RunInMaintenance:      false}}}
}

func (res Auth) Post(request *rest.Request, response *rest.Response) (error) {
    // パラメータを取得
    token := request.GetString("Token", "")

    // アカウントを取得
    am := rest.GetAccountManager()
    account, err := am.GetAccountByToken(token)
    if err != nil {
        return err
    }

    // アカウント取得チェック
    if account == nil {
        response.SetErrorMessage("ERR_ACC_111")
    } else if account.IsBan {
        response.SetErrorMessage("ERR_ACC_102")
    }

    // セッションに登録
    if response.ResultCode == rest.ResultOK {

    }

    // 戻り値に値をセット
    response.Result = map[string]interface{}{"AccountInfo": account}




    return nil
}

package v1

import (
    "github.com/shimalab-jp/goliath/rest"
    "reflect"
)

type AccountPassword struct {
    rest.ResourceBase
}

func (res AccountPassword) Define() *rest.ResourceDefine {
    return &rest.ResourceDefine{
        Methods: map[string]rest.ResourceMethodDefine{
            "POST": {
                Summary:         "パスワードリセット",
                Description:     "データ移行用のパスワードをリセットします。",
                UrlParameters:   []rest.UrlParameter{},
                QueryParameters: map[string]rest.QueryParameter{},
                PostParameters:  map[string]rest.PostParameter{},
                Returns: map[string]rest.Return{
                    "AccountInfo": {
                        Type:        reflect.Map,
                        Description: "アカウント情報"}},
                RequireAuthentication: true,
                IsDebugModeOnly:       false,
                RunInMaintenance:      false}}}
}

func (res AccountPassword) Post(request *rest.Request, response *rest.Response) error {
    // アカウントを作成
    am := rest.GetAccountManager(request, response)
    account, err := am.RenewPassword(request.Account.Token)
    if err != nil {
        return err
    }

    // 戻り値に値をセット
    if response.ResultCode == rest.ResultOK {
        response.Result = map[string]interface{}{"AccountInfo": account.Output()}
    }

    return nil
}

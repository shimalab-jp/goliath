package resources

import (
    "github.com/shimalab-jp/goliath/config"
    "github.com/shimalab-jp/goliath/rest"
    "github.com/shimalab-jp/goliath/util"
    "reflect"
)

type Reference struct {
    rest.ResourceBase
    Resources *map[string]*rest.IRestResource
}

func (res Reference) Define() (*rest.ResourceDefine) {
    return &rest.ResourceDefine{
        Path:    "/reference",
        Methods: map[string]rest.ResourceMethodDefine{
            "GET": {
                Summary:       "RESTリソース一覧取得",
                Description:   "APIリファレンス用に、RESTリソースの一覧を取得します。",
                UrlParameters: map[string]rest.Parameter{},
                PostParameters: map[string]rest.Parameter{},
                Returns: map[string]rest.Return{
                    "Resources": {
                        Type:        reflect.Map,
                        Description: "リファレンスデータ"}},
                RequireAuthentication: false,
                IsDebugModeOnly:       false,
                RunInMaintenance:      true}}}
}

func (res Reference) Get(request *rest.Request, response *rest.Response) (error) {
    // リファレンス情報を取得
    envName := config.Values.Server.Reference.Environment
    envCode := util.ToEnvironmentCode(envName)
    className := ""
    switch envCode {
    case util.EnvironmentDemo:
        envName = "DEMO"
        className = "env_name_demo"
        break
    case util.EnvironmentDevelop1:
        envName = "DEVELOP1"
        className = "env_name_dev1"
        break
    case util.EnvironmentDevelop2:
        envName = "DEVELOP2"
        className = "env_name_dev2"
        break
    case util.EnvironmentTest:
        envName = "TEST"
        className = "env_name_test"
        break
    case util.EnvironmentAppleReview:
        envName = "APPLE REVIEW"
        className = "env_name_apl"
        break
    case util.EnvironmentStaging:
        envName = "STAGING"
        className = "env_name_stg"
        break
    case util.EnvironmentProduction:
        envName = "PRODUCTION"
        className = "env_name_prd"
        break
    default:
        envName = "LOCAL"
        className = "env_name_loc"
        break
    }

    // リファレンスデータを作成
    ret := map[string]rest.ResourceDefine{}
    for i, v := range *res.Resources {
        if v != nil && *v != nil {
            ret[i] = *(*v).Define()
        }
    }

    // 結果に格納
    response.Result = map[string]interface{}{
        "Name": config.Values.Server.Reference.Name,
        "EnvName": config.Values.Server.Reference.Environment,
        "EnvCode": envCode,
        "EnvClass": className,
        "Logo": config.Values.Server.Reference.Logo,
        "Resources": ret}
    return nil
}

package resources

import (
    "github.com/shimalab-jp/goliath/rest"
    "reflect"
)

type Reference struct {
    rest.ResourceBase
    Resources *map[string]*rest.IRestResource
}

func (res Reference) Define() (*rest.ResourceInfo) {
    return &rest.ResourceInfo{
        Path:    "/reference",
        Methods: map[string]rest.ResourceDefine{
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
    ret := map[string]rest.ResourceInfo{}
    for i, v := range *res.Resources {
        if v != nil && *v != nil {
            ret[i] = *(*v).Define()
        }
    }
    response.Result = map[string]interface{}{ "Resources": ret }
    return nil
}

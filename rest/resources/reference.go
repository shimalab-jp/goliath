package resources

import (
    "github.com/shimalab-jp/goliath/rest"
)

type Reference struct {
    rest.ResourceBase
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
                        Type:        "array",
                        Description: "リファレンスデータ"}},
                RequireAuthentication: false,
                IsDebugModeOnly:       false,
                RunInMaintenance:      true}}}
}

func (res Reference) Get(request *rest.Request, response *rest.Response) (error) {
    return nil
}

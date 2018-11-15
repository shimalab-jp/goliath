package v1

import (
    "github.com/shimalab-jp/goliath/rest"
)

type DebugCache struct {
    rest.ResourceBase
}

func (res DebugCache) Define() (*rest.ResourceDefine) {
    return &rest.ResourceDefine{
        Methods: map[string]rest.ResourceMethodDefine{
            "DELETE": {
                Summary:               "キャッシュクリア",
                Description:           "memcached及び内部キャッシュの全ての値をクリアします。",
                UrlParameters:         []rest.UrlParameter{},
                PostParameters:        map[string]rest.PostParameter{},
                Returns:               map[string]rest.Return{},
                RequireAuthentication: false,
                IsDebugModeOnly:       false,
                RunInMaintenance:      true}}}
}

func (res DebugCache) Delete(request *rest.Request, response *rest.Response) (error) {
    return nil
}

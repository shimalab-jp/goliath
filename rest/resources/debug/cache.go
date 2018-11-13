package debug

import (
    "github.com/shimalab-jp/goliath/rest"
)

type Cache struct {
    rest.ResourceBase
}

func (res Cache) Define() (*rest.ResourceDefine) {
    return &rest.ResourceDefine{
        Path:    "/debug/cache",
        Methods: map[string]rest.ResourceMethodDefine{
            "DELETE": {
                Summary:               "キャッシュクリア",
                Description:           "memcached及び内部キャッシュの全ての値をクリアします。",
                UrlParameters:         map[string]rest.Parameter{},
                PostParameters:        map[string]rest.Parameter{},
                Returns:               map[string]rest.Return{},
                RequireAuthentication: false,
                IsDebugModeOnly:       false,
                RunInMaintenance:      true}}}
}

func (res Cache) Delete(request *rest.Request, response *rest.Response) (error) {
    return nil
}

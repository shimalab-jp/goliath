package rest

import (
    "fmt"
    "github.com/pkg/errors"
    "net/http"
    "reflect"
    "strings"
    "sync"

    "github.com/shimalab-jp/goliath/log"
)

type Return struct {
    Type        reflect.Kind
    Description string
}

type Parameter struct {
    Type        reflect.Kind
    Description string
    Range       []float64
    Select      []interface{}
    Regex       string
    Default     interface{}
    Require     bool
}

type ResourceMethodDefine struct {
    Summary               string
    Description           string
    UrlParameters         map[string]Parameter
    PostParameters        map[string]Parameter
    Returns               map[string]Return
    RequireAuthentication bool
    IsDebugModeOnly       bool
    IsAdminModeOnly       bool
    RunInMaintenance      bool
}

type ResourceDefine struct {
    Path     string
    Methods  map[string]ResourceMethodDefine
}

type IRestResource interface {
    Define() (*ResourceDefine)
    Get(request *Request, response *Response) (error)
    Post(request *Request, response *Response) (error)
    Delete(request *Request, response *Response) (error)
    Put(request *Request, response *Response) (error)
}

type ResourceBase struct{}

func (res *ResourceBase) Define() (*ResourceDefine) {
    return &ResourceDefine{ Methods: map[string]ResourceMethodDefine{} }
}

func (res *ResourceBase) Get(request *Request, response *Response) (error) {
    response.StatusCode = http.StatusMethodNotAllowed
    response.ResultCode = ResultNotImplemented
    return nil
}

func (res *ResourceBase) Post(request *Request, response *Response) (error) {
    response.StatusCode = http.StatusMethodNotAllowed
    response.ResultCode = ResultNotImplemented
    return nil
}

func (res *ResourceBase) Delete(request *Request, response *Response) (error) {
    response.StatusCode = http.StatusMethodNotAllowed
    response.ResultCode = ResultNotImplemented
    return nil
}

func (res *ResourceBase) Put(request *Request, response *Response) (error) {
    response.StatusCode = http.StatusMethodNotAllowed
    response.ResultCode = ResultNotImplemented
    return nil
}

type resourceManager struct {
    resourceMutex *sync.Mutex
    resources     *map[string]*IRestResource
}

func createResourceManager() (*resourceManager) {
    return &resourceManager{
        resourceMutex: &sync.Mutex{},
        resources:     &map[string]*IRestResource{}}
}

func (resManager *resourceManager) initialize() {

}

func (resManager *resourceManager) FindResource(path string) (*IRestResource) {
    findPath := "/" + strings.TrimLeft(path, "/")

    resManager.resourceMutex.Lock()
    ret, ok := (*resManager.resources)[findPath]
    resManager.resourceMutex.Unlock()

    if ok {
        return ret
    } else {
        return nil
    }
}

func (resManager *resourceManager) Append(resource *IRestResource) (error) {
    var err error = nil

    if resource == nil {
        err = errors.New("resource is nil.")
    } else {
        resManager.resourceMutex.Lock()

        path := (*resource).Define().Path
        findPath := strings.ToLower("/" + strings.TrimLeft(path, "/"))
        if _, ok := (*resManager.resources)[findPath]; !ok {
            (*resManager.resources)[findPath] = resource
            define := (*resource).Define()
            for method := range (*define).Methods {
                log.D("[ResourceManager] Append resource %s:%s", strings.ToUpper(method), findPath)
            }
        } else {
            err = errors.New(fmt.Sprintf("path '%s' is duplicated.", path))
        }

        resManager.resourceMutex.Unlock()
    }

    return err
}

func (resManager *resourceManager) GetAllResources() (*map[string]*IRestResource) {
    return resManager.resources
}

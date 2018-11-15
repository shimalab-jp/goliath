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

type UrlParameter struct {
    Name        string
    Type        reflect.Kind
    Description string
    Range       []float64
    Select      []interface{}
    Regex       string
    Default     interface{}
    Require     bool
}

type PostParameter struct {
    Type              reflect.Kind
    IsMultilineString bool
    Description       string
    Range             []float64
    Select            []interface{}
    Regex             string
    Default           interface{}
    Require           bool
}

type ResourceMethodDefine struct {
    Summary               string
    Description           string
    UrlParameters         []UrlParameter
    PostParameters        map[string]PostParameter
    Returns               map[string]Return
    RequireAuthentication bool
    IsDebugModeOnly       bool
    IsAdminModeOnly       bool
    RunInMaintenance      bool
}

type ResourceDefine struct {
    Methods map[string]ResourceMethodDefine
}

type IRestResource interface {
    Define() (*ResourceDefine)
    Get(request *Request, response *Response) (error)
    Post(request *Request, response *Response) (error)
    Delete(request *Request, response *Response) (error)
    Put(request *Request, response *Response) (error)
    SetVersion(version uint32)
    GetVersion() (uint32)
    SetPath(path string)
    GetPath() (string)
}

type ResourceBase struct {
    version uint32
    path    string
}

func (res *ResourceBase) Define() (*ResourceDefine) {
    return &ResourceDefine{Methods: map[string]ResourceMethodDefine{}}
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

func (res *ResourceBase) SetVersion(version uint32) {
    res.version = version
}

func (res *ResourceBase) GetVersion() (uint32) {
    return res.version
}

func (res *ResourceBase) SetPath(path string) {
    res.path = path
}

func (res *ResourceBase) GetPath() (string) {
    return res.path
}








type resourceManager struct {
    resourceMutex *sync.Mutex
    resources     *map[uint32]*map[string]*IRestResource
}

func createResourceManager() (*resourceManager) {
    return &resourceManager{
        resourceMutex: &sync.Mutex{},
        resources:     &map[uint32]*map[string]*IRestResource{}}
}

func (resManager *resourceManager) initialize() {

}

func (resManager *resourceManager) FindResource(version uint32, path string) (*IRestResource) {
    findPath := "/" + strings.TrimLeft(path, "/")

    var ret *IRestResource = nil
    resManager.resourceMutex.Lock()
    if ver, ok := (*resManager.resources)[version]; ok {
        if res, ok := (*ver)[findPath]; ok {
            ret = res
        }
    }
    resManager.resourceMutex.Unlock()

    return ret
}

func (resManager *resourceManager) Append(version uint32, path string, resource *IRestResource) (error) {
    var err error = nil

    if resource == nil {
        err = errors.New("resource is nil.")
    } else {
        resManager.resourceMutex.Lock()

        findPath := strings.ToLower("/" + strings.TrimLeft(path, "/"))

        if ver, ok := (*resManager.resources)[version]; ok {
            if _, ok := (*ver)[findPath]; !ok {
                (*ver)[findPath] = resource
                define := (*resource).Define()
                for method := range (*define).Methods {
                    log.D("[ResourceManager] Append resource %s:%s", strings.ToUpper(method), findPath)
                }
                (*resource).SetVersion(version)
                (*resource).SetPath(findPath)
            } else {
                err = errors.New(fmt.Sprintf("path '%s' is duplicated.", path))
            }
        } else {
            (*resManager.resources)[version] = &map[string]*IRestResource{findPath: resource}
            (*resource).SetVersion(version)
            (*resource).SetPath(findPath)
        }

        resManager.resourceMutex.Unlock()
    }

    return err
}

func (resManager *resourceManager) GetAllResources() (*map[uint32]*map[string]*IRestResource) {
    return resManager.resources
}

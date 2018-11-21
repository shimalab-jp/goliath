package rest

import (
    "fmt"
    "github.com/pkg/errors"
    "net/http"
    "reflect"
    "regexp"
    "strings"
    "sync"

    "github.com/shimalab-jp/goliath/log"
)

type ParameterType uint

const (
    Invalid ParameterType = iota
    PostParam
    QueryParam
    UrlParam
)

type Return struct {
    Type        reflect.Kind
    Description string
}

type Parameter interface {
    GetType() reflect.Kind
    GetDescription() string
    GetRange() []float64
    GetSelect() []interface{}
    GetRegex() string
    GetDefault() interface{}
    GetRequire() bool
    GetIsMultiple() bool
}

type UrlParameter struct {
    Type        reflect.Kind
    Description string
    Range       []float64
    Select      []interface{}
    Regex       string
    Default     interface{}
    Require     bool
    IsMultiple  bool
}

func (p *UrlParameter) GetType() reflect.Kind {
    return p.Type
}

func (p *UrlParameter) GetDescription() string {
    return p.Description
}

func (p *UrlParameter) GetRange() []float64 {
    return p.Range
}

func (p *UrlParameter) GetSelect() []interface{} {
    return p.Select
}

func (p *UrlParameter) GetRegex() string {
    return p.Regex
}

func (p *UrlParameter) GetDefault() interface{} {
    return p.Default
}

func (p *UrlParameter) GetRequire() bool {
    return p.Require
}

func (p *UrlParameter) GetIsMultiple() bool {
    return p.IsMultiple
}

type QueryParameter struct {
    Type        reflect.Kind
    Description string
    Range       []float64
    Select      []interface{}
    Regex       string
    Default     interface{}
    Require     bool
}

func (p *QueryParameter) GetType() reflect.Kind {
    return p.Type
}

func (p *QueryParameter) GetDescription() string {
    return p.Description
}

func (p *QueryParameter) GetRange() []float64 {
    return p.Range
}

func (p *QueryParameter) GetSelect() []interface{} {
    return p.Select
}

func (p *QueryParameter) GetRegex() string {
    return p.Regex
}

func (p *QueryParameter) GetDefault() interface{} {
    return p.Default
}

func (p *QueryParameter) GetRequire() bool {
    return p.Require
}

func (p *QueryParameter) GetIsMultiple() bool {
    return false
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

func (p *PostParameter) GetType() reflect.Kind {
    return p.Type
}

func (p *PostParameter) GetDescription() string {
    return p.Description
}

func (p *PostParameter) GetRange() []float64 {
    return p.Range
}

func (p *PostParameter) GetSelect() []interface{} {
    return p.Select
}

func (p *PostParameter) GetRegex() string {
    return p.Regex
}

func (p *PostParameter) GetDefault() interface{} {
    return p.Default
}

func (p *PostParameter) GetRequire() bool {
    return p.Require
}

func (p *PostParameter) GetIsMultiple() bool {
    return false
}

type ResourceMethodDefine struct {
    Summary               string
    Description           string
    UrlParameters         []UrlParameter
    QueryParameters       map[string]QueryParameter
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

type ResourceManager struct {
    initialized   bool
    resourceMutex *sync.Mutex
    resources     *map[uint32]*map[string]*IRestResource
}

func createResourceManager() (*ResourceManager) {
    return &ResourceManager{
        initialized:   true,
        resourceMutex: &sync.Mutex{},
        resources:     &map[uint32]*map[string]*IRestResource{}}
}

func (rm *ResourceManager) FindResource(version uint32, path string) (string, *IRestResource) {
    if !rm.initialized {
        return "", nil
    }

    findPath := "/" + strings.TrimRight(strings.TrimLeft(path, "/"), "/")

    apiPath := ""
    var ret *IRestResource = nil
    rm.resourceMutex.Lock()

    if ver, ok := (*rm.resources)[version]; ok {
        // urlパラメータ無しとして検索
        if res, ok := (*ver)[findPath]; ok {
            apiPath = findPath
            ret = res
        }

        // urlパラメータ付きとして検索
        if ret == nil {
            for p, res := range *ver {
                searchPattern := fmt.Sprintf("^%s/.+$", p)
                regex := regexp.MustCompile(searchPattern)
                regexResult := regex.FindAllStringSubmatch(findPath, -1)
                if len(regexResult) > 0 && len(regexResult[0]) > 0 {
                    apiPath = p
                    ret = res
                    break
                }
            }
        }
    }

    rm.resourceMutex.Unlock()
    return apiPath, ret
}

func (rm *ResourceManager) Append(version uint32, path string, resource *IRestResource) (error) {
    if !rm.initialized {
        return errors.New("[ResourceManager] This structure has not been initialized.")
    }

    var err error = nil

    if resource == nil {
        err = errors.New("resource is nil.")
    } else {
        rm.resourceMutex.Lock()

        findPath := strings.ToLower("/" + strings.TrimLeft(path, "/"))

        if ver, ok := (*rm.resources)[version]; ok {
            if _, ok := (*ver)[findPath]; !ok {
                // 既存リソースとの重複チェック
                for p := range *ver {
                    searchPattern := fmt.Sprintf("^%s/.+$", p)
                    regex := regexp.MustCompile(searchPattern)
                    regexResult := regex.FindAllStringSubmatch(findPath, -1)
                    if len(regexResult) > 0 && len(regexResult[0]) > 0 {
                        err = errors.New(fmt.Sprintf("[ResourceManager] Failed to append resource. Duplicate key '%s' and '%s' detected.", path, p))
                        break
                    }
                }
                {
                    searchPattern := fmt.Sprintf("^%s/.+$", findPath)
                    regex := regexp.MustCompile(searchPattern)
                    for p := range *ver {
                        regexResult := regex.FindAllStringSubmatch(p, -1)
                        if len(regexResult) > 0 && len(regexResult[0]) > 0 {
                            err = errors.New(fmt.Sprintf("[ResourceManager] Failed to append resource. Duplicate key '%s' and '%s' detected.", path, p))
                            break
                        }
                    }
                }

                if err == nil {
                    (*ver)[findPath] = resource
                    define := (*resource).Define()
                    for method := range (*define).Methods {
                        log.D("[ResourceManager] Append resource %s:%s", strings.ToUpper(method), findPath)
                    }
                    (*resource).SetVersion(version)
                    (*resource).SetPath(findPath)
                }
            } else {
                err = errors.New(fmt.Sprintf("[ResourceManager] Failed to append resource. Duplicate key '%s' and '%s' detected.", path, findPath))
            }
        } else {
            (*rm.resources)[version] = &map[string]*IRestResource{findPath: resource}
            (*resource).SetVersion(version)
            (*resource).SetPath(findPath)
        }

        rm.resourceMutex.Unlock()
    }

    return err
}

func (rm *ResourceManager) GetAllResources() (*map[uint32]*map[string]*IRestResource) {
    if !rm.initialized {
        return nil
    }

    return rm.resources
}

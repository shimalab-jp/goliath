package rest

import (
    "fmt"
    "github.com/shimalab-jp/goliath/message"
    "strconv"
)

type Request struct {
    Version        uint32
    BaseUrl        string
    BusinessDay    businessDay
    RemoteAddress  string
    UserAgent      string
    Languages      []message.AcceptLanguage
    Url            string
    NameSpace      string
    Name           string
    SubNames       string
    DisplayName    string
    Method         string
    ContentType    string
    QueryString    string
    Headers        map[string]string
    GetData        map[string]interface{}
    PostData       map[string]interface{}
    OutputFormat   string
    Resource       *IRestResource
    MethodInfo     *ResourceMethodDefine
    ResourceDefine *ResourceDefine
    Token          string
    Account        *accountInfo
    ParseTime      int64
    RequestID      string
    MessageManager *message.GoliathMessageManager
}

func (request *Request) GetPostString(name string, defaultValue string) (string) {
    if value, ok := request.PostData[name]; ok {
        return fmt.Sprint(value)
    }
    return defaultValue
}

func (request *Request) GetPostInt8(name string, defaultValue int8) (int8) {
    ret, err := strconv.ParseInt(request.GetPostString(name, string(defaultValue)), 10, 8)
    if err == nil {
        return int8(ret)
    }
    return defaultValue
}

func (request *Request) GetPostInt16(name string, defaultValue int16) (int16) {
    ret, err := strconv.ParseInt(request.GetPostString(name, string(defaultValue)), 10, 16)
    if err == nil {
        return int16(ret)
    }
    return defaultValue
}

func (request *Request) GetPostInt32(name string, defaultValue int32) (int32) {
    ret, err := strconv.ParseInt(request.GetPostString(name, string(defaultValue)), 10, 32)
    if err == nil {
        return int32(ret)
    }
    return defaultValue
}

func (request *Request) GetPostInt64(name string, defaultValue int64) (int64) {
    ret, err := strconv.ParseInt(request.GetPostString(name, string(defaultValue)), 10, 64)
    if err == nil {
        return int64(ret)
    }
    return defaultValue
}

func (request *Request) GetPostUInt8(name string, defaultValue uint8) (uint8) {
    ret, err := strconv.ParseUint(request.GetPostString(name, string(defaultValue)), 10, 8)
    if err == nil {
        return uint8(ret)
    }
    return defaultValue
}

func (request *Request) GetPostUInt16(name string, defaultValue uint16) (uint16) {
    ret, err := strconv.ParseUint(request.GetPostString(name, string(defaultValue)), 10, 16)
    if err == nil {
        return uint16(ret)
    }
    return defaultValue
}

func (request *Request) GetPostUInt32(name string, defaultValue uint32) (uint32) {
    ret, err := strconv.ParseUint(request.GetPostString(name, string(defaultValue)), 10, 32)
    if err == nil {
        return uint32(ret)
    }
    return defaultValue
}

func (request *Request) GetPostUInt64(name string, defaultValue uint64) (uint64) {
    ret, err := strconv.ParseUint(request.GetPostString(name, string(defaultValue)), 10, 64)
    if err == nil {
        return uint64(ret)
    }
    return defaultValue
}

func (request *Request) GetPostFloat32(name string, defaultValue float32) (float32) {
    ret, err := strconv.ParseFloat(request.GetPostString(name, ""), 32)
    if err == nil {
        return float32(ret)
    }
    return defaultValue
}

func (request *Request) GetPostFloat64(name string, defaultValue float64) (float64) {
    ret, err := strconv.ParseFloat(request.GetPostString(name, ""), 64)
    if err == nil {
        return ret
    }
    return defaultValue
}

func (request *Request) GetPostBoolean(name string, defaultValue bool) (bool) {
    ret, err := strconv.ParseBool(request.GetPostString(name, "false"))
    if err == nil {
        return ret
    }
    return defaultValue
}

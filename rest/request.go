package rest

import (
    "fmt"
    "github.com/shimalab-jp/goliath/log"
    "github.com/shimalab-jp/goliath/message"
    "strconv"
    "strings"
)

type Request struct {
    BusinessDay   businessDay
    RemoteAddress string
    UserAgent     string
    Languages     []message.AcceptLanguage
    Url           string
    NameSpace     string
    Name          string
    SubNames      string
    DisplayName   string
    Method        string
    ContentType   string
    QueryString   string
    Headers       map[string]string
    GetData       map[string]string
    PostData      map[string]interface{}
    OutputFormat  string
    Resource      *IRestResource
    ResourceInfo  *ResourceInfo
    Session       *Session
    ParseTime     int64
    RequestID     string
    MessageManager *message.GoliathMessageManager
}

func (request *Request) GetString(name string, defaultValue string) (string) {
    if strings.ToUpper(request.Method) == "POST" {
        if value, ok := request.PostData[name]; ok {
            return fmt.Sprint(value)
        }
    }

    if value, ok := request.GetData[name]; ok {
        return fmt.Sprint(value)
    }

    log.D(message.SystemMessageManager.Get("DBG_RES_101", name))

    return defaultValue
}

func (request *Request) GetInt8(name string, defaultValue int8) (int8) {
    ret, err := strconv.ParseInt(request.GetString(name, string(defaultValue)), 10, 8)
    if err == nil {
        return int8(ret)
    }
    return defaultValue
}

func (request *Request) GetInt16(name string, defaultValue int16) (int16) {
    ret, err := strconv.ParseInt(request.GetString(name, string(defaultValue)), 10, 16)
    if err == nil {
        return int16(ret)
    }
    return defaultValue
}

func (request *Request) GetInt32(name string, defaultValue int32) (int32) {
    ret, err := strconv.ParseInt(request.GetString(name, string(defaultValue)), 10, 32)
    if err == nil {
        return int32(ret)
    }
    return defaultValue
}

func (request *Request) GetInt64(name string, defaultValue int64) (int64) {
    ret, err := strconv.ParseInt(request.GetString(name, string(defaultValue)), 10, 64)
    if err == nil {
        return int64(ret)
    }
    return defaultValue
}

func (request *Request) GetUInt8(name string, defaultValue uint8) (uint8) {
    ret, err := strconv.ParseUint(request.GetString(name, string(defaultValue)), 10, 8)
    if err == nil {
        return uint8(ret)
    }
    return defaultValue
}

func (request *Request) GetUInt16(name string, defaultValue uint16) (uint16) {
    ret, err := strconv.ParseUint(request.GetString(name, string(defaultValue)), 10, 16)
    if err == nil {
        return uint16(ret)
    }
    return defaultValue
}

func (request *Request) GetUInt32(name string, defaultValue uint32) (uint32) {
    ret, err := strconv.ParseUint(request.GetString(name, string(defaultValue)), 10, 32)
    if err == nil {
        return uint32(ret)
    }
    return defaultValue
}

func (request *Request) GetUInt64(name string, defaultValue uint64) (uint64) {
    ret, err := strconv.ParseUint(request.GetString(name, string(defaultValue)), 10, 64)
    if err == nil {
        return uint64(ret)
    }
    return defaultValue
}

func (request *Request) GetFloat32(name string, defaultValue float32) (float32) {
    ret, err := strconv.ParseFloat(request.GetString(name, ""), 32)
    if err == nil {
        return float32(ret)
    }
    return defaultValue
}

func (request *Request) GetFloat64(name string, defaultValue float64) (float64) {
    ret, err := strconv.ParseFloat(request.GetString(name, ""), 64)
    if err == nil {
        return ret
    }
    return defaultValue
}

func (request *Request) GetBoolean(name string, defaultValue bool) (bool) {
    ret, err := strconv.ParseBool(request.GetString(name, "false"))
    if err == nil {
        return ret
    }
    return defaultValue
}

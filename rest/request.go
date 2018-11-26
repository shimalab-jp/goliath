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
    UrlData        []string
    QueryData      map[string]string
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

func (request *Request) GetParamString(pType ParameterType, name string, defaultValue string) (string, bool) {
    switch pType {
    case PostParam:
        if value, ok := request.PostData[name]; ok {
            return fmt.Sprint(value), true
        }
        break
    case QueryParam:
        if value, ok := request.QueryData[name]; ok {
            return value, true
        }
        break
    case UrlParam:
        val, err := strconv.ParseInt(name, 10, 64)
        if err == nil {
            if len(request.UrlData) < int(val) {
                return request.UrlData[val], true
            }
        }
        break
    }
    return defaultValue, false
}

func (request *Request) GetParamInt8(pType ParameterType, name string, defaultValue int8) (int8, bool) {
    if val, ok := request.GetParamString(pType, name, ""); ok {
        ret, err := strconv.ParseInt(val, 10, 8)
        if err == nil {
            return int8(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetParamInt16(pType ParameterType, name string, defaultValue int16) (int16, bool) {
    if val, ok := request.GetParamString(pType, name, ""); ok {
        ret, err := strconv.ParseInt(val, 10, 16)
        if err == nil {
            return int16(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetParamInt32(pType ParameterType, name string, defaultValue int32) (int32, bool) {
    if val, ok := request.GetParamString(pType, name, ""); ok {
        ret, err := strconv.ParseInt(val, 10, 32)
        if err == nil {
            return int32(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetParamInt64(pType ParameterType, name string, defaultValue int64) (int64, bool) {
    if val, ok := request.GetParamString(pType, name, ""); ok {
        ret, err := strconv.ParseInt(val, 10, 64)
        if err == nil {
            return int64(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetParamUInt8(pType ParameterType, name string, defaultValue uint8) (uint8, bool) {
    if val, ok := request.GetParamString(pType, name, ""); ok {
        ret, err := strconv.ParseUint(val, 10, 8)
        if err == nil {
            return uint8(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetParamUInt16(pType ParameterType, name string, defaultValue uint16) (uint16, bool) {
    if val, ok := request.GetParamString(pType, name, ""); ok {
        ret, err := strconv.ParseUint(val, 10, 16)
        if err == nil {
            return uint16(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetParamUInt32(pType ParameterType, name string, defaultValue uint32) (uint32, bool) {
    if val, ok := request.GetParamString(pType, name, ""); ok {
        ret, err := strconv.ParseUint(val, 10, 32)
        if err == nil {
            return uint32(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetParamUInt64(pType ParameterType, name string, defaultValue uint64) (uint64, bool) {
    if val, ok := request.GetParamString(pType, name, ""); ok {
        ret, err := strconv.ParseUint(val, 10, 64)
        if err == nil {
            return uint64(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetParamFloat32(pType ParameterType, name string, defaultValue float32) (float32, bool) {
    if val, ok := request.GetParamString(pType, name, ""); ok {
        ret, err := strconv.ParseFloat(val, 32)
        if err == nil {
            return float32(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetParamFloat64(pType ParameterType, name string, defaultValue float64) (float64, bool) {
    if val, ok := request.GetParamString(pType, name, ""); ok {
        ret, err := strconv.ParseFloat(val, 64)
        if err == nil {
            return float64(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetParamBoolean(pType ParameterType, name string, defaultValue bool) (bool, bool) {
    if val, ok := request.GetParamString(pType, name, ""); ok {
        ret, err := strconv.ParseBool(val)
        if err == nil {
            return ret, true
        }
    }
    return defaultValue, false
}

func (request *Request) GetUrlParamString(index int, defaultValue string) (string, bool) {
    if len(request.UrlData) > index {
        return request.UrlData[index], true
    }
    return defaultValue, false
}

func (request *Request) GetUrlParamInt8(index int, defaultValue int8) (int8, bool) {
    if val, ok := request.GetUrlParamString(index, ""); ok {
        ret, err := strconv.ParseInt(val, 10, 8)
        if err == nil {
            return int8(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetUrlParamInt16(index int, defaultValue int16) (int16, bool) {
    if val, ok := request.GetUrlParamString(index, ""); ok {
        ret, err := strconv.ParseInt(val, 10, 16)
        if err == nil {
            return int16(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetUrlParamInt32(index int, defaultValue int32) (int32, bool) {
    if val, ok := request.GetUrlParamString(index, ""); ok {
        ret, err := strconv.ParseInt(val, 10, 32)
        if err == nil {
            return int32(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetUrlParamInt64(index int, defaultValue int64) (int64, bool) {
    if val, ok := request.GetUrlParamString(index, ""); ok {
        ret, err := strconv.ParseInt(val, 10, 64)
        if err == nil {
            return int64(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetUrlParamUInt8(index int, defaultValue uint8) (uint8, bool) {
    if val, ok := request.GetUrlParamString(index, ""); ok {
        ret, err := strconv.ParseUint(val, 10, 8)
        if err == nil {
            return uint8(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetUrlParamUInt16(index int, defaultValue uint16) (uint16, bool) {
    if val, ok := request.GetUrlParamString(index, ""); ok {
        ret, err := strconv.ParseUint(val, 10, 16)
        if err == nil {
            return uint16(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetUrlParamUInt32(index int, defaultValue uint32) (uint32, bool) {
    if val, ok := request.GetUrlParamString(index, ""); ok {
        ret, err := strconv.ParseUint(val, 10, 32)
        if err == nil {
            return uint32(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetUrlParamUInt64(index int, defaultValue uint64) (uint64, bool) {
    if val, ok := request.GetUrlParamString(index, ""); ok {
        ret, err := strconv.ParseUint(val, 10, 64)
        if err == nil {
            return uint64(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetUrlParamFloat32(index int, defaultValue float32) (float32, bool) {
    if val, ok := request.GetUrlParamString(index, ""); ok {
        ret, err := strconv.ParseFloat(val, 32)
        if err == nil {
            return float32(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetUrlParamFloat64(index int, defaultValue float64) (float64, bool) {
    if val, ok := request.GetUrlParamString(index, ""); ok {
        ret, err := strconv.ParseFloat(val, 64)
        if err == nil {
            return float64(ret), true
        }
    }
    return defaultValue, false
}

func (request *Request) GetUrlParamBoolean(index int, defaultValue bool) (bool, bool) {
    if val, ok := request.GetUrlParamString(index, ""); ok {
        ret, err := strconv.ParseBool(val)
        if err == nil {
            return ret, true
        }
    }
    return defaultValue, false
}

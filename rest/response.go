package rest

import (
    "encoding/json"
    "github.com/shimalab-jp/goliath/message"
    "time"

    "github.com/shimalab-jp/goliath/log"
)

type MaintenanceInfo struct {
    StartTime int64
    EndTime   int64
    Subject   string
    Body      string
}

type Response struct {
    Method            string
    StatusCode        int
    Name              string
    RemoteAddress     string
    ResultCode        int32
    Result            map[string]interface{}
    Notify            interface{}
    StateData         interface{}
    ErrorCode         string
    ErrorMessage      string
    DebugMessage      string
    LinkUrl           string
    ResourceVersion   string
    Memory            interface{}
    ProcessTime       int64
    Times             map[string]int64
    MaintenanceInfo   MaintenanceInfo
    CurrentServerTime int64
    OutputFormat      string
    ContentType       string
    ContentData       interface{}
    ContentLength     int64
    MessageManager    *message.GoliathMessageManager
}

type OutputData struct {
    Method          string
    ApiName         string
    ResultCode      int32
    Result          map[string]interface{}
    ErrorCode       string
    ErrorMessage    string
    DebugMessage    string
    ResourceVersion string
    LinkUrl         string
    ProcessTime     int64
    ServerTime      int64
    MaintenanceInfo MaintenanceInfo
}

func (response *Response) CreateOutputData() (*OutputData) {
    return &OutputData{
        Method: response.Method,
        ApiName: response.Name,
        ResultCode: response.ResultCode,
        Result: response.Result,
        ErrorCode: response.ErrorCode,
        ErrorMessage: response.ErrorMessage,
        DebugMessage: response.DebugMessage,
        ResourceVersion: response.ResourceVersion,
        LinkUrl:response.LinkUrl,
        ProcessTime: response.ProcessTime,
        ServerTime: time.Now().UnixNano(),
        MaintenanceInfo: response.MaintenanceInfo}
}

func (response *Response) SetErrorMessage(messageCode string, args ...interface{}) {
    errorMessage, resultCode := response.MessageManager.GetWithResultCode(messageCode, args...)
    response.ErrorMessage = errorMessage
    response.ErrorCode = messageCode
    response.ResultCode = resultCode
}

func (response *Response) SetSystemErrorMessage(messageCode string, arg1 []interface{}, debugMessageCode string, arg2 ...interface{}) {
    errorMessage, resultCode := response.MessageManager.GetWithResultCode(messageCode, arg1...)
    response.ErrorMessage = errorMessage
    response.ErrorCode = messageCode
    response.ResultCode = resultCode
    debugMessage := response.MessageManager.Get(messageCode, arg2...)
    response.DebugMessage = debugMessage
}

func (outputData *OutputData) ToJson() (*[]byte) {
    jsonData, err := json.Marshal(outputData)
    if err != nil {
        log.Ee(err)
        return nil
    }
    return &jsonData
}


package message

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "strings"

    "github.com/shimalab-jp/goliath/config"
    "github.com/shimalab-jp/goliath/log"
    "github.com/shimalab-jp/goliath/util"
)

type AcceptLanguage struct {
    Lang string
    Q    float64
}

type messageData struct {
    Code   string
    Result int32
    Text   map[string]string
}

type globalMessageManager struct {
    initialized bool
    system      map[string]messageData
    user        map[string]messageData
}

func (txt *globalMessageManager) initialize() {
    if txt.initialized {
        return
    } else {
        path := strings.TrimSpace(os.ExpandEnv(config.Values.Message.System))
        if util.FileExists(path) {
            raw, err := ioutil.ReadFile(path)
            if err != nil {
                log.Ee(err)
            } else {
                var data []messageData
                err := json.Unmarshal(raw, &data)
                if err != nil {
                    log.Ee(err)
                } else {
                    txt.system = map[string]messageData{}
                    for _, value := range data {
                        if _, ok := txt.system[value.Code]; !ok {
                            txt.system[value.Code] = value
                        }
                    }
                    txt.initialized = true
                }
            }
        } else {
            log.E("Missing system message file! '%s'", path)
        }
    }

    if len(strings.TrimSpace(config.Values.Message.User)) > 0 {
        path := os.ExpandEnv(strings.TrimSpace(config.Values.Message.User))
        if util.FileExists(path) {
            raw, err := ioutil.ReadFile(path)
            if err != nil {
                log.Ee(err)
            } else {
                var data []messageData
                err := json.Unmarshal(raw, &data)
                if err != nil {
                    log.Ee(err)
                } else {
                    txt.user = map[string]messageData{}
                    for _, value := range data {
                        if _, ok := txt.user[value.Code]; !ok {
                            txt.user[value.Code] = value
                        }
                    }
                }
            }
        } else {
            log.E("Missing user message file! '%s'", path)
        }
    }
}

func (txt *globalMessageManager) get(langs *[]AcceptLanguage, messageCode string, args ...interface{}) (string) {
    txt.initialize()

    if messages, ok := txt.system[messageCode]; ok {
        for _, value := range *langs {
            if message, ok := messages.Text[value.Lang]; ok {
                return fmt.Sprintf(message, args...)
            }
        }

        if message, ok := messages.Text[strings.ToLower(config.Values.Message.Default)]; ok {
            return fmt.Sprintf(message, args...)
        }

        for _, value := range messages.Text {
            return fmt.Sprintf(value, args...)
        }

        log.W("Message code '%s' is not defined.", messageCode)
    }

    return ""
}

func (txt *globalMessageManager) getWithResultCode(langs *[]AcceptLanguage, messageCode string, args ...interface{}) (string, int32) {
    txt.initialize()

    if messages, ok := txt.system[messageCode]; ok {
        for _, value := range *langs {
            if message, ok := messages.Text[value.Lang]; ok {
                return fmt.Sprintf(message, args...), messages.Result
            }
        }

        if message, ok := messages.Text[strings.ToLower(config.Values.Message.Default)]; ok {
            return fmt.Sprintf(message, args...), messages.Result
        }

        for _, value := range messages.Text {
            return fmt.Sprintf(value, args...), messages.Result
        }

        log.W("Message code '%s' is not defined.", messageCode)
    }

    return "", 0
}

var manager *globalMessageManager = nil

type GoliathMessageManager struct {
    Languages *[]AcceptLanguage
}

func CreateMessageManager(langs *[]AcceptLanguage) (*GoliathMessageManager) {
    ret := &GoliathMessageManager{Languages: langs}
    ret.Get("INI_TXT_001")
    return ret
}

func (mm *GoliathMessageManager) Get(messageCode string, args ...interface{}) (string) {
    if manager == nil {
        manager = &globalMessageManager{}
    }

    return manager.get(mm.Languages, messageCode, args...)
}

func (mm *GoliathMessageManager) GetWithResultCode(messageCode string, args ...interface{}) (string, int32) {
    if manager == nil {
        manager = &globalMessageManager{}
    }

    return manager.getWithResultCode(mm.Languages, messageCode, args...)
}

var SystemMessageManager *GoliathMessageManager

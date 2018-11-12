package config

import (
    "fmt"
    "github.com/pkg/errors"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "os"
)

type ServerConfig struct {
    Port      uint16
    ApiUrl    string
    TimeZone  string
    LogLevel  uint32
    UserDB    int8
    Debug     bool
    Reference bool
    ReferenceUrl string
    ClearDB   bool
}

type MessageConfig struct {
    Default string
    System  string
    User    string
}

type MemcachedServerConfig struct {
    Host string
    Port uint16
}

type MemcachedConfig struct {
    Prefix     string
    Expiration int32
    Servers    []MemcachedServerConfig
}

type DatabaseConfig struct {
    Name     string
    Driver   string
    Host     string
    Port     uint16
    Scheme   string
    User     string
    Password string
}

func (db *DatabaseConfig) ConnectionString() (string) {
    return fmt.Sprintf(
        "%s:%s@tcp(%s:%d)/%s?parseTime=true",
        db.User,
        db.Password,
        db.Host,
        db.Port,
        db.Scheme)
}

type FcmConfig struct {
    Url       string
    ServerKey string
}

type StoreConfig struct {
    Apple  string
    Google string
}

type GoliathConfig struct {
    Server    ServerConfig
    Message   MessageConfig
    Memcached MemcachedConfig
    Database  []DatabaseConfig
    FCM       FcmConfig
    StoreUrl  StoreConfig
}

type configRoot struct {
    Goliath GoliathConfig
}

var root *configRoot = nil
var Values *GoliathConfig = nil

func Load(configPath string) (error) {
    buf, err := ioutil.ReadFile(os.ExpandEnv(configPath))
    if err != nil {
        return errors.WithStack(err)
    }

    err2 := yaml.Unmarshal(buf, &root)
    if err2 != nil {
        return errors.WithStack(err2)
    }
    Values = &root.Goliath
    return nil
}

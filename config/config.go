package config

import (
    "fmt"
    "github.com/pkg/errors"
    "github.com/shimalab-jp/goliath/util"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "os"
)

type VersionConfig struct {
    Version uint32
    Url     string
}

type ReferenceConfig struct {
    Enable      bool
    Url         string
    WebRoot     string
    Environment string
    Name        string
    Logo        string
    UserAgent   string
}

type DebugConfig struct {
    Enable   bool
    SlowTime uint64
    ClearDB  bool
}

type ServerConfig struct {
    Port      uint16
    IsFastCGI bool
    TimeZone  string
    UserDB    int8
    LogLevel  uint32
    Versions  []VersionConfig
    Reference ReferenceConfig
    Debug     DebugConfig
}

type ClientConfig struct {
    UserAgentPattern string
    MismatchAction   uint8
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

func (db *DatabaseConfig) ConnectionString() string {
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
    Client    ClientConfig
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

func Load(configPath string) error {
    path := os.ExpandEnv(configPath)

    if !util.FileExists(path) {
        return errors.New(fmt.Sprintf("Configuration file '%s' is not exists.", path))
    }

    buf, err := ioutil.ReadFile(path)
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

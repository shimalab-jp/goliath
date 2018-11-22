package rest

import (
    "encoding/json"
    "fmt"
    "github.com/bradfitz/gomemcache/memcache"
    "github.com/shimalab-jp/goliath/config"
)

type Memcached struct {
    instance *memcache.Client
}

func (mem *Memcached) open() {
    var serverList []string
    for _, def := range config.Values.Memcached.Servers {
        serverList = append(serverList, fmt.Sprintf("%s:%d", def.Host, def.Port))
    }
    mem.instance = memcache.New(serverList...)
}

func (mem *Memcached) createKey(key string) string {
    return fmt.Sprintf("%s:%s", config.Values.Memcached.Prefix, key)
}

func (mem *Memcached) Get(key string, value interface{}) error {
    if mem.instance == nil {
        mem.open()
    }
    if mem.instance == nil {
        value = nil
        return nil
    }

    item, err := mem.instance.Get(mem.createKey(key))
    if err == nil {
        _ = json.Unmarshal(item.Value, &value)
        return nil
    } else {
        return err
    }
}

func (mem *Memcached) Set(key string, value interface{}) error {
    if mem.instance == nil {
        mem.open()
    }
    if mem.instance == nil {
        return nil
    }

    buffer, err := json.Marshal(value)
    if err != nil {
        return err
    }

    item := memcache.Item{
        Key:        mem.createKey(key),
        Value:      buffer,
        Expiration: config.Values.Memcached.Expiration}

    return mem.instance.Set(&item)
}

func (mem *Memcached) Delete(key string) error {
    if mem.instance == nil {
        mem.open()
    }
    if mem.instance == nil {
        return nil
    }

    return mem.instance.Delete(mem.createKey(key))
}

func (mem *Memcached) Flush() error {
    if mem.instance == nil {
        mem.open()
    }
    if mem.instance == nil {
        return nil
    }

    return mem.instance.FlushAll()
}

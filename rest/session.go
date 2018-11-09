package rest

import (
    "fmt"
    "github.com/shimalab-jp/goliath/log"
    "github.com/shimalab-jp/goliath/util"
)

const (
    SessionKyeFormat1 = "SESSION1:%d"
    SessionKyeFormat2 = "SESSION2:%s"
)

type Session struct {
    SessionID string
    Account   accountInfo
}

type sessionManager struct {}

func createSessionManager() (*sessionManager) {
    return &sessionManager{}
}

func (sm *sessionManager) Open(info accountInfo) (*Session) {
    // memcachedのインスタンスを作成
    mem := Memcached{}

    // ユーザーIDでセッションキー1を作成
    key1 := fmt.Sprintf(SessionKyeFormat1, info.UserID)

    // 現在のセッションキーを取得
    var currentKey2 string
    {
        err := mem.Get(key1, &currentKey2)
        if err != nil {
            log.We(err)
        }
    }

    // 現在のセッションキーを削除
    if len(currentKey2) > 0 {
        mem.Delete(currentKey2)
    }

    // セッションIDを作成
    sessionID := util.GenerateUuid()

    // 新しいセッションを登録
    key2 := fmt.Sprintf(SessionKyeFormat2, sessionID)
    {
        err := mem.Set(key1, key2)
        if err != nil {
            log.We(err)
        }
    }
    {
        err := mem.Set(key2, info)
        if err != nil {
            log.We(err)
        }
    }

    // 結果を返す
    return &Session{
        SessionID: sessionID,
        Account: info}
}

func (sm *sessionManager) Get(sessionID string) (*Session) {
    // memcachedのインスタンスを作成
    mem := Memcached{}

    // セッションキーを作成
    key2 := fmt.Sprintf(SessionKyeFormat2, sessionID)

    // セッション情報を取得
    var session Session
    err := mem.Get(key2, &session)
    if err != nil {
        log.We(err)
        return nil
    }

    // 結果を返す
    return &session
}

func (sm *sessionManager) Close(info accountInfo) {
    // memcachedのインスタンスを作成
    mem := Memcached{}

    // ユーザーIDでセッションキー1を作成
    key1 := fmt.Sprintf(SessionKyeFormat1, info.UserID)

    // 現在のセッションキーを取得
    var currentKey2 string
    {
        err := mem.Get(key1, &currentKey2)
        if err != nil {
            log.We(err)
        }
    }

    // 現在のセッションキーを削除
    if len(currentKey2) > 0 {
        err := mem.Delete(currentKey2)
        if err != nil {
            log.We(err)
        }
    }

    if len(key1) > 0 {
        err := mem.Delete(key1)
        if err != nil {
            log.We(err)
        }
    }
}

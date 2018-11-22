package rest

import (
    "crypto/sha512"
    "fmt"
    "github.com/pkg/errors"
    "github.com/shimalab-jp/goliath/config"
    "github.com/shimalab-jp/goliath/database"
    "github.com/shimalab-jp/goliath/log"
    "github.com/shimalab-jp/goliath/message"
    "math/rand"
    "strings"
    "sync"
    "time"
)

const (
    accountCreatePlayerIdMaxRetry int    = 3
    accountPlayerIdFormat         string = "%04d-%04d"
    accountPasswordGenerateLength int    = 6
    accountPasswordTokens         string = "ABCDEFGHJKLMNPQRTUVWXYZ2346789"
    accountTokenFormat            string = "%04x%04x%04x:%04x%04x:%04x:%04x:%04x%04x:%04x%04x%04x"
    accountCreateTokenMaxRetry    int    = 3
    accountCreateMaxRetry         int    = 3
)

type AccountOutputInfo struct {
    PlayerID string
    Password string
    Platform int8
    Token    string
}


/********************
 *** Account Info ***
 ********************/
type accountInfo struct {
    UserID         int64
    PlayerID       string
    Password       string
    DatabaseNumber int8
    Platform       int8
    Token          string
    IsBan          bool
    IsAdmin        bool
}

func (am *accountInfo) Output() AccountOutputInfo {
    return AccountOutputInfo{
        PlayerID: am.PlayerID,
        Password: am.Password,
        Platform: am.Platform,
        Token:    am.Token}
}


/*********************
 *** Account Cache ***
 *********************/
type accountCache struct {
    tokenCacheMutex *sync.Mutex
    tokenCache      map[string]int64
}

func (ac *accountCache) getUserIdByTokenFromCache(token string) int64 {
    var ret int64 = 0

    ac.tokenCacheMutex.Lock()
    if userID, ok := ac.tokenCache[token]; ok {
        ret = userID
    }
    ac.tokenCacheMutex.Unlock()

    return ret
}

func (ac *accountCache) setTokenToCache(userID int64, token string) {
    if userID <= 0 {
        return
    }

    ac.tokenCacheMutex.Lock()
    ac.tokenCache[token] = userID
    ac.tokenCacheMutex.Unlock()
}

var cacheInstance *accountCache

func getAccountCache() *accountCache {
    if cacheInstance == nil {
        cacheInstance =         &accountCache{
            tokenCacheMutex:    &sync.Mutex{},
            tokenCache:         map[string]int64{}}
    }
    return cacheInstance
}


/***********************
 *** Account Manager ***
 ***********************/
type AccountManager struct {
    cache    *accountCache
    request  *Request
    response *Response
}

func (am *AccountManager) countTotalUsers(con *database.Connection) (int64, error) {
    result, err := con.Query("SELECT COUNT(*) AS `users` FROM `goliath_dat_account`;")
    if err != nil {
        return 0, err
    }

    if result.MoveFirst() {
        return result.GetInt64("users", 0), nil
    }

    return 0, nil
}

func (am *AccountManager) generatePlayerId(con *database.Connection) (string, error) {
    var playerID = ""
    var count = 0
    for {
        count++
        rand.Seed(time.Now().UnixNano())
        playerID = fmt.Sprintf(accountPlayerIdFormat, rand.Int31n(8999)+1000, rand.Int31n(9999))

        result, err := con.Query("SELECT COUNT(*) AS `users` FROM `goliath_dat_account` WHERE `player_id` = ?;", playerID)
        if err != nil {
            log.We(err)
            playerID = ""
        } else {
            if !result.MoveFirst() {
                playerID = ""
            } else if result.GetUInt64("users", 0) > 0 {
                playerID = ""
            }
        }

        if len(playerID) > 0 {
            break
        }

        if count >= accountCreatePlayerIdMaxRetry {
            return "", errors.New(message.SystemMessageManager.Get("SER_ACC_201", accountCreatePlayerIdMaxRetry))
        }
    }
    return playerID, nil
}

func (am *AccountManager) generatePassword() string {
    ret := ""
    rand.Seed(time.Now().UnixNano())
    runes := []rune(accountPasswordTokens)
    for {
        if len(ret) < accountPasswordGenerateLength {
            pos := rand.Int31n(int32(len(runes)))
            ret = ret + string(runes[pos])
        } else {
            break
        }
    }
    return ret
}

func (am *AccountManager) generateAccountToken(con *database.Connection) (string, error) {
    var token = ""
    var count = 0
    for {
        count++
        rand.Seed(time.Now().UnixNano())
        key := fmt.Sprintf(accountTokenFormat,
            rand.Int31n(0xFFFF), rand.Int31n(0xFFFF), rand.Int31n(0xFFFF),
            rand.Int31n(0xFFFF), rand.Int31n(0xFFFF),
            rand.Int31n(0x0FFF)|0x4000,
            rand.Int31n(0x3FFF)|0x8000,
            rand.Int31n(0xFFFF), rand.Int31n(0xFFFF),
            rand.Int31n(0xFFFF), rand.Int31n(0xFFFF), rand.Int31n(0xFFFF))

        token = strings.ToLower(fmt.Sprintf("%x", sha512.Sum512([]byte(key))))

        result, err := con.Query("SELECT COUNT(*) AS `users` FROM `goliath_dat_account_token` WHERE `token` = ?;", token)
        if err != nil {
            log.We(err)
            token = ""
        } else {
            if !result.MoveFirst() {
                token = ""
            } else if result.GetUInt64("users", 0) > 0 {
                token = ""
            }
        }

        if len(token) > 0 {
            break
        }

        if count >= accountCreateTokenMaxRetry {
            return "", errors.New(message.SystemMessageManager.Get("SER_ACC_221", accountCreateTokenMaxRetry))
        }
    }
    return token, nil
}

func (am *AccountManager) insertAccount(con *database.Connection, playerID string, password string, databaseNumber int8, platform int8, businessDay businessDay) (int64, error) {
    var sql = "INSERT INTO `goliath_dat_account` (`player_id`, `password`, `database_number`, `platform`, `is_ban`, `is_admin`, `create_business_day`, `last_access`) VALUES (?, ?, ?, ?, ?, ?, ?, ?);"
    var args = []interface{}{
        playerID,
        password,
        databaseNumber,
        platform,
        0,
        0,
        businessDay.BusinessDay,
        time.Now().UnixNano()}

    _, err := con.Execute(sql, args...)
    if err != nil {
        return 0, err
    }
    return con.LastInsertedId, nil
}

func (am *AccountManager) renewPassword(con *database.Connection, userID int64, password string) error {
    var sql = "UPDATE `goliath_dat_account` SET `password` = ? WHERE `user_id` = ?;"
    var args = []interface{}{password, userID}

    cnt, err := con.Execute(sql, args...)
    if err != nil {
        return err
    } else if cnt != 1 {
        return errors.New(fmt.Sprintf("Invalid affected count %d/1.", cnt))
    }
    return nil
}

func (am *AccountManager) insertAccountToken(con *database.Connection, token string, userID int64, businessDay businessDay) error {
    var sql = "INSERT INTO `goliath_dat_account_token` (`token`, `user_id`, `is_valid`, `regist_month`) VALUES (?, ?, ?, ?);"
    var args = []interface{}{
        token,
        userID,
        1,
        businessDay.Year*100 + businessDay.Month}

    _, err := con.Execute(sql, args...)
    if err != nil {
        return err
    }
    return nil
}

func (am *AccountManager) renewAccountToken(con *database.Connection, token string, userID int64, businessDay businessDay) error {
    {
        var sql = "UPDATE `goliath_dat_account_token` SET `is_valid` = 0 WHERE `user_id` = ? AND `is_valid` = 1;"
        _, err := con.Execute(sql, userID)
        if err != nil {
            return err
        }
    }

    {
        var sql = "INSERT INTO `goliath_dat_account_token` (`token`, `user_id`, `is_valid`, `regist_month`) VALUES (?, ?, 1, ?);"
        _, err := con.Execute(sql, token, userID, businessDay.Year*100 + businessDay.Month)
        if err != nil {
            return err
        }
    }

    return nil
}

func (am *AccountManager) insertAccountLog(con *database.Connection, userID int64, platform int8, businessDay businessDay) error {
    var sql = "INSERT INTO `goliath_log_create_account` (`business_day`, `user_id`, `platform`) VALUES (?, ?, ?);"
    var args = []interface{}{businessDay.BusinessDay, userID, platform}

    _, err := con.Execute(sql, args...)
    if err != nil {
        return err
    }
    return nil
}

func (am *AccountManager) getAccountByTokenFromMemcached(token string) *accountInfo {
    mem := Memcached{}
    key := fmt.Sprintf("AccountCache:%s", token)
    var account accountInfo
    err := mem.Get(key, &account)
    if err != nil {
        log.We(err)
        return nil
    } else if account.UserID <= 0 {
        return nil
    }
    return &account
}

func (am *AccountManager) setAccountToMemcached(account *accountInfo) {
    if account == nil {
        return
    }

    mem := Memcached{}
    key := fmt.Sprintf("AccountCache:%s", account.Token)
    err := mem.Set(key, account)
    if err != nil {
        log.We(err)
    }
}

func (am *AccountManager) getAccountByToken(con *database.Connection, token string) (*accountInfo, error) {
    // ローカルキャッシュからユーザーIDの取得を試みる
    userID := am.cache.getUserIdByTokenFromCache(token)

    // キャッシュから取得できなかった場合はDBから取得
    if userID == 0 {
        var sql = "SELECT `user_id` FROM `goliath_dat_account_token` WHERE `token` = ? AND is_valid = 1;"
        result, err := con.Query(sql, token)
        if err == nil {
            if result.MoveFirst() {
                userID = result.GetInt64("user_id", 0)
            }
        }

        if userID <= 0 && err == nil {
            return nil, err
        }

        if userID > 0 {
            am.cache.setTokenToCache(userID, token)
        }
    }

    if userID > 0 {
        account := accountInfo{UserID: userID, Token: token}

        var sql = "SELECT `player_id`, `password`, `database_number`, `platform`, `is_ban`, `is_admin` FROM `goliath_dat_account` WHERE `user_id` = ?;"
        result, err := con.Query(sql, userID)
        if err == nil {
            if result.MoveFirst() {
                account.PlayerID = result.GetString("player_id", "")
                account.Password = result.GetString("password", "")
                account.DatabaseNumber = result.GetInt8("database_number", 0)
                account.Platform = result.GetInt8("platform", 0)
                account.IsBan = result.GetBoolean("is_ban", false)
                account.IsAdmin = result.GetBoolean("is_admin", false)
            }
        }
    }

    return nil, nil
}

func (am *AccountManager) getAccountByPlayerID(con *database.Connection, playerID string) (*accountInfo, error) {
    var sql = "SELECT `user_id`, `player_id`, `password`, `database_number`, `platform`, `is_ban`, `is_admin` FROM `goliath_dat_account` WHERE `player_id` = ?;"
    result, err := con.Query(sql, playerID)
    if err == nil {
        if result.MoveFirst() {
            account := accountInfo{}
            account.UserID = result.GetInt64("user_id", 0)
            account.PlayerID = result.GetString("player_id", "")
            account.Password = result.GetString("password", "")
            account.DatabaseNumber = result.GetInt8("database_number", 0)
            account.Platform = result.GetInt8("platform", 0)
            account.IsBan = result.GetBoolean("is_ban", false)
            account.IsAdmin = result.GetBoolean("is_admin", false)
            return &account, nil
        }
    } else {
        return nil, err
    }

    return nil, nil
}

func (am *AccountManager) GetAccountByToken(token string) (*accountInfo, error) {
    account := am.getAccountByTokenFromMemcached(token)
    if account != nil {
        return account, nil
    }

    con, err := database.Connect("goliath")
    if err == nil {
        defer con.Disconnect()

        account, err := am.getAccountByToken(con, token)
        if err != nil {
            return nil, err
        }

        if account != nil {
            am.setAccountToMemcached(account)
        }

        return account, nil
    } else {
        return nil, err
    }
}

func (am *AccountManager) Create(platform int8) (*accountInfo, error) {
    var err error
    var con *database.Connection
    var retry = false
    var retryCount = 0
    var totalUsers int64 = 0
    var databaseNumber int8 = 0
    var playerID = ""
    var password = ""
    var token = ""
    var userID int64 = 0

    for {
        retry = false
        retryCount++

        // 接続
        if !retry {
            con, err = database.Connect("goliath")
            if err != nil {
                log.Ee(err)
                retry = true
            }
        }

        // 登録ユーザー数を取得
        if !retry {
            totalUsers, err = am.countTotalUsers(con)
            if err != nil {
                log.Ee(err)
                retry = true
            }
        }

        // DB番号を作成
        if !retry {
            if config.Values.Server.UserDB <= 0 {
                databaseNumber = 0
            } else {
                databaseNumber = int8(totalUsers % int64(config.Values.Server.UserDB))
            }
        }

        // プレイヤーIDを作成
        if !retry {
            playerID, err = am.generatePlayerId(con)
            if err != nil {
                log.Ee(err)
                retry = true
            }
        }

        // パスワードを作成
        if !retry {
            password = am.generatePassword()
        }

        // アカウントトークンを作成
        if !retry {
            token, err = am.generateAccountToken(con)
            if err != nil {
                log.Ee(err)
                retry = true
            }
        }

        // トランザクションを開始
        if !retry {
            err = con.BeginTransaction()
            if err != nil {
                log.Ee(err)
                retry = true
            }
        }

        // アカウント登録
        if !retry {
            userID, err = am.insertAccount(con, playerID, password, databaseNumber, platform, am.request.BusinessDay)
            if err != nil {
                log.Ee(err)
                retry = true
            }
        }

        // アカウントトークンを登録
        if !retry {
            err = am.insertAccountToken(con, token, userID, am.request.BusinessDay)
            if err != nil {
                log.Ee(err)
                retry = true
            }
        }

        // ログを記録
        if !retry {
            err = am.insertAccountLog(con, userID, platform, am.request.BusinessDay)
            if err != nil {
                log.Ee(err)
                retry = true
            }
        }

        //
        if !retry {
            _ = con.Commit()
            con.Disconnect()
            break
        } else {
            _ = con.Rollback()
            con.Disconnect()
        }

        if retryCount >= accountCreateMaxRetry {
            return nil, errors.New(message.SystemMessageManager.Get("SER_ACC_102", accountCreatePlayerIdMaxRetry))
        }
    }

    // アカウントオブジェクトを返却
    account := &accountInfo{
        UserID:         userID,
        PlayerID:       playerID,
        Password:       password,
        DatabaseNumber: databaseNumber,
        Platform:       platform,
        Token:          token,
        IsBan:          false,
        IsAdmin:        false}

    // キャッシュに反映
    if account != nil {
        am.setAccountToMemcached(account)
    }

    return account, nil
}

func (am *AccountManager) RenewPassword(token string) (*accountInfo, error) {
    var err error
    var con *database.Connection = nil
    var account *accountInfo = nil

    // アカウントを取得
    account, err = am.GetAccountByToken(token)
    if err != err {
        return nil, err
    }

    // パスワード生成
    newPassword := am.generatePassword()

    // 接続
    con, err = database.Connect("goliath")
    if err != nil {
        return nil, err
    }

    // 更新
    err = am.renewPassword(con, account.UserID, newPassword)
    if err != nil {
        return nil, err
    }

    // 切断
    con.Disconnect()

    // 反映
    account.Password = newPassword

    // キャッシュに反映
    if account != nil {
        am.setAccountToMemcached(account)
    }

    // アカウントオブジェクトを返却
    return account, nil
}

func (am *AccountManager) Trans(playerID string, password string, platform int8) (*accountInfo, error) {
    var err error
    var con *database.Connection = nil
    var account *accountInfo = nil

    // 接続
    con, err = database.Connect("goliath")
    if err != nil {
        return nil, err
    }

    // アカウントを取得
    account, err = am.getAccountByPlayerID(con, playerID)
    if err != nil {
        return nil, err
    }
    if account == nil {
        am.response.SetErrorMessage("SER_RES_401")
        return nil, nil
    }

    // パスワードチェック
    if account.Password != password {
        am.response.SetErrorMessage("SER_RES_401")
        return nil, nil
    }

    // BANチェック
    if account.IsBan {
        am.response.SetErrorMessage("ERR_RES_121")
    }

    // 新規トークン発行
    account.Token, err = am.generateAccountToken(con)

    // トランザクション開始
    err = con.BeginTransaction()
    if err != nil {
        return nil, err
    }

    // トークン更新
    err = am.renewAccountToken(con, account.Token, account.UserID, am.request.BusinessDay)
    if err != nil {
        _ = con.Rollback()
        return nil, err
    }

    // コミット
    err = con.Commit()
    if err != nil {
        _ = con.Rollback()
        return nil, err
    }

    // キャッシュに格納
    am.cache.setTokenToCache(account.UserID, account.Token)

    // 結果を返す
    return account, nil
}

func GetAccountManager(request *Request, response *Response) *AccountManager {
    return &AccountManager{cache: getAccountCache(), request: request, response: response}
}

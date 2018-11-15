package rest

import (
    "crypto/sha512"
    "fmt"
    "github.com/pkg/errors"
    "math/rand"
    "strings"
    "sync"
    "time"

    "github.com/shimalab-jp/goliath/config"
    "github.com/shimalab-jp/goliath/database"
    "github.com/shimalab-jp/goliath/log"
    "github.com/shimalab-jp/goliath/message"
)

const (
    createPlayerIdMaxRetry     int    = 3
    playerIdFormat             string = "%04d-%04d"
    passwordGenerateLength     int    = 6
    passwordTokens             string = "ABCDEFGHJKLMNPQRTUVWXYZ2346789"
    accountTokenFormat         string = "%04x%04x%04x:%04x%04x:%04x:%04x:%04x%04x:%04x%04x%04x"
    createAccountTokenMaxRetry int    = 3
    createAccountMaxRetry      int    = 3
)

type AccountOutputInfo struct {
    PlayerID string
    Password string
    Platform int8
    Token    string
}

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

func (am *accountInfo) Output() (AccountOutputInfo) {
    return AccountOutputInfo{
        PlayerID: am.PlayerID,
        Password: am.Password,
        Platform: am.Platform,
        Token:    am.Token}
}

type accountManager struct {
    tokenCacheMutex *sync.Mutex
    tokenCache      map[string]int64
}

var accountManagerInstance *accountManager = nil

func (am *accountManager) countTotalUsers(con *database.Connection) (int64, error) {
    result, err := con.Query("SELECT COUNT(*) AS `users` FROM `goliath_dat_account`;")
    if err != nil {
        return 0, err
    }
    defer result.Rows.Close()

    var users int64 = 0
    for result.Rows.Next() {
        err := result.Rows.Scan(&users)
        if err != nil {
            return 0, err
        }
        break
    }

    return users, nil
}

func (am *accountManager) generatePlayerId(con *database.Connection) (string, error) {
    var playerID = ""
    var count = 0
    for {
        count++
        rand.Seed(time.Now().UnixNano())
        playerID = fmt.Sprintf(playerIdFormat, rand.Int31n(8999)+1000, rand.Int31n(9999))

        result, err := con.Query("SELECT COUNT(*) AS `users` FROM `goliath_dat_account` WHERE `player_id` = ?;", playerID)
        if err != nil {
            log.We(err)
            playerID = ""
        } else {
            for result.Rows.Next() {
                var users = 0
                err := result.Rows.Scan(&users)
                if err != nil {
                    log.We(err)
                    playerID = ""
                }
                break
            }
            result.Rows.Close()
        }

        if len(playerID) > 0 {
            break
        }

        if count >= createPlayerIdMaxRetry {
            return "", errors.New(message.SystemMessageManager.Get("SER_ACC_201", createPlayerIdMaxRetry))
        }
    }
    return playerID, nil
}

func (am *accountManager) generatePassword() (string) {
    ret := ""
    rand.Seed(time.Now().UnixNano())
    runes := []rune(passwordTokens)
    for {
        if len(ret) < passwordGenerateLength {
            pos := rand.Int31n(int32(len(runes)))
            ret = ret + string(runes[pos])
        } else {
            break
        }
    }
    return ret
}

func (am *accountManager) generateAccountToken(con *database.Connection) (string, error) {
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
            for result.Rows.Next() {
                var users = 0
                err := result.Rows.Scan(&users)
                if err != nil {
                    log.We(err)
                    token = ""
                }
                break
            }
            result.Rows.Close()
        }

        if len(token) > 0 {
            break
        }

        if count >= createAccountTokenMaxRetry {
            return "", errors.New(message.SystemMessageManager.Get("SER_ACC_221", createAccountTokenMaxRetry))
        }
    }
    return token, nil
}

func (am *accountManager) insertAccount(con *database.Connection, playerID string, password string, databaseNumber int8, platform int8, businessDay businessDay) (int64, error) {
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

func (am *accountManager) renewPassword(con *database.Connection, userID int64, password string) (error) {
    var sql = "UPDATE `goliath_dat_account` SET `password` = ? WHERE `user_id` = ?;"
    var args = []interface{}{password, userID}

    _, err := con.Execute(sql, args...)
    if err != nil {
        return err
    }
    return nil
}

func (am *accountManager) insertAccountToken(con *database.Connection, token string, userID int64, businessDay businessDay) (error) {
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

func (am *accountManager) insertAccountLog(con *database.Connection, userID int64, platform int8, businessDay businessDay) (error) {
    var sql = "INSERT INTO `goliath_log_create_account` (`business_day`, `user_id`, `platform`) VALUES (?, ?, ?);"
    var args = []interface{}{businessDay.BusinessDay, userID, platform}

    _, err := con.Execute(sql, args...)
    if err != nil {
        return err
    }
    return nil
}

func (am *accountManager) getUserIdFromCache(token string) (int64) {
    var ret int64 = 0

    am.tokenCacheMutex.Lock()
    if userID, ok := am.tokenCache[token]; ok {
        ret = userID
    }
    am.tokenCacheMutex.Unlock()

    return ret
}

func (am *accountManager) setUserIdToCache(userID int64, token string) {
    if userID <= 0 {
        return
    }

    am.tokenCacheMutex.Lock()
    am.tokenCache[token] = userID
    am.tokenCacheMutex.Unlock()
}

func (am *accountManager) getAccountByTokenFromMemcached(token string) (*accountInfo) {
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

func (am *accountManager) setAccountToMemcached(account *accountInfo) {
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

func (am *accountManager) getAccountByToken(con *database.Connection, token string) (*accountInfo, error) {
    // ローカルキャッシュからユーザーIDの取得を試みる
    userID := am.getUserIdFromCache(token)

    // キャッシュから取得できなかった場合はDBから取得
    if userID == 0 {
        var sql = "SELECT `user_id` FROM `goliath_dat_account_token` WHERE `token` = ? AND is_valid = 1;"
        result, err := con.Query(sql, token)
        if err == nil {
            defer result.Rows.Close()
            for result.Rows.Next() {
                err := result.Rows.Scan(&userID)
                if err != nil {
                    return nil, err
                }
            }
        }

        if userID <= 0 && err == nil {
            return nil, err
        }

        if userID > 0 {
            am.setUserIdToCache(userID, token)
        }
    }

    if userID > 0 {
        account := accountInfo{UserID: userID, Token: token}

        var sql = "SELECT `player_id`, `password`, `database_number`, `platform`, `is_ban`, `is_admin` FROM `goliath_dat_account` WHERE `user_id` = ?;"
        result, err := con.Query(sql, userID)
        if err == nil {
            defer result.Rows.Close()
            for result.Rows.Next() {
                var isBan, isAdmin int8
                err := result.Rows.Scan(&account.PlayerID, &account.Password, &account.DatabaseNumber, &account.Platform, &isBan, &isAdmin)
                if err != nil {
                    return nil, err
                }
                account.IsBan = isBan != 0
                account.IsAdmin = isAdmin != 0
                return &account, nil
            }
        }
    }

    return nil, nil
}

func (am *accountManager) GetAccountByToken(token string) (*accountInfo, error) {
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

func (am *accountManager) Create(request *Request, platform int8) (*accountInfo, error) {
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
            userID, err = am.insertAccount(con, playerID, password, databaseNumber, platform, request.BusinessDay)
            if err != nil {
                log.Ee(err)
                retry = true
            }
        }

        // アカウントトークンを登録
        if !retry {
            err = am.insertAccountToken(con, token, userID, request.BusinessDay)
            if err != nil {
                log.Ee(err)
                retry = true
            }
        }

        // ログを記録
        if !retry {
            err = am.insertAccountLog(con, userID, platform, request.BusinessDay)
            if err != nil {
                log.Ee(err)
                retry = true
            }
        }

        //
        if !retry {
            con.Commit()
            con.Disconnect()
            break
        } else {
            con.Rollback()
            con.Disconnect()
        }

        if retryCount >= createAccountMaxRetry {
            return nil, errors.New(message.SystemMessageManager.Get("SER_ACC_102", createPlayerIdMaxRetry))
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

func (am *accountManager) RenewPassword(token string) (*accountInfo, error) {
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

func GetAccountManager() (*accountManager) {
    if accountManagerInstance == nil {
        accountManagerInstance = &accountManager{
            tokenCacheMutex: &sync.Mutex{},
            tokenCache:      map[string]int64{}}
    }
    return accountManagerInstance
}

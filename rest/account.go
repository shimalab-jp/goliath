package rest

import (
    "fmt"
    "github.com/pkg/errors"
    "math/rand"
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

type accountInfo struct {
    UserID         int64
    PlayerID       string
    Password       string
    DatabaseNumber int8
    Platform       int8
    Token          string
    IsBan          bool
    isAdmin        bool
}

type accountManager struct{}

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
        token = fmt.Sprintf(accountTokenFormat,
            rand.Int31n(0xFFFF), rand.Int31n(0xFFFF), rand.Int31n(0xFFFF),
            rand.Int31n(0xFFFF), rand.Int31n(0xFFFF),
            rand.Int31n(0x0FFF)|0x4000,
            rand.Int31n(0x3FFF)|0x8000,
            rand.Int31n(0xFFFF), rand.Int31n(0xFFFF),
            rand.Int31n(0xFFFF), rand.Int31n(0xFFFF), rand.Int31n(0xFFFF))

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
    var args = []interface{}{ businessDay.BusinessDay, userID, platform }

    _, err := con.Execute(sql, args...)
    if err != nil {
        return err
    }
    return nil
}

func (am *accountManager) getAccountByToken(con *database.Connection, token string) (*accountInfo, error) {
    var userID int64 = 0

    {
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
    }

    if userID > 0 {
        account := accountInfo{ UserID: userID, Token: token }

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
                account.isAdmin = isAdmin != 0
                return &account, nil
            }
        }
    }

    return nil, nil
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
    return &accountInfo{
        UserID:         userID,
        PlayerID:       playerID,
        Password:       password,
        DatabaseNumber: databaseNumber,
        Platform:       platform,
        Token:          token,
        IsBan:          false,
        isAdmin:        false}, nil
}

func (am *accountManager) GetAccountByToken(token string) (*accountInfo, error) {
    con, err := database.Connect("goliath")
    if err == nil {
        defer con.Disconnect()
        return am.getAccountByToken(con, token)
    } else {
        return nil, err
    }
}

func GetAccountManager() (*accountManager) {
    return &accountManager{}
}
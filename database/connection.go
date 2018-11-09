package database

import (
    "database/sql"
    "fmt"
    "github.com/pkg/errors"
    "runtime"
    "strings"

    _ "github.com/go-sql-driver/mysql"
    "github.com/shimalab-jp/goliath/config"
)

type Connection struct {
    Name string
    Instance *sql.DB
    Transaction *sql.Tx
    LastInsertedId int64
}

type ResultSet struct {
    Rows *sql.Rows
}

func getConnectionInfo(databaseName string) (*config.DatabaseConfig) {
    if config.Values != nil {
        for _, database := range config.Values.Database {
            if strings.ToLower(database.Name) == strings.ToLower(databaseName) {
                return &database
            }
        }
    }
    return nil
}

func destroyConnection(con *Connection) {
    if con != nil {
        con.Disconnect()
    }
    con = nil
}

func destroyResult(set *ResultSet) {
    if set != nil || set.Rows != nil {
        set.Rows.Close()
        set.Rows = nil
    }
    set = nil
}

func (con *Connection) BeginTransaction() (error)  {
    if con == nil || con.Instance == nil {
        return errors.New(" database connection is nil.")
    }

    tx, err := con.Instance.Begin()
    if err != nil {
        return errors.WithStack(err)
    }
    con.Transaction = tx
    return nil
}

func (con *Connection) Commit() (error) {
    if con == nil || con.Instance == nil {
        return errors.New(" database connection is nil.")
    }
    if con.Transaction != nil {
        err := con.Transaction.Commit()
        con.Transaction = nil
        return errors.WithStack(err)
    }
    return nil
}

func (con *Connection) Rollback() (error) {
    if con == nil || con.Instance == nil {
        return errors.New(" database connection is nil.")
    }
    if con.Transaction != nil {
        err := con.Transaction.Rollback()
        con.Transaction = nil
        return errors.WithStack(err)
    }
    return nil
}

func (con *Connection) Disconnect() {
    if con != nil && con.Transaction != nil {
        con.Rollback()
    }
    if con != nil && con.Instance != nil {
        con.Instance.Close()
        con.Instance = nil
    }
}

func (con *Connection) Query(query string, args ...interface{}) (*ResultSet, error) {
    if con == nil || con.Instance == nil {
        return nil, errors.New(" database connection is nil.")
    }

    result := &ResultSet{}
    runtime.SetFinalizer(result, destroyResult)

    var rows *sql.Rows
    var err error

    if len(args) > 0 {
        var stmt *sql.Stmt

        if con != nil && con.Transaction != nil {
            stmt, err = con.Transaction.Prepare(query)
            if err != nil {
                return nil, errors.WithStack(err)
            }
        } else if con != nil && con.Instance != nil {
            stmt, err = con.Instance.Prepare(query)
            if err != nil {
                return nil, errors.WithStack(err)
            }
        }

        rows, err = stmt.Query(args...)
    } else {
        if con != nil && con.Transaction != nil {
            rows, err = con.Transaction.Query(query)
        } else if con != nil && con.Instance != nil {
            rows, err = con.Instance.Query(query)
        }
    }

    if err != nil {
        return nil, errors.WithStack(err)
    }
    result.Rows = rows

    return result, nil
}

func (con *Connection) Execute(query string, args ...interface{}) (int64, error) {
    if con == nil || con.Instance == nil {
        return 0, errors.New(" database connection is nil.")
    }

    if len(args) > 0 {
        var stmt *sql.Stmt
        var err error

        if con != nil && con.Transaction != nil {
            stmt, err = con.Transaction.Prepare(query)
            if err != nil {
                return 0, errors.WithStack(err)
            }
        } else if con != nil && con.Instance != nil {
            stmt, err = con.Instance.Prepare(query)
            if err != nil {
                return 0, errors.WithStack(err)
            }
        }

        result, err := stmt.Exec(args...)
        if err != nil {
            return 0, errors.WithStack(err)
        }

        affected, err := result.RowsAffected()
        if err != nil {
            return 0, errors.WithStack(err)
        }

        con.LastInsertedId, err = result.LastInsertId()
        if err != nil {
            return 0, errors.WithStack(err)
        }

        return affected, nil
    } else {
        var result sql.Result
        var err error

        if con != nil && con.Transaction != nil {
            result, err = con.Transaction.Exec(query)
            if err != nil {
                return 0, errors.WithStack(err)
            }
        } else if con != nil && con.Instance != nil {
            result, err = con.Instance.Exec(query)
            if err != nil {
                return 0, errors.WithStack(err)
            }
        }

        affected, err := result.RowsAffected()
        if err != nil {
            return 0, errors.WithStack(err)
        }

        con.LastInsertedId, err = result.LastInsertId()
        if err != nil {
            return 0, errors.WithStack(err)
        }

        return affected, nil
    }
}

func Connect(databaseName string) (*Connection, error) {
    // 指定された名前のDB接続情報を取得
    connectionInfo := getConnectionInfo(databaseName)
    if connectionInfo == nil {
        return nil, errors.New(fmt.Sprintf("'%s' database define not found.", databaseName))
    }

    // 接続
    db, err := sql.Open(connectionInfo.Driver, connectionInfo.ConnectionString())
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // コネクションオブジェクトを作成
    result := &Connection{
        Name: connectionInfo.Name,
        Instance: db,
        Transaction:nil,
        LastInsertedId: 0 }
    runtime.SetFinalizer(result, destroyConnection)

    return result, nil
}

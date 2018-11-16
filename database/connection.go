package database

import (
    "database/sql"
    "fmt"
    "github.com/pkg/errors"
    "github.com/shimalab-jp/goliath/log"
    "reflect"
    "runtime"
    "strings"

    _ "github.com/go-sql-driver/mysql"
    "github.com/shimalab-jp/goliath/config"
)

type Connection struct {
    Name           string
    Instance       *sql.DB
    Transaction    *sql.Tx
    LastInsertedId int64
}

type ResultSet struct {
    cols     []string
    colTypes []reflect.Kind
    rowCount int
    rowIndex int
    dataSet  []map[string]interface{}
}

func createResultSet(rows *sql.Rows) (*ResultSet) {
    rs := ResultSet{}
    if rows != nil {
        // カラムの一覧を取得
        {
            cols, err := rows.Columns()
            if err != nil {
                log.Ee(err)
                return nil
            }
            rs.cols = cols
        }

        // カラムの型の一覧を取得
        colTypes, err := rows.ColumnTypes()
        if err != nil {
            log.Ee(err)
            return nil
        }

        // 全値を取得
        rs.rowCount = 0
        for rows.Next() {
            // rowの格納配列を作成
            var temp []interface{}
            for _, t := range colTypes {
                switch t.ScanType().Kind() {
                case reflect.Bool:
                    temp = append(temp, false)
                    break
                case reflect.Int:
                    v := int(0)
                    temp = append(temp, &v)
                    break
                case reflect.Int8:
                    v := int8(0)
                    temp = append(temp, &v)
                    break
                case reflect.Int16:
                    v := int16(0)
                    temp = append(temp, &v)
                    break
                case reflect.Int32:
                    v := int32(0)
                    temp = append(temp, &v)
                    break
                case reflect.Int64:
                    v := int64(0)
                    temp = append(temp, &v)
                    break
                case reflect.Uint:
                    v := uint(0)
                    temp = append(temp, &v)
                    break
                case reflect.Uint8:
                    v := uint8(0)
                    temp = append(temp, &v)
                    break
                case reflect.Uint16:
                    v := uint16(0)
                    temp = append(temp, &v)
                    break
                case reflect.Uint32:
                    v := uint32(0)
                    temp = append(temp, &v)
                    break
                case reflect.Uint64:
                    v := uint64(0)
                    temp = append(temp, &v)
                    break
                case reflect.Float32:
                    v := float32(0)
                    temp = append(temp, &v)
                    break
                case reflect.Float64:
                    v := float64(0)
                    temp = append(temp, &v)
                    break
                case reflect.String:
                    v := ""
                    temp = append(temp, &v)
                    break
                default:
                    v := ""
                    temp = append(temp, &v)
                    break
                }
            }

            // 1行読取
            err := rows.Scan(temp...)
            if err != nil {
                log.Ee(err)
                return nil
            }

            row := map[string]interface{}{}
            i := 0
            for _, t := range colTypes {
                row[t.Name()] = temp[i]
                i++
            }
            rs.dataSet = append(rs.dataSet, row)
            rs.rowCount++
        }
    }

    return &rs
}

func (rs *ResultSet) EoF() (bool) {
    return rs.rowIndex >= len(rs.dataSet)
}

func (rs *ResultSet) MoveFirst() (bool) {
    rs.rowIndex = 0
    return !rs.EoF()
}

func (rs *ResultSet) MoveNext() (bool) {
    if rs.rowIndex < len(rs.dataSet) - 1 {
        rs.rowIndex++
        return true
    } else {
        return false
    }
}

func (rs *ResultSet) getValue(name string) (interface{}, bool) {
    if int(rs.rowIndex) < len(rs.dataSet) {
        if val, ok := rs.dataSet[rs.rowIndex][name]; ok {
            return val, true
        } else {
            return nil, false
        }
    } else {
        return nil, false
    }
}

func (rs *ResultSet) GetString(name string, defaultValue string) (string) {
    if v, ok := rs.getValue(name); ok {
        switch v.(type) {
        case string:
            return v.(string)
        default:
            return fmt.Sprintf("%x", v)
        }
    }
    return defaultValue
}

func (rs *ResultSet) GetInt(name string, defaultValue int) (int) {
    if v, ok := rs.getValue(name); ok {
        switch v.(type) {
        case int:
        case int8:
        case int16:
        case int32:
        case int64:
            return v.(int)
        case uint:
        case uint8:
        case uint16:
        case uint32:
        case uint64:
        case float32:
        case float64:
            return v.(int)
        }
    }
    return defaultValue
}

func (rs *ResultSet) GetInt8(name string, defaultValue int8) (int8) {
    if v, ok := rs.getValue(name); ok {
        switch v.(type) {
        case int:
        case int8:
            return v.(int8)
        case uint:
        case int16:
        case int32:
        case int64:
        case uint8:
        case uint16:
        case uint32:
        case uint64:
        case float32:
        case float64:
            return v.(int8)
        }
    }
    return defaultValue
}

func (rs *ResultSet) GetInt16(name string, defaultValue int16) (int16) {
    if v, ok := rs.getValue(name); ok {
        switch v.(type) {
        case int:
        case int16:
            return v.(int16)
        case uint:
        case int8:
        case int32:
        case int64:
        case uint8:
        case uint16:
        case uint32:
        case uint64:
        case float32:
        case float64:
            return v.(int16)
        }
    }
    return defaultValue
}

func (rs *ResultSet) GetInt32(name string, defaultValue int32) (int32) {
    if v, ok := rs.getValue(name); ok {
        switch v.(type) {
        case int:
        case int32:
            return v.(int32)
        case uint:
        case int8:
        case int16:
        case int64:
        case uint8:
        case uint16:
        case uint32:
        case uint64:
        case float32:
        case float64:
            return v.(int32)
        }
    }
    return defaultValue
}

func (rs *ResultSet) GetInt64(name string, defaultValue int64) (int64) {
    if v, ok := rs.getValue(name); ok {
        switch v.(type) {
        case int:
        case int64:
            return v.(int64)
        case uint:
        case int8:
        case int16:
        case int32:
        case uint8:
        case uint16:
        case uint32:
        case uint64:
        case float32:
        case float64:
            return v.(int64)
        }
    }
    return defaultValue
}

func (rs *ResultSet) GetUInt(name string, defaultValue uint) (uint) {
    if v, ok := rs.getValue(name); ok {
        switch v.(type) {
        case int:
        case int8:
        case int16:
        case int32:
        case int64:
            return v.(uint)
        case uint:
        case uint8:
        case uint16:
        case uint32:
        case uint64:
        case float32:
        case float64:
            return v.(uint)
        }
    }
    return defaultValue
}

func (rs *ResultSet) GetUInt8(name string, defaultValue uint8) (uint8) {
    if v, ok := rs.getValue(name); ok {
        switch v.(type) {
        case uint:
        case uint8:
            return v.(uint8)
        case int:
        case int8:
        case int16:
        case int32:
        case int64:
        case uint16:
        case uint32:
        case uint64:
        case float32:
        case float64:
            return v.(uint8)
        }
    }
    return defaultValue
}

func (rs *ResultSet) GetUInt16(name string, defaultValue uint16) (uint16) {
    if v, ok := rs.getValue(name); ok {
        switch v.(type) {
        case uint:
        case uint16:
            return v.(uint16)
        case int:
        case int8:
        case int16:
        case int32:
        case int64:
        case uint8:
        case uint32:
        case uint64:
        case float32:
        case float64:
            return v.(uint16)
        }
    }
    return defaultValue
}

func (rs *ResultSet) GetUInt32(name string, defaultValue uint32) (uint32) {
    if v, ok := rs.getValue(name); ok {
        switch v.(type) {
        case uint:
        case uint32:
            return v.(uint32)
        case int:
        case int8:
        case int16:
        case int32:
        case int64:
        case uint8:
        case uint16:
        case uint64:
        case float32:
        case float64:
            return v.(uint32)
        }
    }
    return defaultValue
}

func (rs *ResultSet) GetUInt64(name string, defaultValue uint64) (uint64) {
    if v, ok := rs.getValue(name); ok {
        switch v.(type) {
        case uint:
        case uint64:
            return v.(uint64)
        case int:
        case int8:
        case int16:
        case int32:
        case int64:
        case uint8:
        case uint16:
        case uint32:
        case float32:
        case float64:
            return v.(uint64)
        }
    }
    return defaultValue
}

func (rs *ResultSet) GetFloat32(name string, defaultValue float32) (float32) {
    if v, ok := rs.getValue(name); ok {
        switch v.(type) {
        case float64:
        case float32:
            return v.(float32)
        case int:
        case int8:
        case int16:
        case int32:
        case int64:
        case uint:
        case uint8:
        case uint16:
        case uint32:
        case uint64:
            return v.(float32)
        }
    }
    return defaultValue
}

func (rs *ResultSet) GetFloat64(name string, defaultValue float64) (float64) {
    if v, ok := rs.getValue(name); ok {
        switch v.(type) {
        case float32:
        case float64:
            return v.(float64)
        case int:
        case int8:
        case int16:
        case int32:
        case int64:
        case uint:
        case uint8:
        case uint16:
        case uint32:
        case uint64:
            return v.(float64)
        }
    }
    return defaultValue
}

func (rs *ResultSet) GetBoolean(name string, defaultValue bool) (bool) {
    if v, ok := rs.getValue(name); ok {
        switch v.(type) {
        case bool:
            return v.(bool)
        case int:
        case int8:
        case int16:
        case int32:
        case int64:
            return v.(int) != 0
        case uint:
        case uint8:
        case uint16:
        case uint32:
        case uint64:
            return v.(uint) != 0
        }
    }
    return defaultValue
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

func (con *Connection) BeginTransaction() (error) {
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

    defer rows.Close()

    result := createResultSet(rows)

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
        Name:           connectionInfo.Name,
        Instance:       db,
        Transaction:    nil,
        LastInsertedId: 0}
    runtime.SetFinalizer(result, destroyConnection)

    return result, nil
}

package goliath

import (
    "encoding/json"
    "fmt"
    "github.com/shimalab-jp/goliath/config"
    "github.com/shimalab-jp/goliath/database"
    "github.com/shimalab-jp/goliath/log"
    "github.com/shimalab-jp/goliath/message"
    "github.com/shimalab-jp/goliath/rest"
    "github.com/shimalab-jp/goliath/util"
    "io/ioutil"
    "mime"
    "net"
    "net/http"
    "net/http/fcgi"
    "os"
    "path"
    "strings"
)

const (
    DataVersion    string = "1.0.0"
    DataVersionNum uint32 = 1
)

func initDatabase() error {
    var err error
    var con *database.Connection

    con, err = database.Connect("goliath")
    if err != nil { return err }

    err = con.BeginTransaction()
    if err != nil { return err }

    if config.Values.Server.Debug.Enable && config.Values.Server.Debug.ClearDB {
        if err == nil {
            _, err = con.Execute("DROP TABLE IF EXISTS `goliath_mst_version`;")
            if err == nil {
                log.I("[INITDB] DROP `goliath_mst_version` table.")
            } else {
                log.E("[INITDB] DROP FAILED `goliath_mst_version` table.")
            }
        }
        if err == nil {
            _, err = con.Execute("DROP TABLE IF EXISTS `goliath_dat_account`;")
            if err == nil {
                log.I("[INITDB] DROP `goliath_dat_account` table.")
            } else {
                log.E("[INITDB] DROP FAILED `goliath_dat_account` table.")
            }
        }
        if err == nil {
            _, err = con.Execute("DROP TABLE IF EXISTS `goliath_dat_account_token`;")
            if err == nil {
                log.I("[INITDB] DROP `goliath_dat_account_token` table.")
            } else {
                log.E("[INITDB] DROP FAILED `goliath_dat_account_token` table.")
            }
        }
        if err == nil {
            _, err = con.Execute("DROP TABLE IF EXISTS `goliath_log_install`;")
            if err == nil {
                log.I("[INITDB] DROP `goliath_log_install` table.")
            } else {
                log.E("[INITDB] DROP FAILED `goliath_log_install` table.")
            }
        }
        if err == nil {
            _, err = con.Execute("DROP TABLE IF EXISTS `goliath_log_create_account`;")
            if err == nil {
                log.I("[INITDB] DROP `goliath_log_create_account` table.")
            } else {
                log.E("[INITDB] DROP FAILED `goliath_log_create_account` table.")
            }
        }
        if err == nil {
            _, err = con.Execute("DROP TABLE IF EXISTS `goliath_log_hau`;")
            if err == nil {
                log.I("[INITDB] DROP `goliath_log_hau` table.")
            } else {
                log.E("[INITDB] DROP FAILED `goliath_log_hau` table.")
            }
        }
        if err == nil {
            _, err = con.Execute("DROP TABLE IF EXISTS `goliath_mst_api_switch`;")
            if err == nil {
                log.I("[INITDB] DROP `goliath_mst_api_switch` table.")
            } else {
                log.E("[INITDB] DROP FAILED `goliath_mst_api_switch` table.")
            }
        }
        if err == nil {
            _, err = con.Execute("DROP TABLE IF EXISTS `goliath_mst_client_version`;")
            if err == nil {
                log.I("[INITDB] DROP `goliath_mst_client_version` table.")
            } else {
                log.E("[INITDB] DROP FAILED `goliath_mst_client_version` table.")
            }
        }
        if err == nil {
            _, err = con.Execute("DROP TABLE IF EXISTS `goliath_mst_maintenance`;")
            if err == nil {
                log.I("[INITDB] DROP `goliath_mst_maintenance` table.")
            } else {
                log.E("[INITDB] DROP FAILED `goliath_mst_maintenance` table.")
            }
        }

        if err == nil {
            mem := rest.Memcached{}
            err = mem.Flush()
        }
    }

    if err == nil {
        _, err = con.Execute("CREATE TABLE IF NOT EXISTS `goliath_mst_version` (" +
            "`version_num` int(10) unsigned NOT NULL, " +
            "`version` varchar(32) CHARACTER SET ascii NOT NULL, " +
            "`created` datetime NOT NULL DEFAULT current_timestamp(), " +
            "`modified` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(), " +
            "PRIMARY KEY (`version_num`)" +
            ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
        if err == nil {
            log.I("[INITDB] CREATE `goliath_mst_version` table.")
            _, err = con.Execute(
                "REPLACE INTO `goliath_mst_version` (`version_num`, `version`) VALUES (?, ?);",
                DataVersionNum, DataVersion)
            if err == nil {
                log.I("[INITDB] INSERT version data.")
            } else {
                log.E("[INITDB] INSERT FAILED version data.")
            }
        } else {
            log.E("[INITDB] CREATE FAILED `goliath_mst_version` table.")
        }
    }

    if err == nil {
        _, err = con.Execute("CREATE TABLE IF NOT EXISTS `goliath_dat_account` (" +
            "`user_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT, " +
            "`player_id` varchar(32) NOT NULL DEFAULT '', " +
            "`password` varchar(32) NOT NULL DEFAULT '', " +
            "`database_number` tinyint(2) NOT NULL DEFAULT 0, " +
            "`platform` tinyint(1) unsigned NOT NULL DEFAULT 0, "  +
            "`is_ban` tinyint(1) unsigned NOT NULL DEFAULT 0, " +
            "`is_admin` tinyint(1) unsigned NOT NULL DEFAULT 0, " +
            "`create_business_day` int(10) unsigned NOT NULL DEFAULT 0, " +
            "`last_access` bigint(20) unsigned NOT NULL DEFAULT 0, " +
            "`created` datetime NOT NULL DEFAULT current_timestamp(), " +
            "`modified` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(), " +
            "PRIMARY KEY (`user_id`), " +
            "UNIQUE KEY `IDX_PLAYER_ID` (`player_id`), " +
            "KEY `IDX_CREATE_BIZ_DATE` (`create_business_day`), " +
            "KEY `IDX_COUNT1` (`is_ban`), " +
            "KEY `IDX_COUNT2` (`is_ban`, `is_admin`), " +
            "KEY `IDX_KPI_FLASH` (`platform`, `created`), " +
            "KEY `IDX_JOIN` (`is_ban`, `user_id`)" +
            ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
        if err == nil {
            log.I("[INITDB] CREATE `goliath_dat_account` table.")
        } else {
            log.E("[INITDB] CREATE FAILED `goliath_dat_account` table.")
        }
    }

    if err == nil {
        _, err = con.Execute("CREATE TABLE IF NOT EXISTS `goliath_dat_account_token` (" +
            "`token` varchar(512) CHARACTER SET ascii NOT NULL," +
            "`user_id` bigint(20) unsigned NOT NULL," +
            "`is_valid` tinyint(1) unsigned NOT NULL," +
            "`regist_month` int(10) unsigned NOT NULL COMMENT 'yyyymm'," +
            "`created` datetime NOT NULL DEFAULT current_timestamp()," +
            "`modified` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()," +
            "PRIMARY KEY (`token`)," +
            "KEY `IDX_COUNT1` (`user_id`, `regist_month`)," +
            "KEY `IDX_TRANSIT` (`user_id`, `is_valid`)" +
            ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
        if err == nil {
            log.I("[INITDB] CREATE `goliath_dat_account_token` table.")
        } else {
            log.E("[INITDB] CREATE FAILED `goliath_dat_account_token` table.")
        }
    }

    if err == nil {
        _, err = con.Execute("CREATE TABLE IF NOT EXISTS `goliath_log_create_account` (" +
            "`business_day` int(10) unsigned NOT NULL, " +
            "`user_id` bigint(20) unsigned NOT NULL, " +
            "`platform` tinyint(1) unsigned NOT NULL, " +
            "`created` datetime NOT NULL DEFAULT current_timestamp(), " +
            "`modified` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(), " +
            "PRIMARY KEY (`user_id`), " +
            "KEY `IDX_COUNT1` (`business_day`, `platform`)" +
            ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
        if err == nil {
            log.I("[INITDB] CREATE `goliath_log_create_account` table.")
        } else {
            log.E("[INITDB] CREATE FAILED `goliath_log_create_account` table.")
        }
    }

    if err == nil {
        _, err = con.Execute("CREATE TABLE IF NOT EXISTS `goliath_log_hau` (" +
            "`access_date` int(10) unsigned NOT NULL, " +
            "`access_hour` int(10) unsigned NOT NULL, " +
            "`user_id` bigint(20) unsigned NOT NULL, " +
            "`platform` tinyint(1) unsigned NOT NULL, " +
            "PRIMARY KEY (`access_date`, `access_hour`, `user_id`), " +
            "KEY `IDX_COUNT1` (`platform`)" +
            ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
        if err == nil {
            log.I("[INITDB] CREATE `goliath_log_hau` table.")
        } else {
            log.E("[INITDB] CREATE FAILED `goliath_log_hau` table.")
        }
    }

    if err == nil {
        _, err = con.Execute("CREATE TABLE IF NOT EXISTS `goliath_mst_api_switch` (" +
            "`api_name` varchar(256) CHARACTER SET ascii NOT NULL, " +
            "`enable` tinyint(1) unsigned NOT NULL DEFAULT 1, " +
            "`created` datetime NOT NULL DEFAULT current_timestamp(), " +
            "`modified` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(), " +
            "PRIMARY KEY (`api_name`) " +
            ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
        if err == nil {
            log.I("[INITDB] CREATE `goliath_mst_api_switch` table.")
        } else {
            log.E("[INITDB] CREATE FAILED `goliath_mst_api_switch` table.")
        }
    }

    if err == nil {
        _, err = con.Execute("CREATE TABLE IF NOT EXISTS `goliath_mst_client_version` (" +
            "`id` int(10) unsigned NOT NULL AUTO_INCREMENT, " +
            "`platform` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT '1:iOS\n2:Android', " +
            "`major` int(10) unsigned NOT NULL DEFAULT 0, " +
            "`minor` int(10) unsigned NOT NULL DEFAULT 0, " +
            "`revision` int(10) unsigned NOT NULL DEFAULT 0, " +
            "`resource_version` varchar(16) NOT NULL DEFAULT '', " +
            "`start_time` bigint(20) unsigned NOT NULL, " +
            "`end_time` bigint(20) unsigned NOT NULL, " +
            "`description` text NOT NULL, " +
            "`created` datetime NOT NULL DEFAULT current_timestamp()," +
            "`modified` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()," +
            "PRIMARY KEY (`id`)," +
            "KEY `IDX_SEARCH1` (`platform`,`end_time`)" +
            ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
        if err == nil {
            log.I("[INITDB] CREATE `goliath_mst_client_version` table.")
        } else {
            log.E("[INITDB] CREATE FAILED `goliath_mst_client_version` table.")
        }

        if err == nil {
            _, _ = con.Execute("REPLACE INTO `goliath_mst_client_version` (`platform`, `major`, `minor`, `revision`, `resource_version`, `start_time`, `end_time`, `description`) VALUES (1, 1, 0, 0, '', 0, 253402268399, 'Init insert');")
            _, _ = con.Execute("REPLACE INTO `goliath_mst_client_version` (`platform`, `major`, `minor`, `revision`, `resource_version`, `start_time`, `end_time`, `description`) VALUES (2, 1, 0, 0, '', 0, 253402268399, 'Init insert');")
            _, _ = con.Execute("REPLACE INTO `goliath_mst_client_version` (`platform`, `major`, `minor`, `revision`, `resource_version`, `start_time`, `end_time`, `description`) VALUES (1, 1, 0, 1, '', 0, 253402268399, 'Init insert');")
            _, _ = con.Execute("REPLACE INTO `goliath_mst_client_version` (`platform`, `major`, `minor`, `revision`, `resource_version`, `start_time`, `end_time`, `description`) VALUES (2, 1, 0, 1, '', 0, 253402268399, 'Init insert');")
        }
    }

    if err == nil {
        _, err = con.Execute("CREATE TABLE IF NOT EXISTS `goliath_mst_maintenance` (" +
            "`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT, " +
            "`start_time` bigint(20) NOT NULL DEFAULT 0, " +
            "`end_time` bigint(20) NOT NULL DEFAULT 0, " +
            "`subject` varchar(256) NOT NULL DEFAULT '', " +
            "`body` text NOT NULL, " +
            "`admin_id` varchar(32) NOT NULL DEFAULT '', " +
            "`created` datetime NOT NULL DEFAULT current_timestamp(), " +
            "`modified` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(), " +
            "PRIMARY KEY (`id`), " +
            "KEY `IDX_SEARCH1` (`start_time`, `end_time`)" +
            ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
        if err == nil {
            log.I("[INITDB] CREATE `goliath_mst_maintenance` table.")
        } else {
            log.E("[INITDB] CREATE FAILED `goliath_mst_maintenance` table.")
        }
    }

    if err == nil {
        err = con.Commit()
    } else {
        _ = con.Rollback()
    }

    con.Disconnect()
    con = nil

    return err
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
    err := rest.GetEngine().Execute(r, w)
    if err != nil {
        log.Ee(err)
    }
}

func referenceHandler(w http.ResponseWriter, r *http.Request) {
    accessPath := r.URL.Path[len(config.Values.Server.Reference.Url):]
    if len(accessPath) <= 0 {
        accessPath = "index.html"
    }

    webroot := strings.TrimRight(config.Values.Server.Reference.WebRoot, "/")
    if !util.DirectoryExists(webroot) {
        log.E("[Reference] Web root directory '%s' is not exists.", webroot)
        w.WriteHeader(http.StatusForbidden)
        return
    }
    vPath := fmt.Sprintf("%s/%s", webroot, accessPath)
    realPath := os.ExpandEnv(vPath)

    if accessPath == "config.json" {
        // リファレンス情報を取得
        envName := config.Values.Server.Reference.Environment
        envCode := util.ToEnvironmentCode(envName)
        className := ""
        switch envCode {
        case util.EnvironmentDemo:
            envName = "DEMO"
            className = "env_name_demo"
            break
        case util.EnvironmentDevelop1:
            envName = "DEVELOP1"
            className = "env_name_dev1"
            break
        case util.EnvironmentDevelop2:
            envName = "DEVELOP2"
            className = "env_name_dev2"
            break
        case util.EnvironmentTest:
            envName = "TEST"
            className = "env_name_test"
            break
        case util.EnvironmentAppleReview:
            envName = "APPLE REVIEW"
            className = "env_name_apl"
            break
        case util.EnvironmentStaging:
            envName = "STAGING"
            className = "env_name_stg"
            break
        case util.EnvironmentProduction:
            envName = "PRODUCTION"
            className = "env_name_prd"
            break
        default:
            envName = "LOCAL"
            className = "env_name_loc"
            break
        }

        // リファレンスデータを作成
        ret := map[uint32]map[string]map[string]rest.ResourceDefine{}

        for ver, vermap := range *rest.GetEngine().GetResourceManager().GetAllResources() {
            if _, ok := ret[ver]; !ok {
                ret[ver] = map[string]map[string]rest.ResourceDefine{}
            }

            for s, pathmap := range *vermap {
                if pathmap != nil && *pathmap != nil {
                    dir, _ := path.Split((*pathmap).GetPath())

                    if group, ok := ret[ver][dir]; ok {
                        group[s] = *(*pathmap).Define()
                    } else {
                        ret[ver][dir] = map[string]rest.ResourceDefine{s:*(*pathmap).Define()}
                    }
                }
            }
        }

        // jsonデータを作成
        data := map[string]interface{}{
            "Name": config.Values.Server.Reference.Name,
            "EnvName": config.Values.Server.Reference.Environment,
            "EnvCode": envCode,
            "EnvClass": className,
            "Logo": config.Values.Server.Reference.Logo,
            "UserAgent": config.Values.Server.Reference.UserAgent,
            "Versions": config.Values.Server.Versions,
            "Resources": ret}

        // jsonにエンコード
        buffer, err := json.Marshal(data)

        if err != nil {
            w.WriteHeader(http.StatusForbidden)
        } else {
            w.Header().Set("Content-Type", "application/json")
            _, _ = w.Write(buffer)
        }
    } else if util.FileExists(realPath) {
        buffer, err := ioutil.ReadFile(realPath)
        if err != nil {
            w.WriteHeader(http.StatusForbidden)
        } else {
            mimeType := mime.TypeByExtension(path.Ext(realPath))
            w.Header().Set("Content-Type", mimeType)
            _, _ = w.Write(buffer)
        }
    } else {
        w.WriteHeader(http.StatusNotFound)
    }
}

func Initialize(configPath string) error {
    var err error = nil

    // configをロード
    if err == nil {
        err = config.Load(configPath)
    }

    // 出力するログレベルを設定
    if err == nil {
        log.SetLogLevel(config.Values.Server.LogLevel)
    }

    // システム用MessageManagerを初期化
    if err == nil {
        defaultLang := []message.AcceptLanguage{{Lang: strings.ToLower(config.Values.Message.Default), Q: 1}}
        message.SystemMessageManager = message.CreateMessageManager(&defaultLang)
    }

    // データベースを初期化
    if err == nil {
        err = initDatabase()
    }

    // エンジンを初期化
    if err == nil {
        rest.InitializeEngine()
    }

    return err
}

func AppendResource(version uint32, path string, resource rest.IRestResource) error {
    return rest.GetEngine().AppendResource(version, path, &resource)
}

func SetHooks(hooks rest.ExecutionHooks) {
    rest.GetEngine().SetHooks(&hooks)
}

func Listen() error {
    // API用のハンドラを追加
    for _, v := range config.Values.Server.Versions {
        apiUrl := strings.TrimRight(v.Url, "/") + "/"
        http.HandleFunc(apiUrl, requestHandler)
    }

    // リファレンス用
    if config.Values.Server.Reference.Enable {
        // リファレンス用のハンドラを追加
        referenceUrl := strings.TrimRight(config.Values.Server.Reference.Url, "/") + "/"
        http.HandleFunc(referenceUrl, referenceHandler)
    }

    if config.Values.Server.IsFastCGI {
        // FastCGIモードで起動
        listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", config.Values.Server.Port))
        if err != nil {
            return err
        }
        return fcgi.Serve(listener, nil)
    } else {
        // httpサーバーを起動
        return http.ListenAndServe(fmt.Sprintf(":%d", config.Values.Server.Port), nil)
    }
}

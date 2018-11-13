package goliath

import (
    "github.com/shimalab-jp/goliath/config"
    "github.com/shimalab-jp/goliath/database"
    "github.com/shimalab-jp/goliath/log"
    "github.com/shimalab-jp/goliath/message"
    "github.com/shimalab-jp/goliath/rest"
    "github.com/shimalab-jp/goliath/rest/resources/account"
    "github.com/shimalab-jp/goliath/rest/resources/debug"
    "github.com/shimalab-jp/goliath/rest/resources/fcm"
    "strings"
)

func initDatabase() (error) {
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
            "`token` varchar(256) CHARACTER SET ascii NOT NULL," +
            "`user_id` bigint(20) unsigned NOT NULL," +
            "`is_valid` tinyint(1) unsigned NOT NULL," +
            "`regist_month` int(10) unsigned NOT NULL COMMENT 'yyyymm'," +
            "`created` datetime NOT NULL DEFAULT current_timestamp()," +
            "`modified` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp()," +
            "PRIMARY KEY (`token`)," +
            "KEY `IDX_COUNT1` (`user_id`, `regist_month`)," +
            "KEY `IDX_TRANSIT` (`user_id`)" +
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
            con.Execute("REPLACE INTO `goliath_mst_client_version` (`platform`, `major`, `minor`, `revision`, `resource_version`, `start_time`, `end_time`, `description`) VALUES (1, 1, 0, 0, '', 0, 253402268399, 'Init insert');")
            con.Execute("REPLACE INTO `goliath_mst_client_version` (`platform`, `major`, `minor`, `revision`, `resource_version`, `start_time`, `end_time`, `description`) VALUES (2, 1, 0, 0, '', 0, 253402268399, 'Init insert');")
            con.Execute("REPLACE INTO `goliath_mst_client_version` (`platform`, `major`, `minor`, `revision`, `resource_version`, `start_time`, `end_time`, `description`) VALUES (1, 1, 0, 1, '', 0, 253402268399, 'Init insert');")
            con.Execute("REPLACE INTO `goliath_mst_client_version` (`platform`, `major`, `minor`, `revision`, `resource_version`, `start_time`, `end_time`, `description`) VALUES (2, 1, 0, 1, '', 0, 253402268399, 'Init insert');")
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
        con.Rollback()
    }

    con.Disconnect()
    con = nil

    return err
}

func appendBasicResources() (error) {
    var err error = nil

    if err == nil { err = AppendResource(&account.Regist{}) }
    if err == nil { err = AppendResource(&account.Auth{}) }
    if err == nil { err = AppendResource(&account.Password{}) }
    if err == nil { err = AppendResource(&account.Trans{}) }
    if err == nil { err = AppendResource(&fcm.Regist{}) }
    if err == nil { err = AppendResource(&debug.Cache{}) }

    return err
}

func Initialize(configPath string) (error) {
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

    // 標準リソースを追加
    if err == nil {
        err = appendBasicResources()
    }

    return err
}

package rest

import (
    "fmt"
    "github.com/shimalab-jp/goliath/util"
    "sort"
    "strings"

    "github.com/shimalab-jp/goliath/database"
    "github.com/shimalab-jp/goliath/log"
)

type ClientVersion struct {
    major uint32
    minor uint32
    revision uint32
    environment string
    platform string
    resourceVersion string
}

type ClientVersions []ClientVersion

func (cvs ClientVersions) Len() int {
    return len(cvs)
}

func (cvs ClientVersions) Swap(i, j int) {
    cvs[i], cvs[j] = cvs[j], cvs[i]
}

func (cvs ClientVersions) Less(i, j int) bool {
    return cvs[i].CompareTo(&cvs[j]) == -1
}

func (cv *ClientVersion) GetMajor() (uint32) {
    return cv.major
}

func (cv *ClientVersion) GetMinor() (uint32) {
    return cv.minor
}

func (cv *ClientVersion) GetRevision() (uint32) {
    return cv.revision
}

func (cv *ClientVersion) GetVersion() (string) {
    return fmt.Sprintf("%d.%d.%d", cv.major, cv.minor, cv.revision)
}

func (cv *ClientVersion) GetEnvironmentCode() (uint32) {
    return util.ToEnvironmentCode(cv.environment)
}

func (cv *ClientVersion) GetPlatform() (uint32) {
    p := strings.ToLower(cv.platform)
    if p == "ios" || p == "iphone" || p == "ipad" || p == "tvos" || p == "watchos" || p == "apple" {
        return PlatformApple
    } else if p == "android" || p == "google" || p == "nexus" || p == "pixels" {
        return PlatformGoogle
    } else {
        return PlatformNone
    }
}

func (cv *ClientVersion) GetResourceVersion() (string) {
    return cv.resourceVersion
}

func (cv *ClientVersion) CompareTo(target *ClientVersion) (int8) {
    if cv.major == target.major && cv.minor == target.minor && cv.revision == target.revision {
        return 0
    }
    if cv.major < target.major || (cv.major == target.major && cv.minor < target.minor) || (cv.major == target.major && cv.minor == target.minor && cv.revision < target.revision) {
        return -1
    } else if cv.major > target.major || (cv.major == target.major && cv.minor > target.minor) || (cv.major == target.major && cv.minor == target.minor && cv.revision > target.revision) {
        return 1
    }
    return 0
}

func (cv *ClientVersion) GetRequireVersions() (ClientVersions, error) {
    var ret ClientVersions

    mem := Memcached{}
    {
        err := mem.Get("REQUIRE_CLIENT_VERSIONS", ret)
        if err == nil && ret != nil && len(ret) > 0 {
            return ret, nil
        }
    }

    {
        con, err := database.Connect("goliath")
        if err != nil {
            return nil, err
        }

        result, err := con.Query("SELECT `platform`, `major`, `minor`, `revision`, `resource_version` FROM `goliath_mst_client_version` WHERE `start_time` <= UNIX_TIMESTAMP() && UNIX_TIMESTAMP() <= `end_time`")
        if err != nil {
            return nil, err
        }
        defer result.Rows.Close()

        for result.Rows.Next() {
            var platform uint8
            var major, minor, revision uint32
            var resVersion string

            result.Rows.Scan(&platform, &major, &minor, &revision, &resVersion)

            platformStr := "PC"
            if platform == PlatformApple {
                platformStr = "iOS"
            } else if platform == PlatformGoogle {
                platformStr = "Android"
            }

            ret = append(ret, ClientVersion{
                major: major,
                minor: minor,
                revision: revision,
                resourceVersion: resVersion,
                platform: platformStr})
        }

        sort.Sort(ret)
        log.D("%+v", ret)
    }

    {
        err := mem.Set("REQUIRE_CLIENT_VERSIONS", ret)
        if err != nil {
            log.Ee(err)
        }
    }

    return ret, nil
}

func (cv *ClientVersion) GetUpdateRequireVersion() (*ClientVersion, error) {
    // 稼働中のバージョン番号の一覧を取得
    versions, err := cv.GetRequireVersions()
    if err != nil {
        return nil, err
    }

    var reqVer *ClientVersion = nil
    for _, v := range versions {
        // 違うプラットフォームの場合は次へ
        if cv.GetPlatform() != v.GetPlatform() {
            continue
        }

        // バージョン番号を比較
        r := cv.CompareTo(&v)

        if r == 0 {
            // 同じ場合(可動中のバージョンである)
            return nil, nil
        } else if r == -1 {
            // cvが小さい場合(稼働中に大きいバージョンがある)
            reqVer = &v
        } else {
            // cvが大きい場合(大きいバージョン番号は開発中or審査中)
        }
    }
    return reqVer, nil
}

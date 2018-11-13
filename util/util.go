package util

import (
    "fmt"
    "math/rand"
    "os"
    "strings"
    "time"
)

const (
    EnvironmentLocal       = 9
    EnvironmentDemo        = 8
    EnvironmentDevelop1    = 6
    EnvironmentDevelop2    = 5
    EnvironmentTest        = 4
    EnvironmentStaging     = 3
    EnvironmentAppleReview = 2
    EnvironmentProduction  = 1
)

func GenerateUuid() (string) {
    rand.Seed(time.Now().UnixNano())
    return fmt.Sprintf("%04X%04X-%04X-%04X-%04X-%04X%04X%04X",
        // 32 bits for "time_low"
        rand.Int31n(0xffff ), rand.Int31n(0xffff),

        // 16 bits for "time_mid"
        rand.Int31n(0xffff),

        // 16 bits for "time_hi_and_version",
        // four most significant bits holds version number 4
        rand.Int31n(0x0fff) | 0x4000,

        // 16 bits, 8 bits for "clk_seq_hi_res",
        // 8 bits for "clk_seq_low",
        // two most significant bits holds zero and one for variant DCE1.1
        rand.Int31n(0x3fff ) | 0x8000,

        // 48 bits for "node"
        rand.Int31n(0xffff), rand.Int31n(0xffff), rand.Int31n(0xffff))
}

func FileExists(path string) (bool) {
    file, err := os.Stat(path)
    return err == nil && !file.IsDir()
}

func DirectoryExists(path string) (bool) {
    file, err := os.Stat(path)
    return err == nil && file.IsDir()
}

func ToEnvironmentCode(name string) (uint32) {
    e := strings.ToLower(name)
    if e == "dem" || e == "demo" {
        return EnvironmentDemo
    } else if e == "dev" || e == "develop" || e == "development" {
        return EnvironmentDevelop1
    } else if e == "dev1" || e == "develop1" || e == "development1" {
        return EnvironmentDevelop1
    } else if e == "dev2" || e == "develop2" || e == "development2" {
        return EnvironmentDevelop2
    } else if e == "tst" || e == "test" {
        return EnvironmentTest
    } else if e == "stg" || e == "staging" {
        return EnvironmentStaging
    } else if e == "apl" || e == "app" || e == "apple" {
        return EnvironmentAppleReview
    } else if e == "prd" || e == "prod" || e == "product" || e == "production" || e == "live" {
        return EnvironmentProduction
    } else {
        return EnvironmentLocal
    }
}

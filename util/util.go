package util

import (
    "fmt"
    "math/rand"
    "os"
    "time"
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
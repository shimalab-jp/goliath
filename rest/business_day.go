package rest

import (
    "fmt"
    "github.com/shimalab-jp/goliath/config"
    "strconv"
    "time"
)

type businessDay struct {
    BusinessDay int32
    Year        int32
    Month       int32
    Day         int32
    OpenTime    int64
    CloseTime   int64

    YesterdayBusinessDay int32
    YesterdayYear        int32
    YesterdayMonth       int32
    YesterdayDay         int32
    YesterdayOpenTime    int64
    YesterdayCloseTime   int64

    TomorrowBusinessDay int32
    TomorrowYear        int32
    TomorrowMonth       int32
    TomorrowDay         int32
    TomorrowOpenTime    int64
    TomorrowCloseTime   int64

    LastYear  int32
    LastMont  int32
    NextYear  int32
    NextMonth int32
}

func GetBusinessDay(targetTime int64) businessDay {
    currentTime := time.Unix(targetTime, 0)

    tz, _ := time.LoadLocation(config.Values.Server.TimeZone)
    timeObj := currentTime.In(tz)

    todayBusinessDay, _ := strconv.Atoi(timeObj.Format("20060102"))
    todayYear, _ := strconv.Atoi(timeObj.Format("2006"))
    todayMonth, _ := strconv.Atoi(timeObj.Format("01"))
    todayDay, _ := strconv.Atoi(timeObj.Format("02"))
    todayOpenTime, _ := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%04d-%02d-%02d 00:00:00", todayYear, todayMonth, todayDay), tz)
    todayCloseTime, _ := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%04d-%02d-%02d 23:59:59", todayYear, todayMonth, todayDay), tz)

    yesterdayObj := timeObj.AddDate(0, 0, -1)
    yesterdayBusinessDay, _ := strconv.Atoi(yesterdayObj.Format("20060102"))
    yesterdayYear, _ := strconv.Atoi(yesterdayObj.Format("2006"))
    yesterdayMonth, _ := strconv.Atoi(yesterdayObj.Format("01"))
    yesterdayDay, _ := strconv.Atoi(yesterdayObj.Format("02"))
    yesterdayOpenTime, _ := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%04d-%02d-%02d 00:00:00", yesterdayYear, yesterdayMonth, yesterdayDay), tz)
    yesterdayCloseTime, _ := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%04d-%02d-%02d 23:59:59", yesterdayYear, yesterdayMonth, yesterdayDay), tz)

    tomorrowObj := timeObj.AddDate(0, 0, 1)
    tomorrowBusinessDay, _ := strconv.Atoi(tomorrowObj.Format("20060102"))
    tomorrowYear, _ := strconv.Atoi(tomorrowObj.Format("2006"))
    tomorrowMonth, _ := strconv.Atoi(tomorrowObj.Format("01"))
    tomorrowDay, _ := strconv.Atoi(tomorrowObj.Format("02"))
    tomorrowOpenTime, _ := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%04d-%02d-%02d 00:00:00", tomorrowYear, tomorrowMonth, tomorrowDay), tz)
    tomorrowCloseTime, _ := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%04d-%02d-%02d 23:59:59", tomorrowYear, tomorrowMonth, tomorrowDay), tz)

    lastYear := timeObj.AddDate(-1, 0, 0).Year()
    lastMonth := timeObj.AddDate(0, -1, 0).Year()
    nextYear := timeObj.AddDate(1, 0, 0).Year()
    nextMonth := timeObj.AddDate(0, 1, 0).Year()

    return businessDay{
        BusinessDay: int32(todayBusinessDay),
        Year:        int32(todayYear),
        Month:       int32(todayMonth),
        Day:         int32(todayDay),
        OpenTime:    todayOpenTime.Unix(),
        CloseTime:   todayCloseTime.Unix(),

        YesterdayBusinessDay: int32(yesterdayBusinessDay),
        YesterdayYear:        int32(yesterdayYear),
        YesterdayMonth:       int32(yesterdayMonth),
        YesterdayDay:         int32(yesterdayDay),
        YesterdayOpenTime:    yesterdayOpenTime.Unix(),
        YesterdayCloseTime:   yesterdayCloseTime.Unix(),

        TomorrowBusinessDay: int32(tomorrowBusinessDay),
        TomorrowYear:        int32(tomorrowYear),
        TomorrowMonth:       int32(tomorrowMonth),
        TomorrowDay:         int32(tomorrowDay),
        TomorrowOpenTime:    tomorrowOpenTime.Unix(),
        TomorrowCloseTime:   tomorrowCloseTime.Unix(),

        LastYear:  int32(lastYear),
        LastMont:  int32(lastMonth),
        NextYear:  int32(nextYear),
        NextMonth: int32(nextMonth)}
}

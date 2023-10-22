package pkg

import (
	"log"
	"time"
)

func TimeStrToUnixMilli(timeStr string) uint64 {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		log.Println("时间解析失败:", err)
		return 0
	}
	return uint64(t.UnixMilli())
}

func DateStrToUnixMilli(dateStr string) uint64 {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Println("时间解析失败:", err)
		return 0
	}
	return uint64(t.UnixMilli())
}

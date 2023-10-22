package pkg

import (
	"log"
	"strconv"
)

func StrToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		log.Println("数值转换错误:", err)
		return 0
	}
	return f
}

package helpers

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func String2Float(str string) float32 {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		num = 0.0
	}
	return float32(num)
}

func String2Int(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		num = 0
	}
	return num
}

// String2Slice : String to entity "LB CLUB, Wanner Music VN" -> []Brand{}
func String2Slice(str string) []string {
	slice := strings.Split(str, ",")
	for i, name := range slice {
		slice[i] = strings.TrimSpace(name)
	}
	return slice
}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var result string
	for i := 0; i < length; i++ {
		result += string(charset[rand.Intn(len(charset))])
	}
	return result
}

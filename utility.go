package gows

import (
	"math/rand"
	"strconv"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"


func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func Uuid() string {
	time:=time.Now().Unix()
	uuid:=strconv.FormatInt(time,10)
	for i:=0;i<=4;i++{
		uuid+="-"+randString(5)
	}

	return uuid
}

package main

import (
	"time"
)

var cache = make(map[string]cacheValue)
var validOps = [4]string{"GET", "SET", "CLEAR", "DELETE"}

func isValidOp(method string) bool {
	for _, value := range validOps {
		if value == method {
			return true
		}
	}
	return false
}

func setCache(key string, data string, expires int64) string {
	if expires == 0 {
		expires = time.Now().Unix()*1000 + 180000
	}
	value := cacheValue{
		data:    data,
		expires: expires,
	}
	cache[key] = value
	return value.data
}

func getCache(key string) string {
	if val, exists := cache[key]; exists == true {
		return val.data
	}
	return ""
}

func deleteCache(key string) {
	delete(cache, key)
}

func clearCache() {
	cache = make(map[string]cacheValue)
}

func cleanCache() {
	var currentdate = time.Now().Unix() * 1000

	for k, v := range cache {
		if v.expires > 0 && v.expires < currentdate {
			delete(cache, k)
		}
	}
}

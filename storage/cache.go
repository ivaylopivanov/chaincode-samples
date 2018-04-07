package main

import (
	"time"

	"github.com/ivaylopivanov/chaincode-samples/storage/cache"
)

var (
	cleanupInterval   = 30 * time.Second
	defaultExpiration = 10 * time.Minute
	mcache            = cache.New(defaultExpiration, cleanupInterval)
)

func cacheCheck(id, key string) bool {
	key = id + "-" + key
	_, ok := mcache.Get(key)
	if !ok {
		setInCache(key)
	}
	return ok
}

func setInCache(key string) {
	mcache.Add(key, nil, defaultExpiration)
}

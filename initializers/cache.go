package initializers

import (
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	/*
	不过大家要注意，此方法把数据存在一个hashmap里了，不是字符串，所以用的时候一定要想清楚了。

	若是想用字符串的方式，可以用Redigo。*/
	"fmt"
	"strconv"
)

var (
	GlobalCache cache.Cache
	cacheErr error
)

func init()  {
	GlobalCache, cacheErr = cache.NewCache("file", `{"CachePath":"./tmp/cache","FileSuffix":".cache","DirectoryLevel":"2","EmbedExpiry":"10"}`)
	fmt.Println(cacheErr)
}

func GoidMapKey(key string) string {
	id := Goid()

	return key + ":" + strconv.Itoa(id)
}

func RemoveCaches() {
	sessionKey := GoidMapKey("sessionId")
	GlobalCache.Delete(sessionKey)

	requestKey := GoidMapKey("requestId")
	GlobalCache.Delete(requestKey)
}
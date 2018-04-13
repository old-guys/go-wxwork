package lib_wxwork

import (
	"github.com/astaxie/beego"
	"os"
	"log"
	"fmt"
	"wxwork/initializers"
	"time"
	//"strings"
)

var (
	Logger loggerConfig
)

type loggerConfig struct {
	initialized bool
	logger *log.Logger
}

func (c *loggerConfig) getLogger() *log.Logger {
	fileName := beego.AppPath + "/logs/wxwork.log"
	logFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalln("open file error")
	}

	debugLog := log.New(logFile,"", log.Llongfile)

	return debugLog
}

func (c *loggerConfig) Info(args ...interface{}) {
	if !c.initialized {
		c.initialized = true
		c.logger = c.getLogger()
	}

	// session id
	sessionKey := initializers.GoidMapKey("sessionId")
	sessionId := initializers.GlobalCache.Get(sessionKey)

	// request id
	requestKey := initializers.GoidMapKey("requestId")
	requestId := initializers.GlobalCache.Get(requestKey)

	v := make([]interface{}, len(args) + 1)
	v[0] = "\n[" + time.Now().Format("2006-01-02T15:04:05.999999") + "] [sessionID:" + sessionId.(string) + "] [requestId:"+ requestId.(string) + "]\n"

	for i, m := range args {
		v[i + 1] = m
	}

	c.logger.Output(2, fmt.Sprintln(v...))
}

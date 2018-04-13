package main

import (
	_ "wxwork/routers"
	_ "wxwork/initializers"
	_ "wxwork/models"
	_ "wxwork/workers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/satori/go.uuid"
	"fmt"
	"wxwork/initializers"
	"os"
	"strconv"
	"net/http"
	//"github.com/astaxie/beego/grace"
	//"log"
)

var filterRequestId = func(ctx *context.Context){
	requestId := fmt.Sprintf("%v", uuid.Must(uuid.NewV4()))
	ctx.Input.SetData("requestId", requestId)
}

var addCache = func(ctx *context.Context) {
	sessionKey := initializers.GoidMapKey("sessionId")
	initializers.GlobalCache.Put(sessionKey, ctx.Input.CruSession.SessionID(), 0)

	requestKey := initializers.GoidMapKey("requestId")
	initializers.GlobalCache.Put(requestKey, ctx.Input.GetData("requestId"), 0)
}

var removeCaches = func(ctx *context.Context) {
	initializers.RemoveCaches()
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("WORLD!"))
	w.Write([]byte("ospid:" + strconv.Itoa(os.Getpid())))
}

func writePidToFile() {

	fileName := beego.AppPath + "/tmp/pids/server.pid"
	pidFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)

	if err == nil {
		pidFile.Write([]byte(strconv.Itoa(os.Getpid())))
	}
}

func init() {
	// set log
	beego.SetLogger("file", `{"filename": "logs/` + beego.BConfig.RunMode + `.log"}`)
	beego.SetLogFuncCall(true)

	writePidToFile()
}


func main() {
	// session
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "wxwork"

	beego.BConfig.WebConfig.Session.SessionProvider = "file"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = "tmp/sessions"

	beego.BConfig.CopyRequestBody = true

	beego.InsertFilter("/*", beego.BeforeRouter, filterRequestId)
	beego.InsertFilter("/*", beego.BeforeRouter, addCache)
	beego.InsertFilter("/*", beego.FinishRouter, removeCaches)

	//mux := http.NewServeMux()
	//mux.HandleFunc("/hello", handler)
	//err := grace.ListenAndServe("localhost:7510", mux)
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Println("Server on 8080 stopped")
	//os.Exit(0)

	beego.Run()
}


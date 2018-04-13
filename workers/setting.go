package workers

import (
	"github.com/benmanns/goworker"
	"os"
	"os/exec"
	"strconv"
	"github.com/astaxie/beego"
	"strings"
)

func writeWorkerPidToFile() {

	execPath, _ := exec.LookPath(os.Args[0])
	if !strings.Contains(execPath, "wxwork_worker") { return }

	fileName := beego.AppPath + "/tmp/pids/worker.pid"
	pidFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)

	if err == nil {
		pidFile.Write([]byte(strconv.Itoa(os.Getpid())))
	}
}

func init() {
	settings := goworker.WorkerSettings{
		URI:            "redis://localhost:6379/1",
		Connections:    100,
		Queues:         []string{"myqueue", "delimited", "queues"},
		UseNumber:      true,
		ExitOnComplete: false,
		Concurrency:    2,
		Namespace:      "resque:",
		Interval:       5.0,
	}
	goworker.SetSettings(settings)
	goworker.Register("WxworkOrgChangeContact", changeContact)

	writeWorkerPidToFile()
}
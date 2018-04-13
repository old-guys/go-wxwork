### supervisor
	https://blog.csdn.net/hylexus/article/details/78177649?locationNum=7&fps=1
	supervisord -c /usr/local/etc/supervisord.ini
	supervisorctl -c /usr/local/etc/supervisord.ini

### package
	go get github.com/astaxie/beego
	go get github.com/beego/bee
	go get github.com/garyburd/redigo/redis
	go get github.com/imdario/mergo
	go get github.com/jinzhu/gorm
	go get github.com/go-sql-driver/mysql
	go get github.com/satori/go.uuid
	go get gopkg.in/yaml.v2
	go get github.com/json-iterator/go
	go get github.com/benmanns/goworker
	
### compile && run
 	go build -o wxwork main.go
	go build -o wxwork_worker worker.go
	./wxwork
	./wxwork_worker -queues="myqueue"

### deploy
	source scripts/deploy/test.sh

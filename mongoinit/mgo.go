package mongoinit

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"

	"gogs.itcloud.pro/SAS-project/sas/app"
)

var Config *app.ConfigType

//var DBname = cfg.DBname
//var mongoUrl = cfg.Mongourl

//var DB *mgo.Session
//var MongoDBHosts = cfg.Mongourl

// const (
// 	MongoDBHosts = cfg.Mongourl
// 	TestDatabase = "DevBase"
// )
//
//func InitMongo() {
//	cfg = app.GetConfig()
//	go func() {
//		if err := StartMgo(); err != nil {
//			log.Fatal(err)
//		}
//	}()
//
//	go func() {
//		time.Sleep(1 * time.Second)
//		var err error
//		DB, err = Connect()
//		if err != nil {
//			log.Fatal(err)
//		}
//	}()
//	time.Sleep(2 * time.Second)
//
//}

//mongod --config F:\MongoDB\bin\mongodb.config
func StartMgo() error {
	cmd := exec.Command("mongod", "--config", "F:\\MongoDB\\bin\\mongodb.config")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	input := bufio.NewScanner(stdout)
	go func() {
		for input.Scan() {
			fmt.Fprintf(os.Stdout, "MONGO: %s\n", input.Text())
		}
	}()
	return cmd.Wait()
}

//
//func Connect() (*mgo.Session, error) {
//	info := &mgo.DialInfo{
//		Addrs: []string{MongoDBHosts},
//	}
//
//	return mgo.DialWithInfo(info)
//}

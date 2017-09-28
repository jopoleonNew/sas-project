package yad

import (
	"net/http"

	"github.com/jeffail/tunny"
	"github.com/sirupsen/logrus"
)

var YPool *tunny.WorkPool

func InitPool(numWorkers int) {
	var err error
	YPool, err = tunny.CreatePool(numWorkers, func(object interface{}) interface{} {
		r, _ := object.(*http.Request)
		//time.Sleep(1000 * time.Millisecond)
		client := &http.Client{}
		//r.URL.RawQuery = "ID=" + "123"
		resp, err := client.Do(r)
		if err != nil {
			return err
		}
		// Do something that takes a lot of work
		output := resp
		//logrus.Info("\nWorker used with request: ", r, "\n")
		return output
	}).Open()
	if err != nil {
		logrus.Fatal("YandexDirectAPI workerspool InitPool error:", err)
	}

}

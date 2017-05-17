package main

import (
	"log"

	"gogs.itcloud.pro/SAS-project/sas/model"
)

func main() {
	for i := 0; i < 20; i++ {
		log.Println("GenerateKey32chars: ", model.GenerateKey32chars())
		log.Println("RandStringBytes:    ", model.RandStringBytes(32))
	}

}

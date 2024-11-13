package main

import (
	"github.com/whyxn/easynas/backend/pkg/nas"
	"log"
)

func main() {
	zpools, err := nas.ListZpools()
	if err != nil {
		log.Println("Error:", err)
		return
	}

	log.Println(zpools)

	zfsVolumes, err := nas.ListZFSVolumes()
	if err != nil {
		log.Println("Error:", err.Error())
		return
	}

	log.Println(zfsVolumes)
}

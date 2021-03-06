package main

import (
	"bufio"
	"fmt"
	"github.com/punxlab/sadwave-tools/custom-photo-uploader/client"
	"github.com/punxlab/sadwave-tools/custom-photo-uploader/config"
	"github.com/punxlab/sadwave-tools/custom-photo-uploader/photos"
	"log"
	"os"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		printAndWait(err.Error())
		return
	}

	cl := client.NewSadwaveClient(cfg.API)

	pts, err := photos.FromCSV(cfg.SourceCSVPath)
	if err != nil {
		printAndWait(err.Error())
		return
	}

	token, err := cl.RequestToken(cfg.AdminUser, cfg.AdminPassword)
	if err != nil {
		printAndWait(err.Error())
	}

	for _, p := range pts {
		log.Println(fmt.Sprintf("uploading photo %s for event %s", p.PhotoUrl, p.EventUrl))
		err := cl.SetCustomPhoto(p.EventUrl, p.PhotoUrl, token)
		if err != nil {
			log.Println(err)
		}
	}

	printAndWait("completed")
}

func printAndWait(msg string) {
	log.Println(msg)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

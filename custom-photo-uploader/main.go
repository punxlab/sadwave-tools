package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/punxlab/sadwave-tools/custom-photo-uploader/client"
	"github.com/punxlab/sadwave-tools/custom-photo-uploader/config"
	"github.com/punxlab/sadwave-tools/custom-photo-uploader/photos"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		printAndWait(err.Error())
		return
	}

	pts, err := photos.FromCSV(cfg.SourceCSVPath)
	if err != nil {
		printAndWait(err.Error())
		return
	}

	for _, p := range pts {
		log.Println(fmt.Sprintf("uploading photo %s for event %s", p.EventUrl, p.PhotoUrl))
		err := client.NewSadwaveClient(cfg.API).SetCustomPhoto(p.EventUrl, p.PhotoUrl)
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

package main

import (
	"log"
	"log/syslog"

	"net/http"
	"strings"
	"os"
)

const path = "/tmp"

func handler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[1:]
	pp := strings.Split(url, "/")

	if pp[0] == "file" {
		data := r.FormValue("data")
		writeFile(pp[1], data)
	} else {
		log.Print("invalid request path", url)
	}
}

func main() {
	l, e := syslog.New(syslog.LOG_NOTICE, "recfile")
	if e == nil {
		log.SetOutput(l)
	}

    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":9091", nil))
}

func writeFile(fileName, data string) {
	filePath := path+"/"+fileName
	var file *os.File

	s, err := os.Stat(filePath)
	log.Print("temp file stat", s)

	if err != nil {
		log.Print(err)
		if os.IsNotExist(err) {
			log.Print("Creating temp file")
			file, err = os.Create(filePath)
			
			if err != nil {
				log.Fatalf("failed creating file: %s", err)
			}	
		}
	} else {
		file, err = os.OpenFile(path+"/"+fileName, os.O_RDWR, 0644)
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}	
	}
     
    defer file.Close()
     
    _, err = file.WriteAt([]byte(data), 0) // Write at 0 beginning
    if err != nil {
        log.Fatalf("failed writing to file: %s", err)
	}
}
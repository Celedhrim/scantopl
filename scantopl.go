package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jnovack/flag"
	log "github.com/sirupsen/logrus"
)

func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(path.Base(fn), path.Ext(fn))
}

func TitleFromFileName(fn string) string {
	return strings.TrimPrefix(FilenameWithoutExtension(fn), "pl_")
}

func createForm(form map[string]string) (string, io.Reader, error) {
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	defer mp.Close()
	for key, val := range form {
		if strings.HasPrefix(val, "@") {
			val = val[1:]
			file, err := os.Open(val)
			if err != nil {
				return "", nil, err
			}
			defer file.Close()
			part, err := mp.CreateFormFile(key, val)
			if err != nil {
				return "", nil, err
			}
			io.Copy(part, file)
		} else {
			mp.WriteField(key, val)
		}
	}
	return mp.FormDataContentType(), body, nil
}

func uploadFile(document, plurl, pltoken string) {
	form := map[string]string{"document": "@" + document, "title": TitleFromFileName(document)}
	ct, body, err := createForm(form)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", plurl+"/api/documents/post_document/", body)
	req.Header.Set("Content-Type", ct)
	req.Header.Set("Authorization", "Token "+pltoken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Warn("Upload failed with response code: ", resp.StatusCode)
	} else {
		log.Info("Upload succesfully, remove file")
		err := os.Remove(document)
		if err != nil {
			log.Warn(err)
		}
	}
}

func main() {

	// OPTs
	flag.String(flag.DefaultConfigFlagname, "", "path to config file")
	scandir := flag.String("scandir", "/home/scanservjs/output", "Scanserjs ouput directory")
	plurl := flag.String("plurl", "http://localhost:8080", "The paperless instance URL without trailing /")
	pltoken := flag.String("pltoken", "xxxxxxxxxxxxxxxxxx", "Paperless auth token , generated through admin")

	flag.Parse()

	// test if it's a directory
	log.Println("Start watching:", *scandir)

	// The watch loop
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//log.Println("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create && strings.HasPrefix(path.Base(event.Name), "pl_") {
					//little pause to ensure the write operation is finished
					time.Sleep(1 * time.Second)
					log.Info("New file to upload:", event.Name)
					uploadFile(event.Name, *plurl, *pltoken)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Warn("error:", err)
			}
		}
	}()

	//err = watcher.Add("/home/docker_data/scanservjs/output")
	err = watcher.Add(*scandir)

	if err != nil {
		log.Fatal(err)
	}
	<-done
}

package servers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const countBytesBuffer = 10000000

var countUploads = 0
var countDownloads = 0

func httpDownload(w http.ResponseWriter, req *http.Request) {
	countDownloads++
	idDownload := countDownloads
	path := strings.Split(req.URL.Path, "/")
	sizeQuery := path[len(path)-1]
	size, err := strconv.Atoi(sizeQuery)
	if err != nil {
		return
	}
	log.Printf("DOWNLOAD STARTED: #%d, %dMB\n", idDownload, size/1e6)
	countBytesToSend := size
	w.Header().Set("Content-Length", sizeQuery)
	w.Header().Set("Content-Type", "application/octet-stream")
	var bytes []byte = make([]byte, countBytesBuffer)
	tStart := time.Now()
	for countBytesToSend > 0 {
		if countBytesToSend >= countBytesBuffer {
			_, err = w.Write(bytes)
		} else {
			_, err = w.Write(bytes[:countBytesToSend])
			break
		}
		if err != nil {
			break
		}
		countBytesToSend -= countBytesBuffer
	}
	tStop := time.Now()
	dt := tStop.Sub(tStart)
	rateBytesPerSecond := float64(size) / dt.Seconds()
	rateBitsPerSecond := rateBytesPerSecond * 8
	log.Printf("DOWNLOAD FINISHED: #%d, %dMB, %s, %.0fMb/s (%.0fMB/s)\n", idDownload, (size-countBytesToSend)/1e6, dt, rateBitsPerSecond/1e6, rateBytesPerSecond/1e6)
}

func httpUpload(w http.ResponseWriter, req *http.Request) {
	countUploads++
	idUpload := countUploads
	contentLength := int(req.ContentLength)
	countBytesToRead := contentLength
	log.Printf("UPLOAD STARTED: #%d, %dMB\n", idUpload, countBytesToRead/1e6)
	var bytes []byte = make([]byte, countBytesBuffer)
	tStart := time.Now()
	for countBytesToRead > 0 {
		n, err := req.Body.Read(bytes)
		if n > 0 {
			countBytesToRead -= n
		}
		if err != nil {
			break
		}
	}
	tStop := time.Now()
	dt := tStop.Sub(tStart)
	countBytesRead := contentLength - countBytesToRead
	rateBytesPerSecond := float64(countBytesRead) / dt.Seconds()
	rateBitsPerSecond := rateBytesPerSecond * 8
	log.Printf("UPLOAD FINISHED: #%d, %dMB, %s, %.0fMb/s (%.0fMB/s)\n", idUpload, countBytesRead/1e6, dt, rateBitsPerSecond/1e6, rateBytesPerSecond/1e6)
}

func StartServerHTTP(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		log.Printf("The HTTP server has started on port %d\n", portTCP_HTTP)
		defer wg.Done()
		defer log.Printf("The HTTP server has stopped\n")
		http.HandleFunc("/download/", httpDownload)
		http.HandleFunc("/upload", httpUpload)
		portStr := fmt.Sprintf(":%d", portTCP_HTTP)
		err := http.ListenAndServe(portStr, nil)
		if err != nil {
			panic(err)
		}
	}()
}

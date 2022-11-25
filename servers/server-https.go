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

var countUploadsHttps = 0
var countDownloadsHttps = 0

func httpsDownload(w http.ResponseWriter, req *http.Request) {
	countDownloadsHttps++
	idDownload := countDownloadsHttps
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

func httpsUpload(w http.ResponseWriter, req *http.Request) {
	countUploadsHttps++
	idUpload := countUploadsHttps
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

func StartServerHTTPS(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		log.Printf("The HTTPS server has started on port %d\n", portTCP_HTTPS)
		defer wg.Done()
		defer log.Printf("The HTTPS server has stopped\n")
		mux := http.NewServeMux()
		mux.HandleFunc("/download/", httpDownload)
		mux.HandleFunc("/upload", httpUpload)
		portStr := fmt.Sprintf(":%d", portTCP_HTTPS)
		server := &http.Server{Addr: portStr, Handler: mux}
		server.SetKeepAlivesEnabled(false)
		err := server.ListenAndServeTLS("server.crt", "server.key")
		if err != nil {
			panic(err)
		}
	}()
}

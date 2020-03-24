package main

import (
	"fmt"
	"io"
	"log"
	"time"
	"strconv"
	"net/http"
	"crypto/md5"
	"os"
    "text/template"
    "flag"
)

var (
    port string
    uploadDir = "./upload/"
    uploadTmp = `
        <html>
        <head>
            <title>Upload file</title>
        </head>
        <body>
        <form enctype="multipart/form-data" action="/" method="post">
            <input type="file" name="uploadfile" />
            <input type="hidden" name="token" value="{{.}}"/>
            <input type="submit" value="upload" />
        </form>
        </body>
        </html>
    `
)

func init() {
	const (
		defaultPort = "9090"
		usage = "the port of listen"
	)
	flag.StringVar(&port, "port", defaultPort, usage)
    flag.StringVar(&port, "p", defaultPort, usage+" (shorthand)")
    
    if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
        err = os.Mkdir(uploadDir, os.ModeDir | os.ModePerm)
        if err != nil {
            log.Fatal("mkdir err: ", err)
        }
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
       if r.Method == "GET" {
           crutime := time.Now().Unix()
           h := md5.New()
           io.WriteString(h, strconv.FormatInt(crutime, 10))
           token := fmt.Sprintf("%x", h.Sum(nil))

           t, _ := template.New("upload.html").Parse(uploadTmp)
           t.Execute(w, token)
       } else if r.Method == "POST" {
           
           r.ParseMultipartForm(32 << 20)
           file, handler, err := r.FormFile("uploadfile")
           if err != nil {
               fmt.Println(err)
               return
           }
           defer file.Close()
           fmt.Printf("uploading %v to %v", handler.Filename, uploadDir+handler.Filename)
           fmt.Fprintf(w, "%v", handler.Header)
           f, err := os.OpenFile(uploadDir+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
           if err != nil {
               fmt.Println(err)
               return
           }
           defer f.Close()
           io.Copy(f, file)
       } else {

       }
}

func main(){
    flag.Parse()
    http.HandleFunc("/", upload)
    addr := fmt.Sprintf(":%v", port)
    fmt.Println("listen: " + addr)
	err := http.ListenAndServe(addr, nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}


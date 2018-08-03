package web

import (
	"fmt"
    "log"
    "html/template"
    "io"
    "net/http"
    "os"
    "runtime"
    "path"
    "errors"

    "github.com/adyachok/oh/models"
)

var ch chan<- *models.Command

 // taken from: https://undebugable.wordpress.com/2017/04/15/golang-simple-file-upload-using-go-languange/
func upload(w http.ResponseWriter, r *http.Request) {
	path, err := dirname()
	if err != nil {
		panic("Cannot find templates directory. Cannot start server.")
	}
    if r.Method == "GET" {
        t, _ := template.ParseFiles(path + "/upload.gtpl") 
        t.Execute(w, nil)
 
    } else if r.Method == "POST" {
        file, handler, err := r.FormFile("uploadfile")
        if err != nil {
            log.Panic(err)
            return
        }
        defer file.Close()
 
        fmt.Fprintf(w, "%v", handler.Header)
        f, err := os.OpenFile("/tmp/" + handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
        if err != nil {
           log.Panic(err)
            return
        }
        defer f.Close()
 
        io.Copy(f, file)
        ch <- models.NewCommand("upload", "/tmp/" + handler.Filename)
 
    } else {
        log.Println("Unknown HTTP " + r.Method + "  Method")
    }
}

func dirname() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	var err error
    if !ok {
        log.Panic("No caller information")
        err = errors.New("Cannot find dirname")
    }
    return path.Dir(filename), err
}
 
func RunServer(cCh chan<- *models.Command) {
	ch = cCh
    http.HandleFunc("/upload", upload)
    http.ListenAndServe(":8080", nil)
}
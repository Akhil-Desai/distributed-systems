package main

import (
	"fmt"
	"net/http"
    "io"
    "os"
)

var fileMetadataStore map[string]map[string]interface{}

func uploadFile(w http.ResponseWriter, r *http.Request){
    fmt.Println("File upload endpoint hit!")

    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        fmt.Println("Error parsing file", err)
        return
    }

    file,handler,err := r.FormFile("clientFile")
    if err != nil{
        fmt.Println("Error", err)
        return
    }
    defer file.Close()

    dst,err := os.Create("../store/" + handler.Filename)

    if err != nil{
        fmt.Println("Create file err:", err)
    }

    defer dst.Close()

    _,err = io.Copy(dst, file)
    if err != nil {
        fmt.Println("Copy err:", err)
        return
    }
    fmt.Fprintf(w, "File saved: %s", handler.Filename)

    fileMetadataStore[handler.Filename] = map[string]interface{} {
        "size": handler.Size,
        "header": handler.Header,
    }

    return
}

func downloadFile(w http.ResponseWriter, r *http.Request){

    fileName := r.URL.Query().Get("Filename")

    //Check if FileName exist
    exist := fileMetadataStore[fileName]

    if exist == nil {
        fmt.Println("Error file does not exist in the store")
        return
    }

    file,err := os.Open("../store/" + fileName)
    if err != nil {
        fmt.Println("Error", err)
        return
    }
    defer file.Close()

    w.Header().Set("Content-Disposition", "attachment; filename=" + fileName)
    w.Header().Set("Content-Type", "application/octet-stream")

    _,err = io.Copy(w, file)
    if err != nil{
        fmt.Println("Error writing file to client", err)
        return
    }
	return
}

func main(){

    fileMetadataStore = make(map[string]map[string]interface{})

	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/download", downloadFile)

	if err := http.ListenAndServe(":8080",nil ); err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println("Listening at port 8080")

}

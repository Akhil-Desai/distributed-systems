package main

import (
	"fmt"
	"net/http"
)

func uploadFile(w http.ResponseWriter, r *http.Request){
	return
}

func downloadFile(w http.ResponseWriter, r *http.Request){
	return
}

func main(){
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/download", downloadFile)

	if err := http.ListenAndServe(":8080",nil ); err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println("Listening at port 8080")
}

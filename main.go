package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/tamarakaufler/go_dynamodb_search/search"
)

var (
	tmplDir            string
	awsAccessKeyID     string
	awsSecretAccessKey string
)

func init() {
	tmplDir = "templates/"

}

func formHandler(w http.ResponseWriter, r *http.Request) {

	form, err := ioutil.ReadFile(tmplDir + "search_form.html")

	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("404 Problem - " + http.StatusText(404)))
	}

	w.Write(form)
}

func displayHandler(w http.ResponseWriter, r *http.Request) {
	data, err := search.DoDynamoDBScan(r)

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		data.Error = err.Error()
	}

	t, _ := template.ParseFiles(tmplDir + "search_result.html")
	t.Execute(w, data)
}

func main() {
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/process", displayHandler)

	fmt.Println("Starting server")
	http.ListenAndServe("127.0.0.1:9000", nil)
}

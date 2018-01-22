package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var html = `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <title>Go4k8s</title>
</head>
<body style="background-color:#E6E6FA">
  <h3>Hello From Go4k8s Webserver</h3>
  <p>Version: %s</p>
</body>
</html>
`
var version = "1.0.0"

func httpHandler(w http.ResponseWriter, r *http.Request) {
	format := "%s - [%s] \"%s %s %s\" %s\n"
	fmt.Printf(format, version, time.Now().Format(time.RFC1123),
		r.Method, r.URL.Path, r.Proto, r.UserAgent())
	fmt.Fprintf(w, html, version)
}

func main() {

	http.HandleFunc("/", httpHandler)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

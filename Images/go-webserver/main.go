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
  <title>Hello Universe</title>
</head>
<body>
  <h3>Hello Universe</h3>
  <p>Hostname: %s</p>
</body>
</html>
`
var hostname = "localhost"

func httpHandler(w http.ResponseWriter, r *http.Request) {
	format := "%s - [%s] \"%s %s %s\" %s\n"
	fmt.Printf(format, hostname, time.Now().Format(time.RFC1123),
		r.Method, r.URL.Path, r.Proto, r.UserAgent())
	fmt.Fprintf(w, html, hostname)
}

func main() {

	http.HandleFunc("/", httpHandler)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

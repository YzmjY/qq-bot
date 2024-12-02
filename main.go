package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		buf := bytes.Buffer{}
		io.Copy(&buf, r.Body)
		fmt.Println(r.Method)
		fmt.Println(r.URL.Path)
		fmt.Println(r.Header)
		fmt.Println(buf.String())
		w.Write([]byte("Hello World!"))
	})

	http.ListenAndServe(":8080", nil)
}

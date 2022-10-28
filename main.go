package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, _ *http.Request) {
		sites := []string{url("/hello"), url("/echo")}
		msg := fmt.Sprintf("%s\n", html(list(sites)))
		mustWrite(res, msg)
	})

	http.HandleFunc("/echo", func(res http.ResponseWriter, req *http.Request) {
		msg := fmt.Sprintf("%#v", req)
		mustWrite(res, msg)
	})

	http.HandleFunc("/hello", func(res http.ResponseWriter, _ *http.Request) {
		mustWrite(res, "Hello, World")
	})

	http.ListenAndServe(":8080", nil)
}

func url(relativePath string) string {
	return fmt.Sprintf("<a href=\"%s\">%s</a>", relativePath, relativePath)
}

func html(content string) string {
	return fmt.Sprintf("<!DOCTYPE html><html><body>%s</body></html>", content)
}

func list(rows []string) string {
	var contents string
	for _, row := range rows {
		contents += fmt.Sprintf("<li>%s</li>", row)
	}
	list := fmt.Sprintf("<ul>%s</ul>", contents)
	return list
}

func mustWrite(response http.ResponseWriter, msg string) {
	_, err := io.WriteString(response, msg)
	mustOk(err)
}

func mustOk(err error) {
	if err != nil {
		panic(err)
	}
}

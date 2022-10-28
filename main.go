package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/hello", hello)

	http.ListenAndServe(":8080", nil)
}

func index(response http.ResponseWriter, request *http.Request) {
	sites := []string{url("/hello"), url("/echo")}
	mustWrite(response, fmt.Sprintf("%s\n", html(list(sites))))
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

func echo(response http.ResponseWriter, request *http.Request) {
	msg := fmt.Sprintf("%#v", request)
	mustWrite(response, msg)

	fmt.Printf("%#v\n===\n", response)
}

func hello(response http.ResponseWriter, request *http.Request) {
	mustWrite(response, "Hello, World")
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

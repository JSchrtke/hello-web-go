package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/register", register)
	http.HandleFunc("/register_click", register_click)

	http.ListenAndServe(":8080", nil)
}

func index(res http.ResponseWriter, _ *http.Request) {
	sites := []string{makeUrl("/hello"), makeUrl("/echo"), makeUrl("/register")}
	msg := fmt.Sprintf("%s\n", makeHtml(makeList(sites)))
	mustRespond(res, msg)
}

func echo(res http.ResponseWriter, req *http.Request) {
	msg := fmt.Sprintf("%#v", req)
	mustRespond(res, msg)
}

func hello(res http.ResponseWriter, _ *http.Request) {
	mustRespond(res, "Hello, World")
}

func register(res http.ResponseWriter, req *http.Request) {
	mustRespond(res, `
<!DOCTYPE html>
<html>
	<body>
		<h1>Register user</h1>
		<form action="/register_click" method="POST">
			<label for="fname">
				First name:
			</label>
			<input type="text" id="fname" name="fname">
			<br><br>
			<label for="lname">
				Last name:
			</label>
			<input type="text" id="lname" name="lname">
			<br><br>
			<label for="email">
				Email:
			</label>
			<input type="text" id="email" name="email">
			<br><br>
			<input type="submit" value="Register">
		</form>
	</body>
</html>`,
	)
}

func register_click(res http.ResponseWriter, req *http.Request) {
	buf, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(
			res,
			fmt.Sprintf("error reading request: %s\n", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	if len(buf) == 0 {
		http.Error(
			res, "error: no request data", http.StatusInternalServerError,
		)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(res).Encode(makeUser(buf))
	if err != nil {
		http.Error(
			res,
			fmt.Sprintf("error encoding response: %s", err.Error()),
			http.StatusInternalServerError,
		)
	}
}

func makeUrl(s string) string {
	return fmt.Sprintf("<a href=\"%s\">%s</a>", s, s)
}

func makeHtml(s string) string {
	return fmt.Sprintf("<!DOCTYPE html><html><body>%s</body></html>", s)
}

func makeList(rows []string) string {
	var contents string
	for _, row := range rows {
		contents += fmt.Sprintf("<li>%s</li>", row)
	}
	list := fmt.Sprintf("<ul>%s</ul>", contents)
	return list
}

func mustRespond(response http.ResponseWriter, msg string) {
	_, err := io.WriteString(response, msg)
	mustOk(err)
}

func mustOk(err error) {
	if err != nil {
		panic(err)
	}
}

type User struct {
	FirstName string
	LastName  string
	Email     string
}

func makeUser(names []byte) User {
	fields := strings.Split(string(names), "&")
	firstName := strings.Split(fields[0], "=")[1]
	lastName := strings.Split(fields[1], "=")[1]
	email := strings.Split(fields[2], "=")[1]

	return User{
		FirstName: firstName,
		Email:     email,
		LastName:  lastName,
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	err := os.Remove("log.txt")
	if err != nil {
		log.Fatal(err)
	}

	logfile, err := os.OpenFile(
		"log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755,
	)
	if err != nil {
		log.Fatalf("error creating log file: %s", err.Error())
	}
	log.SetOutput(logfile)

	http.HandleFunc("/", index)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/register", register)
	http.HandleFunc("/register_click", register_click)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("error serving server: %s\n", err.Error())
	}
}

func index(res http.ResponseWriter, _ *http.Request) {
	sites := []string{makeUrl("/hello"), makeUrl("/echo"), makeUrl("/register")}
	msg := fmt.Sprintf("%s\n", makeHtml(makeList(sites)))
	_, err := io.WriteString(res, msg)
	if err != nil {
		errmsg := fmt.Sprintf("error writing response: %s\n", err.Error())
		http.Error(res, errmsg, http.StatusInternalServerError)
		log.Fatalf(errmsg)
		fmt.Println(errmsg)
	}
}

func echo(res http.ResponseWriter, req *http.Request) {
	msg := fmt.Sprintf("%#v", req)
	_, err := io.WriteString(res, msg)
	if err != nil {
		errmsg := fmt.Sprintf("error writing response: %s\n", err.Error())
		http.Error(res, errmsg, http.StatusInternalServerError)
		log.Fatalf(errmsg)
		fmt.Println(errmsg)
	}
}

func hello(res http.ResponseWriter, _ *http.Request) {
	_, err := io.WriteString(res, "Hello, World!")
	if err != nil {
		errmsg := fmt.Sprintf("error writing response: %s\n", err.Error())
		http.Error(res, errmsg, http.StatusInternalServerError)
		log.Fatalf(errmsg)
		fmt.Println(errmsg)
	}
}

func register(res http.ResponseWriter, req *http.Request) {
	msg := `
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
</html>`

	_, err := io.WriteString(res, msg)
	if err != nil {
		errmsg := fmt.Sprintf("error writing response: %s\n", err.Error())
		http.Error(res, errmsg, http.StatusInternalServerError)
		log.Fatalf(errmsg)
		fmt.Println(errmsg)
	}
}

func register_click(res http.ResponseWriter, req *http.Request) {
	buf, err := io.ReadAll(req.Body)
	if err != nil {
		errmsg := fmt.Sprintf("error reading request: %s", err.Error())
		http.Error(
			res,
			errmsg,
			http.StatusInternalServerError,
		)
		log.Fatal(errmsg)
		fmt.Println(errmsg)
	}
	if len(buf) == 0 {
		errmsg := "error: no request data"
		http.Error(
			res, errmsg, http.StatusInternalServerError,
		)
		log.Fatal(errmsg)
		fmt.Println(errmsg)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(res).Encode(makeUser(buf))
	if err != nil {
		errmsg := fmt.Sprintf("error encoding response: %s", err.Error())
		http.Error(
			res,
			errmsg,
			http.StatusInternalServerError,
		)
		log.Fatal(errmsg)
		fmt.Println(errmsg)
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

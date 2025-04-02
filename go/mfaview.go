package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"mfaviewresource"
)

// Function to create table with hyperlink button and message
func messageTable(w http.ResponseWriter, urlPath string, buttonMessage string, descriptionMessage string) {
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table class=defaultTableColor>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th><a href=\""+urlPath+"\" class=\"pageButton\">"+buttonMessage+"</a></th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>"+descriptionMessage+"&#11107</th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<br>")
}

// Function to add HTML form authentication input section
func inputAuth(w http.ResponseWriter) {
	fmt.Fprintf(w, "  <label for=\"email\"><b>Enter Email Address:</b>")
	fmt.Fprintf(w, "  </label><br>")
	fmt.Fprintf(w, "  <input type=\"email\" id=\"email\" name=\"email\">")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "  <label for=\"password\"><b>Enter Password:</b>")
	fmt.Fprintf(w, "  </label><br>")
	fmt.Fprintf(w, "  <input type=\"password\" id=\"password\" name=\"password\">")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "  <label for=\"2FA\"><b>Enter 2FA Code:</b>")
	fmt.Fprintf(w, "  </label><br>")
	fmt.Fprintf(w, "  <input type=\"text\" id=\"2FA\" name=\"2FA\">")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<br>")
}

// Function to validate user input
func validateInput() {


}

// Function to create red warning box around message argument
func textBox(w http.ResponseWriter, textBoxMessage string) {
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<div class=leftToCentre>")
	fmt.Fprintf(w, "  <table class=errorTableColor>")
	fmt.Fprintf(w, "    <tr>")
	fmt.Fprintf(w, "      <th>")
	fmt.Fprintf(w, textBoxMessage)
	fmt.Fprintf(w, "      </th>")
	fmt.Fprintf(w, "    </tr>")
	fmt.Fprintf(w, "  </table>")
	fmt.Fprintf(w, "</div>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<br>")
}

var startHTML string
var endHTML string

func main() {
	err := godotenv.Load("/usr/local/etc/mfaview/env/mfaview.env")
	if err != nil {
		panic("Error loading mfaview.env file")
	}

	email := os.Getenv("email")
	password := os.Getenv("password")
	address := os.Getenv("address")
	port := os.Getenv("port")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
		}
		startHTML = mfaviewresource.StartHTML()
		endHTML = mfaviewresource.EndHTML()		

		fmt.Fprintf(w, startHTML)
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table class=defaultTableColor>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><a href=\"https://ell.today\" class=\"tableButton\">Written by Elliot Keavney (Website)</a></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><a href=\"https://github.com/Ellwould/mfa-view\" class=\"tableButton\">MFA View Source Code (GitHub)</a></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		messageTable(w, "/add-mfa-account", "Click to add a new MFA account", "Enter email, password & 2FA to<br>generate MFA codes from key")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<div class=login>")
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/\">")
		inputAuth(w)
		fmt.Fprintf(w, "  <input type=\"submit\" value=\"Submit\">")
		fmt.Fprintf(w, "</form>")
		fmt.Fprintf(w, "</div")

		//Get email address and validate
		f1 := r.FormValue("email")
		var inputEmail string
		inputEmail = f1
		validateInputEmail := validator.New()
		validateInputEmailErr := validateInputEmail.Var(inputEmail, "email,required,max=100")

		//Get password and validate
		f2 := r.FormValue("password")
		var inputPassword string
		inputPassword = f2
		validateInputPassword := validator.New()
		validateInputPasswordErr := validateInputPassword.Var(inputPassword, "required,min=20,max=100")

		//Get 2FA and validate
		f3 := r.FormValue("2FA")
		var input2fa string
		input2fa = f3
		validateInput2fa := validator.New()
		validateInput2faErr := validateInput2fa.Var(input2fa, "number,required,min=6,max=6")

		// Need to work on validation more
		if inputEmail == "" && inputPassword == "" && input2fa == "" {
		} else if validateInputEmailErr != nil {
			textBox(w, "Please enter a valid email address, max 100 charecters length")
		} else if validateInputPasswordErr != nil {
			textBox(w, "Password needs to be between 20-100 charecters length")
		} else if validateInput2faErr != nil {
			textBox(w, "MFA code needs to be a 6 digit number")
		} else if inputEmail == email && inputPassword == password && input2fa == "111111" {
			textBox(w, "Correct Credentials Entered")
		} else {
			textBox(w, "Wrong Credentials Entered")
		}
		fmt.Fprint(w, endHTML)
	})

	http.HandleFunc("/add-mfa-account", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
		}
		startHTML = mfaviewresource.StartHTML()
		endHTML = mfaviewresource.EndHTML()
		
		fmt.Fprintf(w, startHTML)
		messageTable(w, "/", "Click to login & view MFA account(s)", "Enter email, password, 2FA,<br>algorithm, duration & MFA secret<br>below to add account")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<div class=login>")
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/add-mfa-account\">")
		inputAuth(w)
		fmt.Fprintf(w, "  <input type=\"submit\" value=\"Submit\">")
		fmt.Fprintf(w, "</form>")
		fmt.Fprintf(w, "</div")
		fmt.Fprintf(w, endHTML)
	})

	socket := address+":"+port
	fmt.Println("MFA View is running on: " + socket)

	// Start server on port specified above
	log.Fatal(http.ListenAndServe(socket, nil))
}

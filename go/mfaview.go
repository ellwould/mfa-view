package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"mfaviewresource"
	"time"
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
	fmt.Fprintf(w, "    <th>"+descriptionMessage+"</th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<br>")
}

// Function to add HTML form authentication input section
func inputForm(w http.ResponseWriter, value string, labelMessage string, inputType string) {
	fmt.Fprintf(w, "  <label for=\""+value+"\"><b>Enter "+labelMessage+":</b>")
	fmt.Fprintf(w, "  </label><br>")
	fmt.Fprintf(w, "  <input type=\""+inputType+"\" id=\""+value+"\" name=\""+value+"\">")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<br>")
}

// Function to validate user input
func validateInput(formValue string, valueType string, minNum string, maxNum string) (validation bool) {
	validateInput := validator.New()
	if valueType == "" {
		validateInputErr := validateInput.Var(formValue, "required,min="+minNum+",max="+maxNum)
		if validateInputErr != nil {
			validation = false
			return
	        } else {
	        	validation = true
	        	return
		}        
	} else {
		validateInputErr := validateInput.Var(formValue, valueType+",required,min="+minNum+",max="+maxNum)
		if validateInputErr != nil {
			validation = false
			return
	        } else {
	        	validation = true
	        	return
		}
	}
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

//	email := os.Getenv("email")
//	password := os.Getenv("password")
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
		currentTime := time.Now().Format("15:04:05")
		messageTable(w, "/add-mfa-account", "Click to add a new MFA account", "MFA code(s) expire on the minute<br>&#128347 & 30 seconds past a minute &#128353<br><br>Time at page load = "+currentTime+"<br><br>Enter email, password & 2FA to<br>generate MFA code(s) from key(s)")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<div class=login>")
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/\">")
		inputForm(w, "email", "Email Address", "email",)
		inputForm(w, "password", "Password", "password",)
		inputForm(w, "2FA", "2FA Code", "text",)
		fmt.Fprintf(w, "  <input type=\"submit\" value=\"Submit\">")
		fmt.Fprintf(w, "</form>")
		fmt.Fprintf(w, "</div")

		// Get email address, password and 2FA from HTTP POST
		inputEmail := r.FormValue("email")
		inputPassword := r.FormValue("password")
		input2fa := r.FormValue("2FA")
		
		// Validate email address, password and 2FA
		validationEmail := validateInput(inputEmail, "email", "6", "320")
		validationPassword := validateInput(inputPassword, "", "20", "100")
		validation2fa := validateInput(input2fa, "number", "6", "6")
		
		// Need to work on validation more
		if inputEmail == "" && inputPassword == "" && input2fa == "" {
		} else if validationEmail == false {
			textBox(w, "Please enter a valid email address, max 100 charecters length")
		} else if validationPassword == false {
			textBox(w, "Password needs to be between 20-100 charecters length")
		} else if validation2fa == false {
			textBox(w, "MFA code needs to be a 6 digit number")
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
		messageTable(w, "/", "Click to login & view MFA account(s)", "Enter email, password, MFA account<br>name, MFA secret key, select SHA &<br>enter 2FA code to add a new account")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<div class=login>")
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/add-mfa-account\">")
		inputForm(w, "email", "Email Address", "email",)
		inputForm(w, "password", "Password", "password",)		
		inputForm(w, "account", "New MFA Account Name", "text",)
		inputForm(w, "MFA", "New MFA Secret Key", "text",)
		fmt.Fprintf(w, "  <label for=\"sha\"><b>Secure Hash Algorithm (SHA):</b>")
		fmt.Fprintf(w, "  </label><br>")
		fmt.Fprintf(w, "  <select id=\"sha\" name=\"sha\">")
		var shaArray = [3]string{"sha1", "sha256", "sha512"}				
		for i := 0; i < len(shaArray); i++ {
			fmt.Fprintf(w, "<option value="+shaArray[i]+">"+shaArray[i]+"</option>")
		}
		fmt.Fprintf(w, "  </select>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<br>")
		inputForm(w, "2FA", "2FA Code", "text",)
		fmt.Fprintf(w, "  <input type=\"submit\" value=\"Submit\">")
		fmt.Fprintf(w, "</form>")
		fmt.Fprintf(w, "</div")
		
		// Get email address, password and 2FA from HTTP POST
                inputEmail := r.FormValue("email")
                inputPassword := r.FormValue("password")
                input2fa := r.FormValue("2FA")
                
                // Validate email address, password and 2FA
                validationEmail := validateInput(inputEmail, "email", "6", "320")
                validationPassword := validateInput(inputPassword, "", "20", "100")
                validation2fa := validateInput(input2fa, "number", "6", "6")		

                if inputEmail == "" && inputPassword == "" && input2fa == "" {
                } else if validationEmail == false {
                	textBox(w, "Please enter a valid email address, max 100 charecters length")
		} else if validationPassword == false {
		        textBox(w, "Password needs to be between 20-100 charecters length")
		} else if validation2fa == false {
                        textBox(w, "MFA code needs to be a 6 digit number")
		} else {
			textBox(w, "Wrong Credentials Entered")
		}
		fmt.Fprintf(w, endHTML)
	})

	socket := address+":"+port
	fmt.Println("MFA View is running on: " + socket)

	// Start server on port specified above
	log.Fatal(http.ListenAndServe(socket, nil))
}

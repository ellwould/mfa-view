package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"mfaviewresource"
	"net/http"
	"os"
	"slices"
	"time"
)

// American National Standards Institute (ANSI) reset colour code
var resetColour = "\033[0m"

// American National Standards Institute (ANSI) text colour codes
var textBoldWhite = "\033[1;37m"
var textBoldBlack = "\033[1;30m"

// American National Standards Institute (ANSI) background colour codes
var bgRed = "\033[41m"
var bgGreen = "\033[42m"
var bgYellow = "\033[43m"

// Clear screen function for GNU/Linux OS's
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

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

// Function to generate new hash and salted password using Bcrypt
func genPasswd(passwd []byte) string {
	hashAndSalt, err := bcrypt.GenerateFromPassword(passwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hashAndSalt)
}

// Recursive function to add a user
// Needs alot of work this function:
func createUser() {
	clearScreen()
	fmt.Println(bgGreen + textBoldWhite)
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(" □                                                       □ ")
	fmt.Println(" □               email address is required               □ ")
	fmt.Println(" □                                                       □ ")
	fmt.Println(" □   email address can be changed later in mfaview.env   □ ")
	fmt.Println(" □                                                       □ ")
	fmt.Println(" □            (Type \"exit\" to quit program)            □ ")
	fmt.Println(" □                                                       □ ")
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(resetColour)
	fmt.Println(bgRed + textBoldWhite)
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(" □                                                                                                               □ ")
	fmt.Println(" □                                            A password is required                                             □ ")
	fmt.Println(" □                                                                                                               □ ")
	fmt.Println(" □                           The password stored on the disk will be hashed and salted                           □ ")
	fmt.Println(" □   The password will be used to encrypt & decrypt Multi-factor authentication (MFA) secret keys on the disk.   □ ")
	fmt.Println(" □     A lost password will result in unreadable Multi-factor authentication (MFA) secret keys on the disk!      □ ")
	fmt.Println(" □      Store the password securly, do not lose the password, there is no way to recover a lost password!!       □ ")
	fmt.Println(" □                                                                                                               □ ")
	fmt.Println(" □                                        (Type \"exit\" to quit program)                                        □ ")
	fmt.Println(" □                                                                                                               □ ")
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(resetColour)
	fmt.Println("")
	fmt.Println(bgYellow + textBoldBlack)
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(" □                                                                                       □ ")
	fmt.Println(" □   Add the two-factor authentication (2FA) secret key generated below to another app   □ ")
	fmt.Println(" □                                                                                       □ ")
	fmt.Println(" □      MFA View aims to be as secure as possible it requires another app to provide     □ ")
	fmt.Println(" □                    a two-factor authentication (2FA) code to login.                   □ ")
	fmt.Println(" □   The two-factor authentication (2FA) secret key can be changed later in mfaview.env  □ ")
	fmt.Println(" □                                                                                       □ ")
	fmt.Println(" □                            (Type \"exit\" to quit program)                            □ ")
	fmt.Println(" □                                                                                       □ ")
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(resetColour)
	fmt.Scan()

}

// Recursive function to change an existing known password
func changePassword() {
	clearScreen()
	fmt.Println(bgRed + textBoldWhite)
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(" □                                                                   □ ")
	fmt.Println(" □      The change password option is set to yes in mfaview.env      □ ")
	fmt.Println(" □                                                                   □ ")
	fmt.Println(" □    To change the password the exisitng password must be known     □ ")
	fmt.Println(" □                                                                   □ ")
	fmt.Println(" □                  (Type \"exit\" to quit program)                  □ ")
	fmt.Println(" □                                                                   □ ")
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(resetColour)
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
	addAccount := os.Getenv("add_account")

	if email == "" && password == "" {
		createUser()
	} else {

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
			messageTable(w, "/add-account", "Click to add a new MFA account", "MFA code(s) expire on the minute<br>&#128347 & 30 seconds past a minute &#128353<br><br>Time at page load = "+currentTime+"<br><br>Enter email, password & 2FA to<br>generate MFA code(s) from key(s)")
			fmt.Fprintf(w, "<br>")
			fmt.Fprintf(w, "<div class=login>")
			fmt.Fprintf(w, "<form method=\"POST\" action=\"/\">")
			inputForm(w, "email", "Email Address", "email")
			inputForm(w, "password", "Password", "password")
			inputForm(w, "2FA", "2FA Code", "text")
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
				// Testing genPassword function
				p := genPasswd([]byte("p"))
				textBox(w, p)
			}
			fmt.Fprint(w, endHTML)
		})

		if addAccount == "yes" || addAccount == "Yes" || addAccount == "YES" {

			http.HandleFunc("/add-account", func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					fmt.Fprintf(w, "ParseForm() err: %v", err)
				}
				startHTML = mfaviewresource.StartHTML()
				endHTML = mfaviewresource.EndHTML()

				fmt.Fprintf(w, startHTML)
				messageTable(w, "/", "Click to login & view MFA account(s)", "Enter email, password, MFA account<br>name, MFA secret key, select SHA &<br>enter 2FA code to add a new account")
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div class=login>")
				fmt.Fprintf(w, "<form method=\"POST\" action=\"/add-account\">")
				inputForm(w, "email", "Email Address", "email")
				inputForm(w, "password", "Password", "password")
				inputForm(w, "account", "New MFA Account Name", "text")
				inputForm(w, "MFA", "New MFA Secret Key", "text")
				fmt.Fprintf(w, "  <label for=\"sha\"><b>Secure Hash Algorithm (SHA):</b>")
				fmt.Fprintf(w, "  </label><br>")
				fmt.Fprintf(w, "  <select id=\"sha\" name=\"sha\">")
				var shaList = []string{"sha1", "sha256", "sha512"}
				for i := 0; i < len(shaList); i++ {
					fmt.Fprintf(w, "<option value="+shaList[i]+">"+shaList[i]+"</option>")
				}
				fmt.Fprintf(w, "  </select>")
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<br>")
				inputForm(w, "2FA", "2FA Code", "text")
				fmt.Fprintf(w, "  <input type=\"submit\" value=\"Submit\">")
				fmt.Fprintf(w, "</form>")
				fmt.Fprintf(w, "</div")

				// Get email address, password and 2FA from HTTP POST
				inputEmail := r.FormValue("email")
				inputPassword := r.FormValue("password")
				inputAccount := r.FormValue("account")
				inputMfa := r.FormValue("MFA")
				inputSha := r.FormValue("sha")
				input2fa := r.FormValue("2FA")

				// Validate email address, password and 2FA
				validationEmail := validateInput(inputEmail, "email", "6", "320")
				validationPassword := validateInput(inputPassword, "", "20", "100")
				validationAccount := validateInput(inputAccount, "", "20", "100")
				validationMfa := validateInput(inputMfa, "", "10", "200")
				validationSha := slices.Contains(shaList, inputSha)
				validation2fa := validateInput(input2fa, "number", "6", "6")

				if inputEmail == "" && inputPassword == "" && inputAccount == "" && inputMfa == "" && inputSha == "" && input2fa == "" {
				} else if validationEmail == false {
					textBox(w, "Please enter a valid email address, max 100 charecters length")
				} else if validationPassword == false {
					textBox(w, "Password needs to be between 20-100 charecters length")
				} else if validationAccount == false {
					textBox(w, "Please enter a valid account name, max 100 charecters length")
				} else if validationMfa == false {
					textBox(w, "New MFA secret key needs to be between 10-200 charecters ")
				} else if validationSha == false {
					textBox(w, "Secure Hash Algorithm can be sha1, sha256 or sha512")
				} else if validation2fa == false {
					textBox(w, "2FA code needs to be a 6 digit number")
				} else {
					textBox(w, "Wrong Credentials Entered")
				}
				fmt.Fprintf(w, endHTML)
			})

		} else {
			http.HandleFunc("/add-account", func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					fmt.Fprintf(w, "ParseForm() err: %v", err)
				}
				startHTML = mfaviewresource.StartHTML()
				endHTML = mfaviewresource.EndHTML()

				fmt.Fprintf(w, startHTML)
				messageTable(w, "/", "Click to login & view MFA account(s)", "Adding a new MFA account is<br>switched off in the configuration file")
				fmt.Fprintf(w, endHTML)
			})
		}
	}

	socket := address + ":" + port
	fmt.Println("MFA View is running on: " + socket)

	// Start server on port specified above
	log.Fatal(http.ListenAndServe(socket, nil))
}

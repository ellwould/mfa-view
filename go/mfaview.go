package main

import (
	"bytes"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/mdp/qrterminal/v3"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"log"
	"mfaviewresource"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Variable for mfaview.env absolute path
var mfaViewEnv string = "/usr/local/etc/mfaview/env/mfaview.env"

// Variable for American National Standards Institute (ANSI) reset colour code
var resetColour = "\033[0m"

// Variables for American National Standards Institute (ANSI) text colour codes
var textBoldWhite = "\033[1;37m"
var textBoldBlack = "\033[1;30m"
var textBoldGreen = "\033[1;32m"

// Variables for American National Standards Institute (ANSI) background colour codes
var bgBlack = "\033[40m"
var bgRed = "\033[41m"
var bgGreen = "\033[42m"
var bgYellow = "\033[43m"
var bgBlue = "\033[44m"
var bgPurple = "\033[45m"
var bgCyan = "\033[46m"
var bgWhite = "\033[47m"

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

// Function to validate user input utlising the Go Validator package
func validateInput(formValue string, valueType string) (validation bool) {
	validateInput := validator.New()
	// Conditional statments are used for each type of value inputted from a user
	if valueType == "email" {
		validateInputErr := validateInput.Var(formValue, "email,required,min=6,max=320")
		if validateInputErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "password" {
		validateInputErr := validateInput.Var(formValue, "required,min=20,max=100")
		if validateInputErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "2FA" {
		validateInputErr := validateInput.Var(formValue, "number,required,min=6,max=6")
		if validateInputErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "text" {
		validateInputErr := validateInput.Var(formValue, "required,min=10,max=200")
		if validateInputErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	}  else if valueType == "IP" {
		validateInputErr := validateInput.Var(formValue, "required,ip_addr")
		if validateInputErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}                
	} else {
		validation = false
		return
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

// Function to draw box with squares around message, must have a message with characters that total a odd number
func messageBoxCli(bgColour string, messageColour string, message string) {
	topBottomSquare := strings.Repeat(" □", (len(message)/2)+6)
	inbetweenSpace := strings.Repeat(" ", len(message)+8)
	fmt.Println(bgColour + messageColour)
	fmt.Println(topBottomSquare + " ")
	fmt.Println(" □" + inbetweenSpace + "□ ")
	fmt.Println(" □    " + message + "    □ ")
	fmt.Println(" □" + inbetweenSpace + "□ ")
	fmt.Println(topBottomSquare + " ")
	fmt.Print(resetColour)
}

// Function to display message on CLI informing the user the configuration file has a wrong value
func invalidEnvCli(message string) {
	clearScreen()
	messageBoxCli(bgRed, textBoldWhite, message)
	os.Exit(0)
}

// Function to terminate the program with exit code 0 to indicate there were no errors
func exitProgramCli() {
	clearScreen()
	messageBoxCli(bgBlue, textBoldWhite, "Program exited.")
	os.Exit(0)
}

// Function to inform the user of wrong input and waits for user to press enter/return key
func invalidInputCli() {
	clearScreen()
	messageBoxCli(bgRed, textBoldWhite, "Invalid input press enter/return to continue.")
	fmt.Scanln()
}

// Function to inform the user to type exit to terminate the program
func typeExitCli() {
	messageBoxCli(bgBlue, textBoldWhite, "(Type exit to quit program)")
}

// Function to replace strings in mfaview.conf
func replaceText(oldText string, newText string) {
	readFile, err := os.ReadFile(mfaViewEnv)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	writeFile := bytes.Replace(readFile, []byte(oldText), []byte(newText), -1)
	if err = os.WriteFile(mfaViewEnv, writeFile, 0666); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

// Function to generate a 2FA key and URL
func gen2faKey() (string, string) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "MFA View Server",
		AccountName: "MFA View Server Account",
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})

	if err != nil {
		panic(err)
	}

	return key.Secret(), key.URL()
}

// A recursive function to add a user email
func addEmailCli() {
	clearScreen()
	typeExitCli()
	messageBoxCli(bgCyan, textBoldWhite, "Email address is required")
	messageBoxCli(bgRed, textBoldWhite, "Email address can be manually changed later in "+mfaViewEnv)
	fmt.Println(textBoldBlack)
	fmt.Printf("   Please enter an email address between 6-320\n   characters: ")
	var addEmail string
	fmt.Scan(&addEmail)
	fmt.Println(resetColour)
	validationAddEmail := validateInput(addEmail, "email")
	if addEmail == "exit" || addEmail == "Exit" || addEmail == "EXIT" {
		exitProgramCli()
	} else if validationAddEmail == false {
		invalidInputCli()
		addEmailCli()
	} else {
		replaceText("email_not_set", addEmail)
	}
}

// A recursive function to add a user password, a hashed and salted value of the password will be stored
func addPasswordCli() {
	clearScreen()
	typeExitCli()
	messageBoxCli(bgPurple, textBoldWhite, "A password is required ")
	fmt.Println(bgRed + textBoldWhite)
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(" □                                                                                                               □ ")
	fmt.Println(" □                           The password stored on the disk will be hashed and salted                           □ ")
	fmt.Println(" □   The password will be used to encrypt & decrypt Multi-factor authentication (MFA) secret keys on the disk.   □ ")
	fmt.Println(" □     A lost password will result in unreadable Multi-factor authentication (MFA) secret keys on the disk!      □ ")
	fmt.Println(" □      Store the password securly, do not lose the password, there is no way to recover a lost password!!       □ ")
	fmt.Println(" □                                                                                                               □ ")
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Print(resetColour)
	fmt.Println(textBoldBlack)
	fmt.Printf("    Please enter a password between 20-100\n    characters: ")
	var addPassword string
	fmt.Scan(&addPassword)
	if addPassword == "exit" || addPassword == "Exit" || addPassword == "EXIT" {
		exitProgramCli()
	}        
	validationAddPassword := validateInput(addPassword, "password")
	fmt.Println("")
	fmt.Printf("    Please re-enter the password: ")
	var checkPassword string
	fmt.Scan(&checkPassword)
	validationCheckPassword := validateInput(checkPassword, "password")
	fmt.Println(resetColour)
	if validationAddPassword == false || validationCheckPassword == false {
		invalidInputCli()
		addPasswordCli()
	} else if addPassword != checkPassword {
		clearScreen()
		messageBoxCli(bgRed, textBoldWhite, "Passwords entered do not match, press enter/return to continue.")
		fmt.Scanln()
		addPasswordCli()
	} else {
		hashedPassword := genPasswd([]byte(addPassword))
		replaceText("password_not_set", hashedPassword)
	}
}

// A recursive function to add a 2FA secret key, MFA View needs a username(email), password and 2FA code to login to view other 2FA/MFA codes
func add2faCli() {
	clearScreen()
	typeExitCli()
	messageBoxCli(bgBlack, textBoldWhite, "Scan the QR (Quick Response) code with another authenticator app or copy the secret key to another authenticator app ")
	fmt.Println("")
	generated2faKey, generated2faUrl := gen2faKey()
	qrConfig := qrterminal.Config{
		Level:     qrterminal.L,
		Writer:    os.Stdout,
		BlackChar: qrterminal.WHITE,
		WhiteChar: qrterminal.BLACK,
		QuietZone: 0,
	}
	qrterminal.GenerateWithConfig(generated2faUrl, qrConfig)
	fmt.Println(bgBlack + textBoldWhite)
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(" □                                                       □ ")
	fmt.Println(" □    Secret Key: " + resetColour + bgWhite + textBoldBlack + " " + generated2faKey + " " + resetColour + bgBlack + textBoldWhite + "     □ ")
	fmt.Println(" □                                                       □ ")
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(resetColour + bgRed + textBoldWhite)
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(" □                                                                                                   □ ")
	fmt.Println(" □                  MFA View aims to be as secure as possible, it requires another                   □ ")
	fmt.Println(" □           authenticator app to provide a 2FA (Two-Factor Authentication) code to login.           □ ")
	fmt.Println(" □    The 2FA secret key can be manually changed later in " + mfaViewEnv + ".    □ ")
	fmt.Println(" □                                                                                                   □ ")
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(resetColour + textBoldBlack)
	fmt.Printf("      Please enter the 6 digit generated\n      code from the your other 2FA/MFA app: ")
	var test2faCode string
	fmt.Scan(&test2faCode)
	validationTest2faCode := validateInput(test2faCode, "2FA")
	fmt.Println("")
	if test2faCode == "exit" || test2faCode == "Exit" || test2faCode == "EXIT" {
		exitProgramCli()
	} else if validationTest2faCode == false {
		invalidInputCli()
		add2faCli()
	} else {
		correct2faCode := totp.Validate(test2faCode, generated2faKey)
		if correct2faCode {
			clearScreen()
			replaceText("2fa_not_set", generated2faKey)
			messageBoxCli(bgGreen, textBoldWhite, "Successfully added 2FA secret key")
			os.Exit(0)
		} else {
			clearScreen()
			messageBoxCli(bgRed, textBoldWhite, "2FA code entered is incorrect, press enter/return to continue")
			fmt.Scanln()                
			add2faCli()
		}
	}
}

// Basic function with no parameters to call the addEmailCli, addPasswordCli and add2faCli functions
func createUserCli() {
	addEmailCli()
	addPasswordCli()
	add2faCli()
}

// Basic function with no parameters to call the addPasswordCli and add2faCli functions
func createPassword2faCli() {
	addPasswordCli()
	add2faCli()
}

// Recursive function to change an existing known password
// This function needs alot of work:
func changePasswordCli() {
	clearScreen()
	fmt.Println(bgRed + textBoldWhite)
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(" □                                                                   □ ")
	fmt.Println(" □      The change_password option is set to yes in mfaview.env      □ ")
	fmt.Println(" □                                                                   □ ")
	fmt.Println(" □    To change the password the exisitng password must be known     □ ")
	fmt.Println(" □                                                                   □ ")
	fmt.Println(" □                  (Type \"exit\" to quit program)                    □ ")
	fmt.Println(" □                                                                   □ ")
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(resetColour)
}

// Declare two string variables for HTML
var startHTML string
var endHTML string

// Main function
func main() {
	err := godotenv.Load(mfaViewEnv)
	if err != nil {
		panic("Error loading mfaview.env file")
	}

	envEmail := os.Getenv("email")
	envPassword := os.Getenv("password")
	env2faKey := os.Getenv("2fa_key")
	envChangePassword := os.Getenv("change_password")
	envAddress := os.Getenv("address")
	envPort := os.Getenv("port")
	envAddAccount := os.Getenv("add_account_page")

	validationEnvEmail := validateInput(envEmail, "email")
	validationEnvAddress := validateInput(envAddress, "IP")

	envPortInt, err := strconv.Atoi(envPort)
	if err != nil {
		invalidEnvCli("Port must be a number in " + mfaViewEnv)
	}

	if envEmail == "email_not_set" && envPassword == "password_not_set" && env2faKey == "2fa_not_set" {
		createUserCli()
	} else if validationEnvEmail == false {
		invalidEnvCli("Email address stored in " + mfaViewEnv + " is invalid")
	} else if validationEnvEmail == true && envPassword == "password_not_set" && env2faKey == "2fa_not_set" {
		createPassword2faCli()
	} else if env2faKey == "2fa_not_set" {
		add2faCli()	
	} else if envAddress != "localhost" {
		if validationEnvAddress == false {
			invalidEnvCli("Address in " + mfaViewEnv + " must be a valid Internet Protocol (IP) address or localhost")
		}
	} else if envPortInt <= 0 || envPortInt >= 65536 {
		invalidEnvCli("Port number in " + mfaViewEnv + " must be between 1 and 35535")
	} else if envChangePassword == "yes" || envChangePassword == "Yes" || envChangePassword == "YES" {
		changePasswordCli()
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
			fmt.Fprintf(w, "<div>")
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
			validationEmail := validateInput(inputEmail, "email")
			validationPassword := validateInput(inputPassword, "password")
			validation2fa := validateInput(input2fa, "2FA")

			// Conditional statement to validate input
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

		// Conditional statement to activate the add account section of the website based on value inside the
		// mfaview.env configuration file
		if envAddAccount == "yes" || envAddAccount == "Yes" || envAddAccount == "YES" {

			http.HandleFunc("/add-account", func(w http.ResponseWriter, r *http.Request) {
				if err := r.ParseForm(); err != nil {
					fmt.Fprintf(w, "ParseForm() err: %v", err)
				}
				startHTML = mfaviewresource.StartHTML()
				endHTML = mfaviewresource.EndHTML()

				fmt.Fprintf(w, startHTML)
				messageTable(w, "/", "Click to login & view MFA account(s)", "Enter email, password, MFA account<br>name, MFA secret key, select SHA &<br>enter 2FA code to add a new account")
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div>")
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

				// Get email address, password, MFA account name, MFA secret key, SHA type and 2FA from HTTP POST
				inputEmail := r.FormValue("email")
				inputPassword := r.FormValue("password")
				inputAccount := r.FormValue("account")
				inputMfa := r.FormValue("MFA")
				inputSha := r.FormValue("sha")
				input2fa := r.FormValue("2FA")

				// Returns false or true based on if input is valid
				validationEmail := validateInput(inputEmail, "email")
				validationPassword := validateInput(inputPassword, "password")
				validationAccount := validateInput(inputAccount, "text")
				validationMfa := validateInput(inputMfa, "text")
				validationSha := slices.Contains(shaList, inputSha)
				validation2fa := validateInput(input2fa, "2FA")

				// Conditional statement to validate input
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
			// HTTP server informing the user add account is not available
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

	// IP address and port number, value taken from mfaview.env
	socket := envAddress + ":" + envPort
	messageBoxCli(bgGreen, textBoldWhite, "MFA View is running on "+socket)
	fmt.Println("")

	// Start server on port specified above
	log.Fatal(http.ListenAndServe(socket, nil))
}

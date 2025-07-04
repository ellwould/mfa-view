package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ellwould/aestext"
	"github.com/ellwould/csvcell"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/mdp/qrterminal/v3"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Constant for mfaview.env absolute path
const mfaViewEnv string = "/etc/mfaview/env/mfaview.env"

// Constant for directory path that contains the file key.csv
const dirKeyCSV string = "/etc/mfaview/key"

// Constant for key.csv file
const fileKeyCSV string = "key.csv"

// Constant for directory path that contains the files mfaview-start.html and mfaview-end.html
const dirHTML string = "/etc/mfaview/html"

// Constant for fileStartHTML file
const fileStartHTML string = "mfaview-start.html"

// Constant for fileEndHTML file
const fileEndHTML string = "mfaview-end.html"

// Variable for American National Standards Institute (ANSI) reset colour code
const resetColour string = "\033[0m"

// Variables for American National Standards Institute (ANSI) text colour codes
const (
	textBoldWhite string = "\033[1;37m"
	textBoldBlack string = "\033[1;30m"
	textBoldGreen string = "\033[1;32m"
)

// Variables for American National Standards Institute (ANSI) background colour codes

const (
	bgBlack  string = "\033[40m"
	bgRed    string = "\033[41m"
	bgGreen  string = "\033[42m"
	bgYellow string = "\033[43m"
	bgBlue   string = "\033[44m"
	bgPurple string = "\033[45m"
	bgCyan   string = "\033[46m"
	bgWhite  string = "\033[47m"
)

// Clear screen function for GNU/Linux OS's
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// Function to retrieve file contents
func getFile(rootDirPath string, fileName string) (result string) {

	// Go introduced OpenRoot in version 1.24, it restricts file operations to a single directory
	rootDir, err := os.OpenRoot(rootDirPath)

	if err != nil {
		panic("Directory path does not exist")
	}

	defer rootDir.Close()

	file, err := rootDir.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var fileSlice []string

	for scanner.Scan() {
		line := scanner.Text()
		fileSlice = append(fileSlice, line)
	}

	result = strings.Join(fileSlice, "")
	return result
}

// Function to empty a file
func wipeKeyFile() {

	// Go introduced OpenRoot in version 1.24, it restricts file operations to a single directory
	rootDir, err := os.OpenRoot(dirKeyCSV)

	if err != nil {
		panic("Directory path does not exist")
	}

	defer rootDir.Close()

	file, err := rootDir.Create(fileKeyCSV)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
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
func validateInput(value string, valueType string) (validation bool) {
	validateInput := validator.New()
	// Conditional statments are used for each type of value inputted from a user
	if valueType == "email" {
		validateInputErr := validateInput.Var(value, "email,required,min=6,max=320,excludes=0x2C")
		if validateInputErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "password" {
		validateInputErr := validateInput.Var(value, "required,min=16,max=32,excludes=0x2C")
		if validateInputErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "hashedPassword" {
		validateInputErr := validateInput.Var(value, "required,min=60,max=60")
		if validateInputErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "2FACode" {
		validateInputErr := validateInput.Var(value, "number,required,min=6,max=6")
		if validateInputErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "secretKey" {
		validateInputErr := validateInput.Var(value, "required,min=16,max=128,excludes=0x2C")
		if validateInputErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "accountName" {
		validateInputErr := validateInput.Var(value, "required,min=1,max=100,excludes=0x2C")
		if validateInputErr != nil {
			validation = false
			return
		} else {
			validation = true
			return
		}
	} else if valueType == "IP" {
		validateInputErr := validateInput.Var(value, "required,ip_addr")
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

// Function to pad out password with 0's for length between 16 and 32 characters
func zeroPad(password string) string {
	zeroQuantity := 32 - len(password)
	newPassword := password + strings.Repeat("0", zeroQuantity)
	return newPassword
}

// Function to generate new hash and salted password using Bcrypt
func genPassword(password []byte) string {
	hashAndSalt, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hashAndSalt)
}

func comparePassword(password []byte, hashedPassword []byte) bool {
	result := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if result == nil {
		return true
	} else {
		return false
	}
}

// Function to draw box with squares around message, must have a message with characters that total a odd number
func messageBoxCLI(bgColour string, messageColour string, message string) {
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
func invalidEnvCLI(message string) {
	clearScreen()
	messageBoxCLI(bgRed, textBoldWhite, message)
	fmt.Println("")
	os.Exit(0)
}

// Function to terminate the program with exit code 0 to indicate there were no errors
func exitProgramCLI() {
	clearScreen()
	messageBoxCLI(bgBlue, textBoldWhite, "Program exited ")
	fmt.Println("")
	os.Exit(0)
}

// Function to inform the user of wrong input and waits for user to press enter/return key
func invalidInputCLI() {
	clearScreen()
	messageBoxCLI(bgRed, textBoldWhite, "Invalid input press enter/return to continue ")
	fmt.Scanln()
}

// Function to inform the user to type exit to terminate the program
func typeExitCLI() {
	messageBoxCLI(bgBlue, textBoldWhite, "(Type exit to quit program)")
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
func gen2FAKey() (string, string) {
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
func addEmailCLI() {
	clearScreen()
	typeExitCLI()
	messageBoxCLI(bgCyan, textBoldWhite, "Email address is required - no commas (,) allowed")
	messageBoxCLI(bgRed, textBoldWhite, "Email address can be manually changed later in "+mfaViewEnv)
	fmt.Println(textBoldBlack)
	fmt.Printf("   Please enter an email address between 6-320\n   characters: ")
	var addEmail string
	fmt.Scan(&addEmail)
	fmt.Println(resetColour)
	validationAddEmail := validateInput(addEmail, "email")
	if addEmail == "exit" || addEmail == "Exit" || addEmail == "EXIT" {
		exitProgramCLI()
	} else if validationAddEmail == false {
		invalidInputCLI()
		addEmailCLI()
	} else {
		replaceText("email_not_set", addEmail)
	}
}

// A recursive function to add a user password, a hashed and salted value of the password will be stored
func addPasswordCLI() {
	clearScreen()
	typeExitCLI()
	messageBoxCLI(bgPurple, textBoldWhite, "A password is required - no commas (,) allowed ")
	fmt.Println(bgRed + textBoldWhite)
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(" □                                                                                                               □ ")
	fmt.Println(" □                           The password stored on the disk will be hashed and salted                           □ ")
	fmt.Println(" □   The password will be used to encrypt & decrypt Multi-factor authentication (MFA) secret keys on the disk    □ ")
	fmt.Println(" □     A lost password will result in unreadable Multi-factor authentication (MFA) secret keys on the disk!      □ ")
	fmt.Println(" □      Store the password securly, do not lose the password, there is no way to recover a lost password!!       □ ")
	fmt.Println(" □                                                                                                               □ ")
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Print(resetColour)
	fmt.Println(textBoldBlack)
	fmt.Printf("    Please enter a password between 16-32\n    characters: ")
	var addPassword string
	fmt.Scan(&addPassword)
	if addPassword == "exit" || addPassword == "Exit" || addPassword == "EXIT" {
		exitProgramCLI()
	}
	validationAddPassword := validateInput(addPassword, "password")
	fmt.Println("")
	fmt.Printf("    Please re-enter the password: ")
	var checkPassword string
	fmt.Scan(&checkPassword)
	validationCheckPassword := validateInput(checkPassword, "password")
	fmt.Println(resetColour)
	if validationAddPassword == false || validationCheckPassword == false {
		invalidInputCLI()
		addPasswordCLI()
	} else if addPassword != checkPassword {
		clearScreen()
		messageBoxCLI(bgRed, textBoldWhite, "Passwords entered do not match, press enter/return to continue ")
		fmt.Scanln()
		addPasswordCLI()
	} else {
		paddedPassword := zeroPad(addPassword)
		hashedPassword := genPassword([]byte(paddedPassword))
		replaceText("password_not_set", hashedPassword)
	}
}

// A recursive function to add a 2FA secret key, MFA View needs a username(email), password and 2FA code to login to view other 2FA/MFA codes
func add2FACLI() {
	clearScreen()
	typeExitCLI()
	messageBoxCLI(bgBlack, textBoldWhite, "Scan the QR (Quick Response) code with another authenticator app or copy the secret key to another authenticator app ")
	fmt.Println("")
	generated2FAKey, generated2FAURL := gen2FAKey()
	qrConfig := qrterminal.Config{
		Level:     qrterminal.L,
		Writer:    os.Stdout,
		BlackChar: qrterminal.WHITE,
		WhiteChar: qrterminal.BLACK,
		QuietZone: 0,
	}
	qrterminal.GenerateWithConfig(generated2FAURL, qrConfig)
	fmt.Println(bgBlack + textBoldWhite)
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(" □                                                       □ ")
	fmt.Println(" □    Secret Key: " + resetColour + bgWhite + textBoldBlack + " " + generated2FAKey + " " + resetColour + bgBlack + textBoldWhite + "     □ ")
	fmt.Println(" □                                                       □ ")
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(resetColour + bgRed + textBoldWhite)
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(" □                                                                                                   □ ")
	fmt.Println(" □                  MFA View aims to be as secure as possible, it requires another                   □ ")
	fmt.Println(" □           authenticator app to provide a 2FA (Two-Factor Authentication) code to login            □ ")
	fmt.Println(" □    The 2FA secret key can be manually changed later in " + mfaViewEnv + "     □ ")
	fmt.Println(" □                                                                                                   □ ")
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(resetColour + textBoldBlack)
	fmt.Printf("      Please enter the 6 digit generated\n      code from the your other 2FA/MFA app: ")
	var test2FACode string
	fmt.Scan(&test2FACode)
	validationTest2FACode := validateInput(test2FACode, "2FACode")
	fmt.Println("")
	if test2FACode == "exit" || test2FACode == "Exit" || test2FACode == "EXIT" {
		exitProgramCLI()
	} else if validationTest2FACode == false {
		invalidInputCLI()
		add2FACLI()
	} else {
		correct2FACode := totp.Validate(test2FACode, generated2FAKey)
		if correct2FACode {
			clearScreen()
			replaceText("2fa_not_set", generated2FAKey)
			messageBoxCLI(bgGreen, textBoldWhite, "Successfully added 2FA secret key")
			fmt.Println("")
			os.Exit(0)
		} else {
			clearScreen()
			messageBoxCLI(bgRed, textBoldWhite, "2FA code entered is incorrect, press enter/return to continue")
			fmt.Scanln()
			add2FACLI()
		}
	}
}

// Basic function with no parameters to call the addEmailCli, addPasswordCli and add2faCli functions
func createUserCLI() {
	addEmailCLI()
	addPasswordCLI()
	add2FACLI()
}

// Basic function with no parameters to call the addPasswordCli and add2faCli functions
func createPassword2FACLI() {
	addPasswordCLI()
	add2FACLI()
}

// Recursive function to change an existing known password
func changePasswordCLI(envPassword string) {
	clearScreen()
	typeExitCLI()
	fmt.Println(bgRed + textBoldWhite)
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(" □                                                                   □ ")
	fmt.Println(" □      The change_password option is set to yes in mfaview.env      □ ")
	fmt.Println(" □    To change the password the exisitng password must be known     □ ")
	fmt.Println(" □  New password must be 16-32 characters & not contain a comma (,)  □ ")
	fmt.Println(" □                                                                   □ ")
	fmt.Println(" □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ □ ")
	fmt.Println(resetColour)
	var currentPassword string
	fmt.Printf(textBoldBlack)
	fmt.Printf("      Please enter current password: ")
	fmt.Scan(&currentPassword)
	validationCurrentPassword := validateInput(currentPassword, "password")
	correctPassword := comparePassword([]byte(zeroPad(currentPassword)), []byte(envPassword))
	if currentPassword == "exit" || currentPassword == "Exit" || currentPassword == "EXIT" {
		exitProgramCLI()
	} else if validationCurrentPassword == false {
		clearScreen()
		messageBoxCLI(bgRed, textBoldWhite, "Current password needs to be 16-32 characters, press enter/return to continue")
		fmt.Scanln()
		changePasswordCLI(envPassword)
	} else if correctPassword == false {
		clearScreen()
		messageBoxCLI(bgRed, textBoldWhite, "Password incorrect ")
		fmt.Scanln()
		changePasswordCLI(envPassword)
	}
	fmt.Println("")
	var newPassword string
	fmt.Printf("      Please enter new password: ")
	fmt.Scan(&newPassword)
	validationNewPassword := validateInput(newPassword, "password")
	if newPassword == "exit" || newPassword == "Exit" || newPassword == "EXIT" {
		exitProgramCLI()
	} else if validationNewPassword == false {
		clearScreen()
		messageBoxCLI(bgRed, textBoldWhite, "New password needs to be 16-32 characters, press enter/return to continue")
		fmt.Scanln()
		changePasswordCLI(envPassword)
	}
	fmt.Println("")
	var newPasswordCheck string
	fmt.Printf("      Please re-enter the new password: ")
	fmt.Scan(&newPasswordCheck)
	validationNewPasswordCheck := validateInput(newPasswordCheck, "password")
	if newPasswordCheck == "exit" || newPasswordCheck == "Exit" || newPasswordCheck == "EXIT" {
		exitProgramCLI()
	} else if validationNewPasswordCheck == false {
		clearScreen()
		messageBoxCLI(bgRed, textBoldWhite, "New re-entered password needs to be 16-32 characters, press enter/return to continue")
		fmt.Scanln()
		changePasswordCLI(envPassword)
	}
	if newPassword != newPasswordCheck {
		clearScreen()
		messageBoxCLI(bgRed, textBoldWhite, "New passwords entered do not match, press enter/return to continue ")
		fmt.Scanln()
		changePasswordCLI(envPassword)
	}

	wipeFile := true

	readKeyCSV := csvcell.ReadCSV(dirKeyCSV, fileKeyCSV)
	for _, readKeyCSV := range readKeyCSV {
		accountName := strings.Join((readKeyCSV[0:][0:1]), ", ")
		oldDecryptedKey := aestext.DecText(strings.Join((readKeyCSV[0:][1:2]), ", "), zeroPad(currentPassword))
		sha := strings.Join((readKeyCSV[0:][2:3]), ", ")
		dateAdded := strings.Join((readKeyCSV[0:][3:4]), ", ")
		if wipeFile == true {
			wipeKeyFile()
			wipeFile = false
		}
		data := accountName + "," + aestext.EncText(oldDecryptedKey, zeroPad(newPassword)) + "," + sha + "," + dateAdded
		csvcell.WriteCSV(dirKeyCSV, fileKeyCSV, 0, data, 0)
	}

	newHashedPassword := genPassword([]byte(zeroPad(newPassword)))
	replaceText(envPassword, newHashedPassword)
	clearScreen()
	messageBoxCLI(bgGreen, textBoldWhite, "Successfully updated password")
	fmt.Println("")
	os.Exit(0)
}

// Function to generate passcode from MFA secret key
func readMFA(key string, sha otp.Algorithm) (passcode string) {
	passcode, err := totp.GenerateCodeCustom(key, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: sha,
	})

	if err != nil {
		panic(err)
	}

	return passcode
}

func javaScriptCopy(w http.ResponseWriter, sha string) {
	fmt.Fprintf(w, "<script>")
	fmt.Fprintf(w, "  function cp"+sha+"() {")
	fmt.Fprintf(w, "    navigator.clipboard.writeText('"+sha+"');")
	fmt.Fprintf(w, "  }")
	fmt.Fprintf(w, "</script>")
}

func copyHTMLButton(w http.ResponseWriter, sha string) {
	fmt.Fprintf(w, "&nbsp <button onclick=cp"+sha+"() class=\"copyButton\">&#10697</button>")
}

// Declare two string variables for HTML
var startHTML string
var endHTML string

// Main function
func main() {
	startHTML = getFile(dirHTML, fileStartHTML)
	endHTML = getFile(dirHTML, fileEndHTML)

	err := godotenv.Load(mfaViewEnv)
	if err != nil {
		panic("Error loading mfaview.env file")
	}

	envEmail := os.Getenv("email")
	envPassword := os.Getenv("password")
	env2FAKey := os.Getenv("2fa_key")
	envChangePassword := os.Getenv("change_password")
	envAddress := os.Getenv("address")
	envPort := os.Getenv("port")
	envAddAccount := os.Getenv("add_account_page")

	validationEnvEmail := validateInput(envEmail, "email")
	validationEnvPassword := validateInput(envPassword, "hashedPassword")
	validationEnv2FA := validateInput(env2FAKey, "secretKey")
	validationEnvAddress := validateInput(envAddress, "IP")

	envPortInt, err := strconv.Atoi(envPort)
	if err != nil {
		invalidEnvCLI("Port must be a number in " + mfaViewEnv)
	}

	if envPortInt <= 0 || envPortInt >= 65536 {
		invalidEnvCLI("Port number in " + mfaViewEnv + " must be between 1 and 65535")
	} else if envAddress != "localhost" {
		if validationEnvAddress == false {
			invalidEnvCLI("Address in " + mfaViewEnv + " must be a valid Internet Protocol (IP) address or localhost")
		}
	} else if envEmail == "email_not_set" && envPassword == "password_not_set" && env2FAKey == "2fa_not_set" {
		createUserCLI()
	} else if validationEnvEmail == false {
		invalidEnvCLI("Email address stored in " + mfaViewEnv + " is invalid")
	} else if validationEnvEmail == true && envPassword == "password_not_set" && env2FAKey == "2fa_not_set" {
		createPassword2FACLI()
	} else if env2FAKey == "2fa_not_set" {
		add2FACLI()
	} else if validationEnvPassword == false {
		invalidEnvCLI("Password stored in " + mfaViewEnv + " is invalid ")
	} else if validationEnv2FA == false {
		invalidEnvCLI("2FA secret key stored stored in " + mfaViewEnv + " is invalid")
	} else if envChangePassword == "yes" || envChangePassword == "Yes" || envChangePassword == "YES" {
		changePasswordCLI(envPassword)
	} else {

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
			}

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
			fmt.Fprintf(w, "<br>")
			fmt.Fprintf(w, "<br>")
			fmt.Fprintf(w, "<br>")

			// Get email address, password and 2FA from HTTP POST
			inputEmail := r.FormValue("email")
			inputPassword := r.FormValue("password")
			input2FA := r.FormValue("2FA")

			// Validate email address, password and 2FA
			validationEmail := validateInput(inputEmail, "email")
			validationPassword := validateInput(inputPassword, "password")
			validation2FA := validateInput(input2FA, "2FACode")

			// Conditional statement to validate input
			if inputEmail == "" && inputPassword == "" && input2FA == "" {
			} else if validationEmail == false {
				textBox(w, "Please enter a valid email address, max 320 characters in length")
			} else if validationPassword == false {
				textBox(w, "Password needs to be between 16-32 characters in length")
			} else if validation2FA == false {
				textBox(w, "MFA code needs to be a 6 digit number")
			} else if inputEmail == envEmail {
				correctPassword := comparePassword([]byte(zeroPad(inputPassword)), []byte(envPassword))
				correct2FA := totp.Validate(input2FA, env2FAKey)
				if correctPassword == true && correct2FA == true {
					fmt.Fprintf(w, "<table class=mfaTableTextColor>")
					fmt.Fprintf(w, "  <tr>")
					fmt.Fprintf(w, "    <th class=accountNameTitleColor>Account Name</th>")
					fmt.Fprintf(w, "    <th class=sha1TitleColor>SHA1 Code<br>(Default)</th>")
					fmt.Fprintf(w, "    <th class=sha256TitleColor>SHA256 Code</th>")
					fmt.Fprintf(w, "    <th class=sha512TitleColor>SHA512 Code</th>")
					fmt.Fprintf(w, "    <th class=dateAddedTitleColor>Date Added<br>(DD-MM-YYYY)</th>")
					fmt.Fprintf(w, "  </tr>")

					// Read key.csv file
					readKeyCSV := csvcell.ReadCSV(dirKeyCSV, fileKeyCSV)
					for _, readKeyCSV := range readKeyCSV {
						accountName := strings.Join((readKeyCSV[0:][0:1]), ", ")
						decryptedKey := aestext.DecText(strings.Join((readKeyCSV[0:][1:2]), ", "), zeroPad(inputPassword))
						sha := strings.Join((readKeyCSV[0:][2:3]), ", ")
						dateAdded := strings.Join((readKeyCSV[0:][3:4]), ", ")

						fmt.Fprintf(w, "  <tr>")
						fmt.Fprintf(w, "    <td class=accountNameValueColor><b>"+accountName+"</b></td>")
						if sha == "SHA1" {
							sha1 := readMFA(decryptedKey, otp.AlgorithmSHA1)
							javaScriptCopy(w, sha1)
							fmt.Fprintf(w, "    <td class=sha1CodeColor><b>"+sha1+"</b>")
							copyHTMLButton(w, sha1)
							fmt.Fprintf(w, "    </td>")
						} else {
							fmt.Fprintf(w, "    <td class=sha1CodeColor><b>&#9473&#9473</b></td>")
						}
						if sha == "SHA256" {
							sha256 := readMFA(decryptedKey, otp.AlgorithmSHA256)
							javaScriptCopy(w, sha256)
							fmt.Fprintf(w, "    <td class=sha256CodeColor><b>"+sha256+"</b>")
							copyHTMLButton(w, sha256)
							fmt.Fprintf(w, "    </td>")
						} else {
							fmt.Fprintf(w, "    <td class=sha256CodeColor><b>&#9473&#9473</b></td>")
						}
						if sha == "SHA512" {
							sha512 := readMFA(decryptedKey, otp.AlgorithmSHA512)
							javaScriptCopy(w, sha512)
							fmt.Fprintf(w, "    <td class=sha512CodeColor><b>"+sha512+"</b>")
							copyHTMLButton(w, sha512)
							fmt.Fprintf(w, "    </td>")
						} else {
							fmt.Fprintf(w, "    <td class=sha512CodeColor><b>&#9473&#9473</b></td>")
						}
						fmt.Fprintf(w, "    <td class=dateAddedValueColor><b>"+dateAdded+"</b></td>")
						fmt.Fprintf(w, "  </tr>")
					}
					fmt.Fprintf(w, "</table>")
					fmt.Fprintf(w, "<br>")
					fmt.Fprintf(w, "<br>")
				} else {
					textBox(w, "Wrong Credentials Entered")
				}
			} else {
				textBox(w, "Wrong Credentials Entered")
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

				fmt.Fprintf(w, startHTML)
				messageTable(w, "/", "Click to login & view MFA account(s)", "Enter email, password, MFA account<br>name, MFA secret key, select SHA &<br>enter 2FA code to add a new account.<br>Commas (,) not allowed.")
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<div>")
				fmt.Fprintf(w, "<form method=\"POST\" action=\"/add-account\">")
				inputForm(w, "email", "Email Address", "email")
				inputForm(w, "password", "Password", "password")
				inputForm(w, "account", "New MFA Account Name", "text")
				inputForm(w, "MFA", "New MFA Secret Key", "text")
				fmt.Fprintf(w, "  <label for=\"SHA\"><b>Secure Hash Algorithm (SHA):<br>(SHA1 is Default)</b>")
				fmt.Fprintf(w, "  </label><br>")
				fmt.Fprintf(w, "  <select id=\"SHA\" name=\"SHA\">")
				var shaList = []string{"SHA1", "SHA256", "SHA512"}
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
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<br>")

				// Get email address, password, MFA account name, MFA secret key, SHA type and 2FA from HTTP POST
				inputEmail := r.FormValue("email")
				inputPassword := r.FormValue("password")
				inputAccount := r.FormValue("account")
				inputMFA := r.FormValue("MFA")
				inputSHA := r.FormValue("SHA")
				input2FA := r.FormValue("2FA")

				// Returns false or true based on if input is valid
				validationEmail := validateInput(inputEmail, "email")
				validationPassword := validateInput(inputPassword, "password")
				validationAccount := validateInput(inputAccount, "accountName")
				validationMFA := validateInput(inputMFA, "secretKey")
				validationSHA := slices.Contains(shaList, inputSHA)
				validation2FA := validateInput(input2FA, "2FACode")

				// Conditional statement to validate input
				if inputEmail == "" && inputPassword == "" && inputAccount == "" && inputMFA == "" && inputSHA == "" && input2FA == "" {
				} else if validationEmail == false {
					textBox(w, "Please enter a valid email address, max 320 characters in length")
				} else if validationPassword == false {
					textBox(w, "Password needs to be between 16-32 characters in length")
				} else if validationAccount == false {
					textBox(w, "Please enter a valid account name, max 100 characters in length")
				} else if validationMFA == false {
					textBox(w, "MFA secret key needs to be between 16-128 characters")
				} else if validationSHA == false {
					textBox(w, "Secure Hash Algorithm (SHA) can be SHA1, SHA256 or SHA512")
				} else if validation2FA == false {
					textBox(w, "2FA code needs to be a 6 digit number")
				} else if inputEmail == envEmail {
					correctPassword := comparePassword([]byte(zeroPad(inputPassword)), []byte(envPassword))
					correct2FA := totp.Validate(input2FA, env2FAKey)
					if correctPassword == true && correct2FA == true {
						date := time.Now().Local()
						data := inputAccount + "," + aestext.EncText(inputMFA, zeroPad(inputPassword)) + "," + inputSHA + "," + date.Format("02-01-2006")
						csvcell.WriteCSV(dirKeyCSV, fileKeyCSV, 0, data, 0)
						textBox(w, "Account Added")
					} else {
						textBox(w, "Wrong Credentials Entered")
					}
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

				fmt.Fprintf(w, startHTML)
				messageTable(w, "/", "Click to login & view MFA account(s)", "Adding a new MFA account is<br>switched off in the configuration file")
				fmt.Fprintf(w, endHTML)
			})
		}
	}

	// IP address and port number, value taken from mfaview.env
	socket := envAddress + ":" + envPort
	messageBoxCLI(bgGreen, textBoldWhite, "MFA View is running on "+socket)
	fmt.Println("")

	// Start server on port specified above
	log.Fatal(http.ListenAndServe(socket, nil))
}

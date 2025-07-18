# MFA View

### A Go HTTP program to generate MFA codes from a stored secret key and can also add a new MFA key.
### MFA keys are encrypted and the user's password is hashed and salted.

#### (Built using operating system - Ubuntu 22.04.5 LTS)

<br>
<br>

### Access demo site at https://mfa.ell.today
#### Demo username is: `demo@ell.today`
#### Demo password is: `passwordpasswordpassword1$`
#### Demo 2FA secret key is: `IHJU2OCCJCGAH5CMTSIBPU66SVLDVIRV`
#### or scan the QR code below for secret key

<br>
<br>

#### MFA View can be installed easily with the BASH script - [install-mfaview.sh](https://github.com/ellwould/mfa-view/blob/main/install-mfaview.sh)

#### MFA View can be un-installed easily with the BASH script - [uninstall-mfaview.sh](https://github.com/ellwould/mfa-view/blob/main/uninstall-mfaview.sh)

---

## Website:

### Login page:

![image](https://github.com/ellwould/mfa-view/blob/main/image/WEB_login_page.png)

### Logging in with no MFA accounts present:

![image](https://github.com/ellwould/mfa-view/blob/main/image/WEB_no_MFA_accounts_added.png)

### Adding a MFA acount:
#### (In this example I generated a random secret key for the purpose of demonstrating - `V2J6E2BTG3LPGRWB63CVCKXXRGENYS3K`)

![image](https://github.com/ellwould/mfa-view/blob/main/image/WEB_adding_a_MFA_account.png)

### Logging in with a MFA account present:
#### (Java Script used to create a copy to clipboard button)

![image](https://github.com/ellwould/mfa-view/blob/main/image/WEB_new_MFA_code_avaiable.png)

### Screenshot of he mfaview-key.csv file (located in /etc/mfaview/key) showing the secret key is encrypted:
#### (AES encryption used)

![image](https://github.com/ellwould/mfa-view/blob/main/image/CLI_encrypted_secret_key.png)

### Add MFA account page switched off in the configuration file:

![image](https://github.com/ellwould/mfa-view/blob/main/image/WEB_account_page_switched_off.png)

### Examples of validation on the login page:

![image](https://github.com/ellwould/mfa-view/blob/main/image/WEB_wrong_details_entered.png)

![image](https://github.com/ellwould/mfa-view/blob/main/image/WEB_password_entered_wrong_length.png)

![image](https://github.com/ellwould/mfa-view/blob/main/image/WEB_no_2FA_entered.png)

---

## CLI:

### Adding an email address:
#### (`demo@ell.today` was used)

![image](https://github.com/ellwould/mfa-view/blob/main/image/CLI_creating_an_email.png)

### Adding a password:
#### (Easy example password used was `passwordpasswordpassword1$`):

![image](https://github.com/ellwould/mfa-view/blob/main/image/CLI_creating_a_password.png)

### Generating a 2FA secret key:
#### (Secret key generated - `IHJU2OCCJCGAH5CMTSIBPU66SVLDVIRV`)

![image](https://github.com/ellwould/mfa-view/blob/main/image/CLI_2FA_Secret_key.png)

### 2FA secret key embedded inside a Quick Response (QR) code for a 3rd party authenticator app to scan:

![image](https://github.com/ellwould/mfa-view/blob/main/image/CLI_Quick_Response_QR_code.png)

### 3rd party application generating 2FA code:

![image](https://github.com/ellwould/mfa-view/blob/main/image/3rd_party_authenticator_app.jpeg)

### MFA View checks the 2FA code generated from a 3rd party authenticator app is correct before adding the secret key to the configuration file:

![image](https://github.com/ellwould/mfa-view/blob/main/image/CLI_2FA_check_code.png)

### MFA View informing the user the 2FA secret key has been added to the configuration file:

![image](https://github.com/ellwould/mfa-view/blob/main/image/CLI_message_2FA_correct.png)

### Example of configuration file (password is hashed and salted):

![image](https://github.com/ellwould/mfa-view/blob/main/image/CLI_configuration_file_after_account_setup.png)

### Systemd after MFA View was installed but before account setup:

![image](https://github.com/ellwould/mfa-view/blob/main/image/CLI_systemctl_before_adding_an_account.png)

### Systemd after account setup is complete:

![image](https://github.com/ellwould/mfa-view/blob/main/image/CLI_systemctl_after_adding_an_account.png)

### Examples of validation when creating an account:

![image](https://github.com/ellwould/mfa-view/blob/main/image/CLI_invalid_input.png)

![image](https://github.com/ellwould/mfa-view/blob/main/image/CLI_password_does_not_match.png)

---

>[!NOTE]
>For a list of abbreviations and there meanings used throughout this repository please refer to this [README](https://github.com/Ellwould/information_technology_and_telecommunication_abbreviations)

<br>

> [!IMPORTANT]
> All third-party product and/or company names and logos are trademarks™ or registered® trademarks and remain the property of their respective holders/owners. Unless specifically identified as such, use of third party trademarks does not imply any affiliation with or endorsement between Elliot Michael Keavney and the owners of those trademarks.

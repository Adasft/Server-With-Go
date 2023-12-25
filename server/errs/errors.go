package errs

import (
	"log"
	"net/http"

	"server/template"
)

const (
	/*
	 * error messages that will be displayed in the client view
	 */
	IncorrectPasswordError = "Wrong password. Please try again or click on the 'Forgot your password?' link. Remember that you only have %d attempts."
	NoUserFoundError       = "No user has been found registered with the email '%s'"
	InvalidEmailError      = "The email '%s' is not a valid email."
	EmptyPasswordError     = "The password cannot be empty."
	InvalidUsernameError   = "The value '%s' is not valid data for the 'username' field."
	ShortPasswordError     = "Password must contain %d letters or more (current letters: %d)."
	DuplicateEmailError    = "The email '%s' is already registered. Please enter a different email address."
	PasswordsNotMatchError = "Passwords do not match"
	AccountBlockedError    = "Your account has been locked due to an excessive number of failed password attempts. To unlock your account, please click <a href='%s'>here</a> to recover your account. We are sorry for any inconvenience this may cause and we are here to help."

	/*
	 * server status error messages
	 */
	InternalServerError = "500 Internal Server Error"

	/*
	 * program error messages
	 */
	RouterIDNotExistError                 = "router ID does not exist"
	DatabaseConnectionNotEstablishedError = "database connection has not been established"
	DatabaseConnectionNotOpenError        = "database connection is not open"
	DatabaseConnectionAlreadyOpenError    = "database connection is already open"
)

func InternalServerErrorHandler(w http.ResponseWriter, err error, BackRoute string) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(err)

	_, err = template.Render(w, &template.InternalServerErrorPageData{
		Title:     "500 Internal Server Error",
		BackRoute: BackRoute,
	}, template.GetView("500"))

	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
	}
}

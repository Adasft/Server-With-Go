package routes

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-session/session"
	"golang.org/x/crypto/bcrypt"
	"server/db"
	"server/errs"
	"server/form"
	"server/routerutils"
	"server/template"
	"server/utils"
)

var templateSignupData = &template.SignupPageData{}

func signupHandlerGet(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)

	if err != nil {
		errs.InternalServerErrorHandler(w, err, SignupPath)
		return
	}

	_, ok := store.Get("user_id")

	if ok {
		http.Redirect(w, r, HomePath, http.StatusSeeOther)
		return
	}

	_, err = template.Render(w, templateSignupData, template.GetView("index"), template.GetLayout("signup"))

	if err != nil {
		errs.InternalServerErrorHandler(w, err, SignupPath)
	}

	templateSignupData.EnableErrorView(false)
	templateSignupData.ClearErrors()
}

func signupHandlerPost(w http.ResponseWriter, r *http.Request) {
	connection, err := db.HandlerConnector.GetConnection()

	if err != nil {
		log.Fatal(err)
	}

	signupFormFields, err := validateSignupFormFields(r, connection)

	if err != nil {
		errs.InternalServerErrorHandler(w, err, SignupPath)
		return
	}

	if templateSignupData.HasErrors() {
		templateSignupData.EnableErrorView(true)
		http.Redirect(w, r, SignupPath, http.StatusSeeOther)
		return
	}

	if err = insertNewUser(connection, signupFormFields); err != nil {
		errs.InternalServerErrorHandler(w, err, SignupPath)
		return
	}

	http.Redirect(w, r, LoginPath, http.StatusSeeOther)
}

func insertNewUser(connection *sql.DB, signupFormFields *form.SignupFormFields) error {
	query := "INSERT INTO users(username, email, password) VALUES (?, ?, ?)"

	hashedPassword, err := generateHashedPassword([]byte(signupFormFields.Password))

	if err != nil {
		return err
	}

	_, err = connection.Exec(query, signupFormFields.Username, signupFormFields.Email, hashedPassword)

	if err != nil {
		return err
	}

	return nil
}

func generateHashedPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

func getSignupFormFields(r *http.Request) (*form.SignupFormFields, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	return &form.SignupFormFields{
		Username:        r.Form.Get(form.UsernameFieldName),
		Email:           r.Form.Get(form.EmailFieldName),
		Password:        r.Form.Get(form.PasswordFieldName),
		ConfirmPassword: r.Form.Get(form.ConfirmPasswordFieldName),
	}, nil
}

func validateSignupFormFields(r *http.Request, connection *sql.DB) (*form.SignupFormFields, error) {
	signupFormFields, err := getSignupFormFields(r)

	if err != nil {
		return nil, err
	}

	if utils.IsEmptyStr(signupFormFields.Username) {
		templateSignupData.PushError(fmt.Sprintf(errs.InvalidUsernameError, signupFormFields.Username))
	}

	if len(signupFormFields.Password) < form.MinNumberCharsPassword {
		templateSignupData.PushError(fmt.Sprintf(
			errs.ShortPasswordError,
			form.MinNumberCharsPassword,
			len(signupFormFields.Password),
		))
	}

	if signupFormFields.Password != signupFormFields.ConfirmPassword {
		templateSignupData.PushError(errs.PasswordsNotMatchError)
	}

	if !utils.IsValidEmail(signupFormFields.Email) {
		templateSignupData.PushError(fmt.Sprintf(errs.InvalidEmailError, signupFormFields.Email))
	}

	if err := checkDuplicateEmail(connection, signupFormFields.Email); err != nil {
		return nil, err
	}

	return signupFormFields, nil
}

func checkDuplicateEmail(connection *sql.DB, email string) error {
	var countOfEmail int
	query := "SELECT COUNT(*) FROM users WHERE email=?"

	if err := connection.QueryRow(query, email).Scan(&countOfEmail); err != nil {
		return err
	}

	if countOfEmail > 0 {
		templateSignupData.PushError(fmt.Sprintf(errs.DuplicateEmailError, email))
	}

	return nil
}

func initSignupRouter(router *routerutils.Router) {
	templateSignupData.FillDefault()
	router.Get(SignupPath, signupHandlerGet, denyAccessIfAlreadyLoggedInMiddleware)
	router.Post(SignupPath, signupHandlerPost, nil)
}

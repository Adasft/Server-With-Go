package routes

import (
	"context"
	"database/sql"
	"errors"
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

const maxLoginAttempts = 2

var templateLoginData = &template.LoginPageData{}

func loginHandlerGet(w http.ResponseWriter, r *http.Request) {
	_, err := template.Render(w, templateLoginData, template.GetView("index"), template.GetLayout("login"))

	if err != nil {
		errs.InternalServerErrorHandler(w, err, LoginPath)
	}

	templateLoginData.EnableErrorView(false)
	templateLoginData.ClearErrors()
}

func loginHandlerPost(w http.ResponseWriter, r *http.Request) {
	connection, err := db.HandlerConnector.GetConnection()

	if err != nil {
		log.Fatal(err)
	}

	loginFormFields, ok, err := validateLoginFormFields(r)

	if err != nil {
		errs.InternalServerErrorHandler(w, err, LoginPath)
		return
	}

	if !ok {
		templateLoginData.EnableErrorView(true)
		http.Redirect(w, r, LoginPath, http.StatusSeeOther)
		return
	}

	user, err := getUserByEmail(connection, loginFormFields.Email)

	if err != nil {
		handleUserLookupError(w, r, loginFormFields.Email, err)
		return
	}

	if user.IsLocked {
		dealWithBlockedAccount()
		http.Redirect(w, r, LoginPath, http.StatusSeeOther)
		return
	}

	if err = compareHashPassword(user.Password, loginFormFields.Password); err != nil {
		handlePasswordComparisonError(w, r, connection, err, user.UserId)
		return
	}

	if err = startSession(w, r, user.UserId); err != nil {
		errs.InternalServerErrorHandler(w, err, LoginPath)
		return
	}

	if err = resetLoginAttempts(connection, user.UserId); err != nil {
		errs.InternalServerErrorHandler(w, err, LoginPath)
		return
	}

	http.Redirect(w, r, HomePath, http.StatusSeeOther)
}

func resetLoginAttempts(connection *sql.DB, userId int) error {
	_, err := connection.Exec("UPDATE users SET login_attempts = 0 WHERE user_id = ?", userId)
	return err
}

func incrementLoginAttempts(connection *sql.DB, userId int) error {
	_, err := connection.Exec("UPDATE users SET login_attempts = login_attempts + 1 WHERE user_id = ?", userId)
	return err
}

func loginAttemptHandler(connection *sql.DB, userId int) (bool, error) {
	if err := incrementLoginAttempts(connection, userId); err != nil {
		return false, err
	}

	return checkAndLockAccount(connection, userId)
}

func dealWithBlockedAccount() {
	templateLoginData.EnableErrorView(true)
	templateLoginData.PushError(fmt.Sprintf(errs.AccountBlockedError, RecoverPath))
}

func checkAndLockAccount(connection *sql.DB, userId int) (bool, error) {
	var (
		loginAttempts int
		isLocked      bool
		err           error
	)

	if err = connection.QueryRow("SELECT login_attempts FROM users WHERE user_id=?", userId).Scan(&loginAttempts); err != nil {
		return false, err
	}

	if loginAttempts > maxLoginAttempts {
		if _, err = connection.Exec("UPDATE users SET is_locked = 1 WHERE user_id=?", userId); err != nil {
			return false, err
		}

		isLocked, err = isAccountBlocked(connection, userId)

		if err != nil {
			return false, err
		}

		if isLocked {
			dealWithBlockedAccount()
		}

	}

	return isLocked, nil
}

func isAccountBlocked(connection *sql.DB, userId int) (bool, error) {
	var isLocked bool

	if err := connection.QueryRow("SELECT is_locked FROM users WHERE user_id=?", userId).Scan(&isLocked); err != nil {
		return false, err
	}

	return isLocked, nil
}

func startSession(w http.ResponseWriter, r *http.Request, userId int) error {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		return err
	}

	store.Set("user_id", userId)

	if err = store.Save(); err != nil {
		return err
	}

	return nil
}

func compareHashPassword(hashedPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}

func handlePasswordComparisonError(w http.ResponseWriter, r *http.Request, connection *sql.DB, err error, userId int) {
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		isLocked, err := loginAttemptHandler(connection, userId)

		if err != nil {
			errs.InternalServerErrorHandler(w, err, LoginPath)
			return
		}

		if !isLocked {
			templateLoginData.EnableErrorView(true)
			templateLoginData.PushError(fmt.Sprintf(errs.IncorrectPasswordError, maxLoginAttempts))
		}

		http.Redirect(w, r, LoginPath, http.StatusSeeOther)
		return
	}

	errs.InternalServerErrorHandler(w, err, LoginPath)
}

func getUserByEmail(connection *sql.DB, email string) (*form.User, error) {
	query := "SELECT * FROM users WHERE email=?"

	var user form.User

	if err := connection.QueryRow(query, email).Scan(
		&user.UserId, &user.Username, &user.Email, &user.Password, &user.LoginAttempts, &user.IsLocked, &user.PhoneId,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func handleUserLookupError(w http.ResponseWriter, r *http.Request, email string, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		templateLoginData.EnableErrorView(true)
		templateLoginData.PushError(fmt.Sprintf(errs.NoUserFoundError, email))
		http.Redirect(w, r, LoginPath, http.StatusSeeOther)
	} else {
		errs.InternalServerErrorHandler(w, err, LoginPath)
	}
}

func getLoginFormFields(r *http.Request) (*form.LoginFormFields, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	return &form.LoginFormFields{
		Email:    r.Form.Get(form.EmailFieldName),
		Password: r.Form.Get(form.PasswordFieldName),
	}, nil
}

func validateLoginFormFields(r *http.Request) (*form.LoginFormFields, bool, error) {
	loginFormFields, err := getLoginFormFields(r)

	if err != nil {
		return nil, false, nil
	}

	if !utils.IsValidEmail(loginFormFields.Email) {
		templateLoginData.PushError(fmt.Sprintf(errs.InvalidEmailError, loginFormFields.Email))
	}

	if utils.IsEmptyStr(loginFormFields.Password) {
		templateLoginData.PushError(errs.EmptyPasswordError)
	}

	if templateLoginData.HasErrors() {
		return nil, false, nil
	}

	return loginFormFields, true, nil
}

func initLoginRouter(router *routerutils.Router) {
	templateLoginData.FillDefault()
	router.Get(LoginPath, loginHandlerGet, denyAccessIfAlreadyLoggedInMiddleware)
	router.Post(LoginPath, loginHandlerPost, nil)
}

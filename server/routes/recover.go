package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"server/db"
	"server/errs"
	"server/form"
	"server/routerutils"
	"server/template"
	"server/utils"
)

var templateRecoveryData template.RecoveryPageData

func recoverHandlerGet(w http.ResponseWriter, r *http.Request) {
	_, err := template.Render(w, nil, template.GetView("recover"))

	if err != nil {
		errs.InternalServerErrorHandler(w, err, RecoverPath)
	}
}

func recoverHandlerPost(w http.ResponseWriter, r *http.Request) {
	connection, err := db.HandlerConnector.GetConnection()

	if err != nil {
		log.Fatal(err)
	}

	recoveryFormFields, err := getRecoveryFormFields(r)

	if err != nil {
		errs.InternalServerErrorHandler(w, err, LoginPath)
		return
	}

	user, err := getUserByRecoveryMethod(connection, recoveryFormFields.Value)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			templateRecoveryData.EnableErrorView(true)
			templateRecoveryData.PushError("Error")
		}
	}

	fmt.Println(user)

}

func getUserByRecoveryMethod(connection *sql.DB, recoveryMethodValue string) (*form.UserWithPhone, error) {
	query := "SELECT * FROM users.*, phone LEFT JOIN phones ON users.phone_id = phones.phone_id WHERE users.email = ? OR phones.phone = ?"

	var user form.UserWithPhone

	if err := connection.QueryRow(query, recoveryMethodValue, recoveryMethodValue).Scan(
		&user.UserId, &user.Username, &user.Email, &user.Password, &user.LoginAttempts, &user.IsLocked, &user.PhoneId, &user.Phone,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func predictTypeOfRecoveryMethod(recoveryMethodValue string) form.RecoveryMethodType {
	var recoveryMethodType form.RecoveryMethodType

	if utils.IsValidEmail(recoveryMethodValue) {
		recoveryMethodType.SetAsEmail()
	} else if utils.IsValidPhoneNumber(recoveryMethodValue) {
		recoveryMethodType.SetAsPhone()
	}

	return recoveryMethodType
}

func getRecoveryFormFields(r *http.Request) (*form.RecoveryFromFields, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	recoveryMethodValue := r.Form.Get(form.RecoveryMethodFieldName)
	recoveryMethodType := predictTypeOfRecoveryMethod(recoveryMethodValue)

	return &form.RecoveryFromFields{
		Value:      recoveryMethodValue,
		MethodType: recoveryMethodType,
	}, nil
}

func initRecoverRouter(router *routerutils.Router) {
	templateRecoveryData.FillDefault()
	router.Get(RecoverPath, recoverHandlerGet, denyAccessIfAlreadyLoggedInMiddleware)
	router.Post(RecoverPath, recoverHandlerPost, nil)
}

package form

type User struct {
	UserId        int
	Username      string
	Email         string
	Password      string
	LoginAttempts int
	IsLocked      bool
}

type LoginFormFields struct {
	Email    string
	Password string
}

type SignupFormFields struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
}

const (
	MinNumberCharsPassword   = 5
	UsernameFieldName        = "username"
	EmailFieldName           = "email"
	PasswordFieldName        = "password"
	ConfirmPasswordFieldName = "confirm_password"
)

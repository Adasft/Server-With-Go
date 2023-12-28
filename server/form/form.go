package form

type User struct {
	UserId        int
	Username      string
	Email         string
	Password      string
	LoginAttempts int
	IsLocked      bool
	PhoneId       int
}

type UserWithPhone struct {
	User
	Phone       string
	CountryCode string
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

type RecoveryMethodType string

func (r *RecoveryMethodType) SetAsEmail() {
	*r = "email"
}

func (r *RecoveryMethodType) SetAsPhone() {
	*r = "phone"
}

type RecoveryFromFields struct {
	Value      string
	MethodType RecoveryMethodType
}

const (
	MinNumberCharsPassword   = 5
	UsernameFieldName        = "username"
	EmailFieldName           = "email"
	PasswordFieldName        = "password"
	ConfirmPasswordFieldName = "confirm_password"
	RecoveryMethodFieldName  = "recovery_method"
)

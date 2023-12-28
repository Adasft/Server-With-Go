package template

import "server/utils"

type PageFormErrors struct {
	ShowErrors bool
	Errors     []string
}

func (f *PageFormErrors) EnableErrorView(state bool) {
	f.ShowErrors = state
}

func (f *PageFormErrors) PushError(errMessage string) {
	utils.Append[string](&f.Errors, errMessage)
}

func (f *PageFormErrors) ClearErrors() {
	f.Errors = []string{}
}

func (f *PageFormErrors) HasErrors() bool {
	return len(f.Errors) > 0
}

type FormPageData struct {
	Title string
	PageFormErrors
}

func (f *FormPageData) SetTitle(title string) {
	f.Title = title
}

func (f *FormPageData) FillDefault() {
	f.ShowErrors = false
	f.Errors = []string{}
}

type LoginPageData struct {
	FormPageData
}

func (l *LoginPageData) FillDefault() {
	l.Title = "Log In"
}

type SignupPageData struct {
	FormPageData
}

func (s *SignupPageData) FillDefault() {
	s.Title = "Sign Up"
}

type RecoveryPageData struct {
	FormPageData
}

func (r *RecoveryPageData) FillDefault() {
	r.Title = "Recovery"
}

type InternalServerErrorPageData struct {
	Title     string
	BackRoute string
}

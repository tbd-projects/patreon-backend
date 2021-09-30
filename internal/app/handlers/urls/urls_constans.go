package urls

type UrlPath string

const (
	Register = "/register"
	Login = "/login"
	Creator = "/creators"
	Logout = "/logout"
	Profile = "/profile"
	ID = "/{id}"
)

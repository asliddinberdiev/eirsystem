package codes

import "net/http"

type Code int

const (
	// SYSTEM -> 0 - 999
	Ok       Code = 200
	InternalError Code = 500
	InvalidParams Code = 400

	// USER -> 1000 - 1999
	UserNotFound      Code = 1001
	UserAlreadyExists Code = 1002
	UserPasswordWrong Code = 1003
	UserInactive      Code = 1004

	// AUTH -> 2000 - 2999
	AuthTokenExpired Code = 2001
	AuthTokenInvalid Code = 2002
	AuthRequired     Code = 2003
)

func (c Code) HTTPStatus() int {
	switch c {
	case Ok:
		return http.StatusOK
	case InternalError:
		return http.StatusInternalServerError
	case InvalidParams, UserAlreadyExists, UserPasswordWrong:
		return http.StatusBadRequest
	case UserNotFound:
		return http.StatusNotFound
	case AuthTokenExpired, AuthTokenInvalid, AuthRequired:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func (c Code) String() string {
	switch c {
	// SYSTEM
	case Ok:
		return "Operation successful"
	case InternalError:
		return "Internal server error"
	case InvalidParams:
		return "Invalid parameters or input"

	// USER
	case UserNotFound:
		return "User not found"
	case UserAlreadyExists:
		return "User already exists"
	case UserPasswordWrong:
		return "Incorrect password"
	case UserInactive:
		return "User account is inactive"

	// AUTH
	case AuthTokenExpired:
		return "Access token has expired"
	case AuthTokenInvalid:
		return "Invalid access token"
	case AuthRequired:
		return "Authorization is required"

	default:
		return "Unknown error"
	}
}

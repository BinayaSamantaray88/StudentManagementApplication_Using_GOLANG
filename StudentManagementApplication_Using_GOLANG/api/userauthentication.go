package api

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func FetchUser(db2 *sql.DB, username string, password string, w http.ResponseWriter) (User, error) {
	var user User
	err := db2.QueryRow("SELECT username,password FROM User WHERE username=?", username).Scan(&user.Username, &user.Password)
	if err != nil {

		return user, err
	}
	return user, nil
}

func AuthenticateUser(user User, password string) bool {
	{

		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return false
		}
		return true

	}

}

func AuthMiddleware(db2 *sql.DB) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username := r.Header.Get("username")
			password := r.Header.Get("password")
			if username == "" || password == "" {
				statuscode := 1
				message := "UNAUTHORIZED! Please Provide Username and Password"
				sendErrorResponse(w, statuscode, message, nil)
				return
			}
			user, err := FetchUser(db2, username, password, w)
			{
				if err != nil {
					statuscode := 1
					message := "User Not Found"
					sendErrorResponse(w, statuscode, message, nil)
					return
				}
			}
			if AuthenticateUser(user, password) {
				next.ServeHTTP(w, r)

			} else {
				statuscode := 1
				message := "Unauthorized"
				sendErrorResponse(w, statuscode, message, nil)
				return
			}

		})
	}
}

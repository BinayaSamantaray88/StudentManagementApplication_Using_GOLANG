package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(db2 *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newuser User
		json.NewDecoder(r.Body).Decode(&newuser)
		if newuser.Username == "" {
			statuscode := 1
			message := "Bad Credentials"
			sendErrorResponse(w, statuscode, message, nil)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newuser.Password), bcrypt.DefaultCost)
		if err != nil {
			statusCode := 1
			message := err.Error()
			sendErrorResponse(w, statusCode, message, nil)
		}
		newuser.Password = string(hashedPassword)
		_, err = db2.Exec("INSERT INTO User(username,password) VALUES(?,?)", newuser.Username, newuser.Password)
		if err != nil {
			statuscode := 1
			sendErrorResponse(w, statuscode, err.Error(), nil)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		statuscode := 0
		message := "Successfully Created"
		sendSuccessResponse(w, statuscode, message, nil)
	}
}

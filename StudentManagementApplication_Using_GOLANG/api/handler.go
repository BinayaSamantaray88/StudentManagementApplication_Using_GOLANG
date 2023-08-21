package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateStudent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newStudent Student
		_ = json.NewDecoder(r.Body).Decode(&newStudent)

		if newStudent.Name == "" || newStudent.Class <= 0 {
			statuscode := 1
			message := "Bad Credentials"
			sendErrorResponse(w, statuscode, message, nil)

			return

		}
		_, err := db.Exec("INSERT INTO Student(name,class) VALUES(?,?)", newStudent.Name, newStudent.Class)
		if err != nil {
			statuscode := 1
			sendErrorResponse(w, statuscode, err.Error(), nil)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		statuscode := 0
		message := "Successfully Created"
		sendSuccessResponse(w, statuscode, message, nil)
		//json.NewEncoder(w).Encode(newStudent)

	}
}

func GetStudent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])

		if err != nil {
			statuscode := 1
			message := "Invalid ID"
			sendErrorResponse(w, statuscode, message, nil)
			return

		}
		if id <= 0 {
			statuscode := 1
			message := "Invalid ID"
			sendErrorResponse(w, statuscode, message, nil)
			return
		}
		var student Student
		err = db.QueryRow("SELECT id,name,class FROM Student WHERE id=?", id).Scan(&student.ID, &student.Name, &student.Class)
		if err != nil {
			statuscode := 1
			message := "Student not Found"
			sendErrorResponse(w, statuscode, message, nil)
			return
		}
		statuscode := 0
		message := "Data Fetched Successfully"
		sendSuccessResponse(w, statuscode, message, student)
	}
}

func GetALLStudents(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var page, perPage int

		page, _ = strconv.Atoi((r.URL.Query().Get("page")))
		perPage, _ = strconv.Atoi((r.URL.Query().Get("perPage")))
		if page <= 0 || perPage <= 0 {
			var err1, err2 error
			page, err1 = strconv.Atoi(r.Header.Get("page"))
			perPage, err2 = strconv.Atoi(r.Header.Get("perPage"))
			fmt.Println(page, perPage)
			if page <= 0 || perPage <= 0 || err1 != nil || err2 != nil {
				statuscode := 1
				message := "Invalid page  or perPage Values"
				sendErrorResponse(w, statuscode, message, nil)
				return
			}
		}
		fmt.Println(page, perPage)
		offset := (page - 1) * perPage
		query := fmt.Sprintf("SELECT id,name, class FROM Student LIMIT %d OFFSET %d", perPage, offset)
		rows, err := db.Query(query)
		if err != nil {
			statuscode := 1
			message := "Internal Server Error"
			sendErrorResponse(w, statuscode, message, nil)
			return
		}
		defer rows.Close()
		var students []Student
		for rows.Next() {
			var student Student
			err := rows.Scan(&student.ID, &student.Name, &student.Class)
			if err != nil {
				statuscode := 1
				message := err.Error()
				sendErrorResponse(w, statuscode, message, nil)
				return
			}
			students = append(students, student)
		}
		var totalCount int
		db.QueryRow("SELECT COUNT(*) FROM Student").Scan(&totalCount)
		totalPages := int(math.Ceil(float64(totalCount) / float64(perPage)))

		pagination := Pagination{
			TotalStudents:   totalCount,
			TotalPages:      totalPages,
			CurrentPage:     page,
			StudentsPerPage: perPage,
			Students:        students,
		}
		pagination2 := PaginationEr{
			TotalPages: totalPages,
		}
		if page > totalPages {
			statuscode := 1
			message := "Please Check the Page Number!.The pagenumber cannot be greater than the total page number"
			sendErrorResponse(w, statuscode, message, pagination2)
			return
		}
		statuscode := 0
		message := "Data Fetched Successfully"
		sendSuccessResponse(w, statuscode, message, pagination)

	}
}
func DeleteStudent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			message := "Invalid ID"
			statuscode := 1

			sendErrorResponse(w, statuscode, message, nil)
			return

		}

		rows, err := db.Query("SELECT id from Student where id=?", id)
		if err != nil {
			statuscode := 1
			message := err.Error()
			sendErrorResponse(w, statuscode, message, nil)
			return
		}
		defer rows.Close()

		if !rows.Next() {
			statuscode := 1
			message := "Student Not Found"
			sendErrorResponse(w, statuscode, message, nil)
			return
		}

		_, err = db.Exec("DELETE FROM Student WHERE id=?", id)
		if err != nil {
			statuscode := 1
			message := err.Error()
			sendErrorResponse(w, statuscode, message, nil)
			return
		}

		message := "Student Deleted Successfully"
		statuscode := 0
		sendSuccessResponse(w, statuscode, message, nil)

	}
}

func UpdateStudent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			message := "Invalid ID"
			statuscode := 1

			sendErrorResponse(w, statuscode, message, nil)
			return

		}
		var updatedStudent Student
		err = json.NewDecoder(r.Body).Decode(&updatedStudent)
		if err != nil {
			message := err.Error()
			statuscode := 1

			sendErrorResponse(w, statuscode, message, nil)
			return
		}
		rows, err := db.Query("SELECT id FROM Student WHERE id=?", id)
		if err != nil {
			message := err.Error()
			statuscode := 1

			sendErrorResponse(w, statuscode, message, nil)
			return
		}
		defer rows.Close()

		if !rows.Next() {
			message := "Student Not Found"
			statuscode := 1

			sendErrorResponse(w, statuscode, message, nil)

			return

		}

		_, err = db.Exec("UPDATE Student SET name=?, class=?, id=? WHERE id=?", updatedStudent.Name, updatedStudent.Class, updatedStudent.ID, id)
		if err != nil {
			statuscode := 1
			message := err.Error()
			sendErrorResponse(w, statuscode, message, nil)
			return
		}

		message := "Student Updated Successfully"
		statuscode := 0
		sendSuccessResponse(w, statuscode, message, nil)

	}
}

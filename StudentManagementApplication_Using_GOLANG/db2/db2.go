package db2

import "database/sql"

func ConnectDB2() (*sql.DB, error) {
	//Update with your MYSQL database details
	db2, err := sql.Open("mysql", "root:1322@tcp(127.0.0.1:3306)/htmlcourse")
	if err != nil {
		return nil, err
	}
	return db2, nil
}

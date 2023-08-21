package api

type Pagination struct {
	TotalStudents   int         `json:"totalStudents"`
	TotalPages      int         `json:"totalPages"`
	CurrentPage     int         `json:"currentPage"`
	StudentsPerPage int         `json:"studentsPerPage"`
	Students        interface{} `json:"students"`
}

type PaginationEr struct {
	TotalPages int `json:"totalPages"`
}

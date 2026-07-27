package module7

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "studentdb"
)

var db *sql.DB
var gormDB *gorm.DB

type Student struct {
	ID    int    `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func connectDB() {

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)

	var err error

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Connection Pooling
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(10 * time.Minute)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("database/sql Connected")

	gormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("GORM Connected")
}

// ---------------- CREATE ----------------

func createStudent(w http.ResponseWriter, r *http.Request) {

	var s Student

	json.NewDecoder(r.Body).Decode(&s)

	// Transaction Example
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = tx.Exec(
		"INSERT INTO students(name,age,email) VALUES($1,$2,$3)",
		s.Name,
		s.Age,
		s.Email,
	)

	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), 500)
		return
	}

	tx.Commit()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Student Created",
	})
}

// ---------------- READ ALL ----------------

func getStudents(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query(
		"SELECT id,name,age,email FROM students ORDER BY id",
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer rows.Close()

	var students []Student

	for rows.Next() {

		var s Student

		rows.Scan(
			&s.ID,
			&s.Name,
			&s.Age,
			&s.Email,
		)

		students = append(students, s)
	}

	json.NewEncoder(w).Encode(students)
}

// ---------------- READ ONE ----------------

func getStudent(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	var s Student

	err := db.QueryRow(
		"SELECT id,name,age,email FROM students WHERE id=$1",
		id,
	).Scan(
		&s.ID,
		&s.Name,
		&s.Age,
		&s.Email,
	)

	if err != nil {
		http.Error(w, "Student Not Found", 404)
		return
	}

	json.NewEncoder(w).Encode(s)
}

// ---------------- UPDATE ----------------

func updateStudent(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	var s Student

	json.NewDecoder(r.Body).Decode(&s)

	_, err := db.Exec(
		`UPDATE students
		 SET name=$1,age=$2,email=$3
		 WHERE id=$4`,
		s.Name,
		s.Age,
		s.Email,
		id,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Student Updated",
	})
}

// ---------------- DELETE ----------------

func deleteStudent(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]

	_, err := db.Exec(
		"DELETE FROM students WHERE id=$1",
		id,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Student Deleted",
	})
}

// ---------------- GORM DEMO ----------------

func gormStudents(w http.ResponseWriter, r *http.Request) {

	var students []Student

	gormDB.Find(&students)

	json.NewEncoder(w).Encode(students)
}

func Start() {

	connectDB()

	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Student CRUD API")
	})

	router.HandleFunc("/students", createStudent).Methods("POST")
	router.HandleFunc("/students", getStudents).Methods("GET")
	router.HandleFunc("/students/{id}", getStudent).Methods("GET")
	router.HandleFunc("/students/{id}", updateStudent).Methods("PUT")
	router.HandleFunc("/students/{id}", deleteStudent).Methods("DELETE")

	// ORM Basics Example
	router.HandleFunc("/gorm/students", gormStudents).Methods("GET")

	fmt.Println("Server running at http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", router))
}

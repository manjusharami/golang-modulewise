package module6

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

var jwtKey = []byte("my-secret-key")

// =====================
// Models
// =====================

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// =====================
// In Memory Database
// =====================

var students = []Student{
	{
		ID:    1,
		Name:  "Alice",
		Age:   20,
		Email: "alice@example.com",
	},
	{
		ID:    2,
		Name:  "Bob",
		Age:   21,
		Email: "bob@example.com",
	},
}

// =====================
// JWT Login
// =====================

func loginHandler(w http.ResponseWriter, r *http.Request) {

	var request LoginRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Demo user
	if request.Username != "admin" ||
		request.Password != "admin123" {

		http.Error(
			w,
			"Invalid Credentials",
			http.StatusUnauthorized,
		)

		return
	}

	claims := jwt.MapClaims{

		"username": request.Username,

		"exp": time.Now().
			Add(time.Hour).
			Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err :=
		token.SignedString(jwtKey)

	if err != nil {

		http.Error(
			w,
			"Could not create token",
			500,
		)

		return
	}

	json.NewEncoder(w).Encode(
		map[string]string{
			"token": tokenString,
		},
	)
}

// =====================
// JWT Middleware
// =====================

func auth(next http.HandlerFunc) http.HandlerFunc {

	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		header :=
			r.Header.Get("Authorization")

		if header == "" {

			http.Error(
				w,
				"Missing Token",
				http.StatusUnauthorized,
			)

			return
		}

		tokenString :=
			strings.TrimPrefix(
				header,
				"Bearer ",
			)

		token, err :=
			jwt.Parse(
				tokenString,
				func(token *jwt.Token) (
					interface{},
					error,
				) {

					if _, ok :=
						token.Method.(*jwt.SigningMethodHMAC); !ok {

						return nil,
							fmt.Errorf(
								"invalid signing method",
							)
					}

					return jwtKey, nil
				},
			)

		if err != nil ||
			!token.Valid {

			http.Error(
				w,
				"Invalid Token",
				http.StatusUnauthorized,
			)

			return
		}

		next(w, r)

	}
}

// =====================
// Logging Middleware
// =====================

func logging(next http.HandlerFunc) http.HandlerFunc {

	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		log.Println(
			r.Method,
			r.URL.Path,
		)

		next(w, r)
	}
}

// =====================
// GET All Students
// =====================

func getStudents(
	w http.ResponseWriter,
	r *http.Request,
) {

	json.NewEncoder(w).
		Encode(students)
}

// =====================
// GET Student By ID
// =====================

func getStudent(
	w http.ResponseWriter,
	r *http.Request,
) {

	id :=
		mux.Vars(r)["id"]

	studentID, err :=
		strconv.Atoi(id)

	if err != nil {

		http.Error(
			w,
			"Invalid ID",
			400,
		)

		return
	}

	for _, student := range students {

		if student.ID == studentID {

			json.NewEncoder(w).
				Encode(student)

			return
		}
	}

	http.Error(
		w,
		"Student not found",
		404,
	)

}

// =====================
// CREATE Student
// =====================

func createStudent(
	w http.ResponseWriter,
	r *http.Request,
) {

	var student Student

	err :=
		json.NewDecoder(r.Body).
			Decode(&student)

	if err != nil {

		http.Error(
			w,
			"Invalid JSON",
			400,
		)

		return
	}

	if student.Name == "" ||
		student.Email == "" {

		http.Error(
			w,
			"Name and Email required",
			400,
		)

		return
	}

	student.ID =
		len(students) + 1

	students =
		append(
			students,
			student,
		)

	json.NewEncoder(w).
		Encode(student)

}

// =====================
// UPDATE Student
// =====================

func updateStudent(
	w http.ResponseWriter,
	r *http.Request,
) {

	id :=
		mux.Vars(r)["id"]

	studentID, _ :=
		strconv.Atoi(id)

	var updated Student

	json.NewDecoder(r.Body).
		Decode(&updated)

	for index, student := range students {

		if student.ID == studentID {

			updated.ID =
				studentID

			students[index] =
				updated

			json.NewEncoder(w).
				Encode(updated)

			return
		}
	}

	http.Error(
		w,
		"Student not found",
		404,
	)

}

// =====================
// DELETE Student
// =====================

func deleteStudent(
	w http.ResponseWriter,
	r *http.Request,
) {

	id :=
		mux.Vars(r)["id"]

	studentID, _ :=
		strconv.Atoi(id)

	for index, student := range students {

		if student.ID == studentID {

			students =
				append(
					students[:index],
					students[index+1:]...,
				)

			json.NewEncoder(w).
				Encode(
					map[string]string{
						"message": "Student deleted",
					},
				)

			return
		}
	}

	http.Error(
		w,
		"Student not found",
		404,
	)

}

// =====================
// Main
// =====================

func Start() {

	router :=
		mux.NewRouter()

	// Public route
	router.HandleFunc(
		"/login",
		loginHandler,
	).Methods("POST")

	// Protected routes

	router.HandleFunc(
		"/students",
		logging(
			auth(getStudents),
		),
	).Methods("GET")

	router.HandleFunc(
		"/students/{id}",
		logging(
			auth(getStudent),
		),
	).Methods("GET")

	router.HandleFunc(
		"/students",
		logging(
			auth(createStudent),
		),
	).Methods("POST")

	router.HandleFunc(
		"/students/{id}",
		logging(
			auth(updateStudent),
		),
	).Methods("PUT")

	router.HandleFunc(
		"/students/{id}",
		logging(
			auth(deleteStudent),
		),
	).Methods("DELETE")

	fmt.Println(
		"Server running on :8080",
	)

	log.Fatal(
		http.ListenAndServe(
			":8080",
			router,
		),
	)

}

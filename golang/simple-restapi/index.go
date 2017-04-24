package main

import (
	"encoding/json"
	"log"
	"net/http"
	"./model"

	"github.com/gorilla/mux"
)

var employees []model.Employee

func GetEmployee(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, emp := range employees {
		if emp.ID == params["id"] {
			json.NewEncoder(w).Encode(emp)
			return
		}
	}
	json.NewEncoder(w).Encode(&model.Employee{})
}

func GetAllEmployees(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(employees)
}

func AddEmployee(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var emp model.Employee
	_ = json.NewDecoder(req.Body).Decode(&emp)
	emp.ID = params["id"]
	employees = append(employees, emp)
	json.NewEncoder(w).Encode(employees)
}

func FireEmployee(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, emp := range employees {
		if emp.ID == params["id"] {
			employees = append(employees[:index], employees[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(employees)
}

func main() {
	router := mux.NewRouter()
	employees = append(employees, model.Employee{ID: "1", Firstname: "Lincoln", Lastname: "Borrow", Age: 22})
	employees = append(employees, model.Employee{ID: "2", Firstname: "Michael", Lastname: "Scofield", Age: 21})
	router.HandleFunc("/employees", GetAllEmployees).Methods("GET")
	router.HandleFunc("/employees/{id}", GetEmployee).Methods("GET")
	router.HandleFunc("/employees/{id}", AddEmployee).Methods("POST")
	router.HandleFunc("/employees/{id}", FireEmployee).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

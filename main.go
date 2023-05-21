package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type Customer struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     int    `json:"phone"`
	Contacted bool   `json:"contacted"`
}

var customersDb = map[string]Customer{
	"e7847fee-3a20e-455e-b1a51-519ba7851c7": {
		ID:        "e7847fee-3a20e-455e-b1a51-519ba7851c7",
		Name:      "Mahsa",
		Role:      "silver",
		Email:     "mahsa@freedom.com",
		Phone:     981111111110,
		Contacted: true,
	},
	"a787fee-3a4ae-455e-b151-519bd9851c7": {
		ID:        "a787fee-3a4ae-455e-b151-519bd9851c7",
		Name:      "Nika",
		Role:      "bronze",
		Email:     "nika@freedom.com",
		Phone:     981111111111,
		Contacted: true,
	},
	"b7847fee-3a0e-6a5e-b151-519bdb151c7": {
		ID:        "b7847fee-3a0e-6a5e-b151-519bdb151c7",
		Name:      "Kian",
		Role:      "gold",
		Email:     "kian@freedom.com",
		Phone:     981111111112,
		Contacted: false,
	},
	"23a847fee-3a0e-455e-b151-519bdb9851c7": {
		ID:        "23a847fee-3a0e-455e-b151-519bdb9851c7",
		Name:      "Sarina",
		Role:      "gold",
		Email:     "sarina@freedom.com",
		Phone:     981111111113,
		Contacted: false,
	},
	"8c47fee-3a0e-455de-b151-519bdb2851c7": {
		ID:        "8c47fee-3a0e-455de-b151-519bdb2851c7",
		Name:      "Hadis",
		Role:      "gold",
		Email:     "hadis@freedom.com",
		Phone:     981111111114,
		Contacted: true,
	},
}

func handleError(err error) {
	fmt.Println("An error occurred:")
	fmt.Println(err)
}

func handleCustomerNotFound(id string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	_, err := fmt.Fprintf(w, "Failed to find the requested customer with ID \"%s\".", id)
	if err != nil {
		handleError(err)
	}
}

func getCustomers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(customersDb)
	if err != nil {
		handleError(err)
	}
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if customer, ok := customersDb[id]; ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(customer)
		if err != nil {
			handleError(err)
			return
		}
	} else {
		handleCustomerNotFound(id, w)
	}
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	var newCustomer Customer
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(err)
		return
	}
	err = json.Unmarshal(reqBody, &newCustomer)
	if err != nil {
		handleError(err)
		return
	}
	id := uuid.New().String()
	newCustomer.ID = id
	customersDb[id] = newCustomer

	err = json.NewEncoder(w).Encode(newCustomer)
	if err != nil {
		handleError(err)
		return
	}
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	if _, ok := customersDb[id]; ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		delete(customersDb, id)
		err := json.NewEncoder(w).Encode(customersDb)
		if err != nil {
			handleError(err)
			return
		}
	} else {
		handleCustomerNotFound(id, w)
	}
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if _, ok := customersDb[id]; ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		var updatedCustomer Customer
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			handleError(err)
			return
		}
		err = json.Unmarshal(reqBody, &updatedCustomer)
		if err != nil {
			handleError(err)
			return
		}
		updatedCustomer.ID = id
		customersDb[id] = updatedCustomer

		err = json.NewEncoder(w).Encode(updatedCustomer)
		if err != nil {
			handleError(err)
			return
		}
	} else {
		handleCustomerNotFound(id, w)
	}
}

func main() {
	router := mux.NewRouter()
	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/", fileServer)
	router.HandleFunc("/customers", getCustomers).Methods("GET")
	router.HandleFunc("/customers", addCustomer).Methods("POST")
	router.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	router.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")
	router.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")

	fmt.Println("Starting service on port 3000...")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		handleError(err)
	}
}

package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetInfraTenants(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

func GetInfraTenantByID(w http.ResponseWriter, r *http.Request) {

}

func CreateInfraTenant(w http.ResponseWriter, r *http.Request) {

}

func UpdateInfraTenant(w http.ResponseWriter, r *http.Request) {

}

func DeleteInfraTenant(w http.ResponseWriter, r *http.Request) {

}

package api

import (
	"Go-Tracker/model"
	"Go-Tracker/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func IdempotentIncrementStorageApi(w http.ResponseWriter, r *http.Request) {
	if result, err := service.IdempotentIncrement(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&model.ErrorResponse{
				ErrorMessage:  err.Error(),
				TransactionTs: time.Now().Format(time.RFC3339),
			},
		)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

func IncrementStorageWithKeyApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	if result, err := service.IdempotentIncrement(key); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&model.ErrorResponse{
				ErrorMessage:  err.Error(),
				TransactionTs: time.Now().Format(time.RFC3339),
			},
		)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}


func DecrementStorageApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	if result, err := service.Decrement(key); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&model.ErrorResponse{
				ErrorMessage:  err.Error(),
				TransactionTs: time.Now().Format(time.RFC3339),
			})
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

func CheckStorageApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	if result, err := service.Check(key); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&model.ErrorResponse{
				ErrorMessage:  err.Error(),
				TransactionTs: time.Now().Format(time.RFC3339),
			})
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}


func DeleteStorageApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	if result := service.Delete(key); !result {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}
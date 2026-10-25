package main

import (
	"fmt"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/health", healthHandler)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("server error: ", err)
	}

}

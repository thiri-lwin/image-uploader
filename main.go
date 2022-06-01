package main

import (
	"fmt"
	"go_skill_test/controller"
	"log"
	"net/http"
)

func main() {
	serviceConfig := controller.Initialize()

	mux := controller.Router()
	fmt.Printf("Server started on http://localhost:%v", serviceConfig.Port+"\n")
	err := http.ListenAndServe(":"+serviceConfig.Port, mux)
	if err != nil {
		log.Fatal(err)
	}
}

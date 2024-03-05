package main

import (
	"api/src/router"
	"fmt"
	"log"
	"net/http"
)

func main() {

	r := router.Gerar()

	fmt.Println("Servidor aberto: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

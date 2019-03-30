package main

import (
	"./p3"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	router1 := p3.NewRouter()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":" + os.Args[1], router1))
		fmt.Println("os Args > 1")
	} else {
		/* Launch My Server */
		log.Fatal(http.ListenAndServe(":6675", router1))
	}
}

package main

import (
	"./p3"
	"log"
	"net/http"
	"os"
)

func main() {
	router1 := p3.NewRouter()
	router2 := p3.NewRouter()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":" + os.Args[1], router1))
		log.Fatal(http.ListenAndServe(":" + os.Args[1], router2))
	} else {
		/* Launch TA Server */
		go http.ListenAndServe(":6670", router1)
		/* Launch My Server */
		log.Fatal(http.ListenAndServe(":6671", router2))
	}
}

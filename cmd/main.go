package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"testTaskIdeasoft/handlers"
	"testTaskIdeasoft/repository"
)

func main() {
	addr := flag.String("addr", ":8080", "Server address")
	flag.Parse()
	cmds := repository.StartRepositoryManager(12000)


	srv := http.Server{
		Addr:    *addr,
		Handler: handlers.NewHandler(cmds),
	}
	fmt.Fprintf(os.Stderr, "Listen on %v", *addr)
	log.Fatal(srv.ListenAndServe())
}

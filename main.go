package main

import (
	"fmt"
	server "gasbin/srv"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	s := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", "0.0.0.0", 6789),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      GetAPI(),
	}
	srv := server.New(s, log.New(os.Stdout, "", log.LstdFlags), 5*time.Second)
	srv.Run()
}

package main

import (
	"fmt"
	"net/http"
	"webServer/app"
	"webServer/config"
)

func main() {
	fmt.Printf("Server Started with port %v\n", config.Port)

	http.HandleFunc("/", app.BasicHandler)
	http.HandleFunc("/counter", app.CounterHandler)
	http.HandleFunc("/math", app.MathHandler)
	http.HandleFunc("/chess", app.ChessHandler)
	http.HandleFunc("/temp", app.TempHandler)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	if err := http.ListenAndServe(config.Port, nil); err != nil {
		fmt.Println("Error - ", err.Error())
	}
}

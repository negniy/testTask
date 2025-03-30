package main

import (
	"net/http"
	"os"
	"task/config"
	"task/handlers"
	"task/repository"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func loggingMidleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.Logger.Infof("Обработка запроса: %s %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
		config.Logger.Infoln("Обработка запроса завершена")
	})
}

func init() {
	config.LoadLoger()
	config.Logger.Debug("Загрузка логера прошла успешно")
	err := godotenv.Load(".\\config\\conf.env")
	if err != nil {
		config.Logger.Fatal("Ошибка загрузки .env файла")
	}
	config.Logger.Debug("Загрузка .env файла прошла успешно")
	repository.LoadDB()
	config.Logger.Debug("Загрузка бд прошла успешно")
}

func main() {
	router := mux.NewRouter()

	router.Use(loggingMidleware)
	router.HandleFunc("/people", handlers.GetPeople).Methods("GET")
	router.HandleFunc("/people", handlers.CreatePerson).Methods("POST")
	router.HandleFunc("/people", handlers.UpdatePerson).Methods("PUT")
	router.HandleFunc("/people", handlers.DeletePerson).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	config.Logger.Infoln("Запуск сервера на порту", port)
	config.Logger.Fatal(srv.ListenAndServe())
}

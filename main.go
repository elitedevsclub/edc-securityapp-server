package main

import (
	"context"
	"edc-security-app/handlers"
	"edc-security-app/repos"
	"edc-security-app/services"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	port   string
	router *mux.Router
	db     *gorm.DB
}

func NewServer() *Server {
	s := &Server{}
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load env vars; ", err)
	}
	s.port = os.Getenv("PORT")
	if err := s.setupDb(); err != nil {
		log.Fatal(err)
	}
	s.router = mux.NewRouter()
	s.mountRoute()
	return s
}

func (s *Server) run() {
	addr := fmt.Sprintf(":%s", s.port)
	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}
	go func() {
		log.Printf("server started at http://locahost:%s", s.port)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("server error: ", err)
		}
	}()

	graceSignal := make(chan os.Signal, 1)
	signal.Notify(graceSignal, syscall.SIGINT, syscall.SIGTERM)
	<-graceSignal
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("server shutdown error: ", err)
	}
}

func (s *Server) setupDb() error {
	db, err := services.Connect(os.Getenv("DATABASE_URI"))
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Server) mountRoute() {
	repo := repos.NewUserRepository(s.db)
	userHandler := handlers.NewUserHttpHandler(repo)

	router := s.router
	router.HandleFunc("/user/new", userHandler.CreateUserHandler).Methods("POST")
}

func main() {
	NewServer().run()
}

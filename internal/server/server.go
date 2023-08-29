package server

import (
  "fmt"
  "log"
  "net/http"
  "database/sql"
  _ "github.com/lib/pq"
  "github.com/gorilla/mux"
)

type Server struct {
  db *sql.DB
  config *config
  router *mux.Router
}

func NewServer(cfgPath string) (*Server, error){
  config, err := NewConfig(cfgPath)
	if err != nil {
		return nil, err
	}
  connStr := fmt.Sprintf("postgresql://%s:%s@postgres/%s?sslmode=disable", config.User, config.Password, config.DBname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
  return &Server {
    db: db,
    config: config
    router: mu.NewRouter()
  }, nil
}

func (serv *Server) Start() {
  addr := serv.getAddr()
  serv.setRouter()
  fmt.Printf("Server is up on %s\n", addr)
  log.Printf("Server is up on %s\n", addr)
  log.Fatal(http.ListenAndServe(addr, serv.router))
}

func (serv *Server) setRouter() {
  http.Handle("/", serv.router)
}

func (serv *Server) getAddr() string {
  return fmt.Sprintf("%s:%s", s.config.Address, s.config.Port)
}

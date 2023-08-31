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
  config, err := CreateConfig(cfgPath)
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
    config: config,
    router: mux.NewRouter(),
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
  serv.router.HandleFunc("/segments/add", serv.addSegment).Methods("POST")
  serv.router.HandleFunc("/segments/delete", serv.delSegment).Methods("POST")
  serv.router.HandleFunc("/users/add", serv.addUser).Methods("POST")
  serv.router.HandleFunc("/users/get", serv.getSegments).Methods("POST")
}

func (serv *Server) getAddr() string {
  return fmt.Sprintf("%s:%s", serv.config.Address, serv.config.Port)
}

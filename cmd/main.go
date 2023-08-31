package main

import (
    "log"
  "https://github.com/DEEMMOONS/avito_backend_internship/internal/server"
)

func main() {
  server, err := server.NewServer("app/configs/config.json")
  if err != nil {
    log.Fatal(err)
  }
  server.Start()
}

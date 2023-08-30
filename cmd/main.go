package main

import (
    "log"
    "internal/server"
)

func main() {
  server, err := NewServer("../../configs")
  if err != nil {
    log.Fatal(err)
  }
  server.Start()
}

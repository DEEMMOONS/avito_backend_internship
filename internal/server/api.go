package server

import (
  "github.com/DEEMMOONS/avito_backend_internship/tree/develop/internal/database"
  "net/http"
  "os"
  "fmt"
)

func (serv *Server) addSegment(w http.ResponseWriter, r *http.Request) {
  inputData := segment{}
  status := getRequestData(w, r, &inputData)
  if !status {
    return
  }
  check, err := database.CheckSeg(serv.db, inputData.Name)
  if err != nil {
    internalServerError(w, err)
    return
  }
  if !check {
    invalidData(w)
    return
  }
  err2 := database.AddSegment(serv.db, inputData.Name)
  if err2 != nil {
    internalServerError(w, err2)
    return
  }
  respondSuccess(w)
}

func (serv *Server) delSegment(w http.ResponseWriter, r *http.Request) {
  inputData := segment{}
  status := getRequestData(w, r, &inputData)
  if !status {
    return
  }
  check, err := database.CheckSeg(serv.db, inputData.Name)
  if err != nil {
    internalServerError(w, err)
    return
  }
  if check {
    invalidData(w)
    return
  }
  err2 := database.DelSegment(serv.db, inputData.Name)
  if err2 != nil {
    internalServerError(w, err2)
    return
  }
  respondSuccess(w)
}

func (serv *Server) addUser(w http.ResponseWriter, r *http.Request) {
  inputData := user{}
  status := getRequestData(w, r, &inputData)
  if !status {
    return
  }
  for _, addSeg := range inputData.AddSegs {
    userStat, segStat, err := database.CheckStatus(serv.db, addSeg, inputData.Id)
    if err != nil {
      internalServerError(w, err)
      return
    }
    if !userStat || segStat {
      invalidData(w)
      return
    }
  }
  for _, delSeg := range inputData.DelSegs {
    userStat, err := database.CheckUser(serv.db, inputData.Id, delSeg)
    if err != nil {
      internalServerError(w, err)
      return
    }
    if userStat {
      invalidData(w)
      return
    }
  }
  err := database.AddUserSegs(serv.db, inputData.Id, inputData.AddSegs, inputData.DelTime)
  if err != nil {
    internalServerError(w, err)
    return
  }
  err2 := database.DelUserSegs(serv.db, inputData.Id, inputData.DelSegs)
  if err2 != nil {
    internalServerError(w, err2)
    return
  }
  respondSuccess(w)
}

func (serv *Server) getSegments(w http.ResponseWriter, r *http.Request) {
  inputData := data{}
  status := getRequestData(w, r, &inputData)
  if !status {
    return
  }
  result, err := database.GetSegments(serv.db, inputData.Id)
  if err != nil {
    internalServerError(w, err)
    return
  }
  makeRespond(w, http.StatusOK, jsonRespond(result))
}

func (serv *Server) getUserStat(w http.ResponseWriter, r *http.Request) {
  inputData := data{}
  status := getRequestData(w, r, &inputData)
  if !status {
    return
  }
  result, err := database.GetUserStat(serv.db, inputData.Id, inputData.Time)
  if err != nil {
    internalServerError(w, err)
    return
  }
  errCSV := makeCSV("temp.csv", result)
  if errCSV != nil {
    internalServerError(w, errCSV)
    return
  }
  defer os.Remove("temp.csv")
  link := fmt.Sprintf("http://%s/%s", r.Host, "temp.csv")
  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, `{"csv_link": "%s"}`, link)
  respondSuccess(w)
}

func makeCSV(name string, content []string) error {
  file, err := os.Create(name)
  if err != nil {
    return err
  }
  defer file.Close()
  for _, str := range content {
    _, err2 := file.WriteString(str + "\n")
    if err2 != nil {
    return err2
    }
  }
  return nil
}

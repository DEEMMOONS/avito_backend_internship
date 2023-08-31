package server

import (
  "internal/database"
  "io"
  "net/http"
  "time"
)

func (serv *Server) addSegment(w http.ResponseWriter, r *http.Request) {
  inputData := segment{}
  status := getRequestData(w, r, &inputData)
  if !status {
    return
  }
  check, err := CheckSeg(serv.db, inputData.name)
  if err != nil {
    internalServerError(w, err)
    return
  }
  if !check {
    invalidData(w)
    return
  }
  err := AddSegment(serv.db, inputData.name)
  if err != nil {
    internalServerError(w, err)
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
  check, err := CheckSeg(serv.db, inputData.name)
  if err != nil {
    internalServerError(w, err)
    return
  }
  if check {
    invalidData(w)
    return
  }
  err := DelSegment(serv.db, inputData.name)
  if err != nil {
    internalServerError(w, err)
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
  for _, addSeg := range inputiData.addSegs {
    userStat, segStat, err := CheckStatus(serv.db, inputData.id, addSeg)
    if err != nil {
      internalServerError(w, err)
      return
    }
    if !userStat || segStat {
      invalidData(w)
      return
    }
  }
  for _, delSeg := range inputData.delSegs {
    userStat, segStat, err := CheckStatus(serv.db, inputData.id, delSeg)
    if err != nil {
      internalServerError(w, err)
      return
    }
    if userStat || segStat {
      invalidData(w)
      return
    }
  }
  err := AddUserSegs(serv.db, inputData.id, inputData.addSegs, inputData.delTime)
  if err != nil {
    internalServerError(w, err)
    return
  }
  err := DelUserSegs(serv.db, inputData.id, inputData.delSegs)
  if err != nil {
    internalServerError(w, err)
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
  result, err := GetSegments(serv.db, inputData.id)
  if err != nil {
    internalServerError(w, err)
    return
  }
  makeRespond(w, http.StatusOK, jsonResult(result))
}

func (serv *Server) getUserStat(w http.ResponseWriter, r *http.Request) {
  inputData := data{}
  status := getRequestData(w, r, &inputData)
  if !status {
    return
  }
  result, err := GetUserStat(serv.db, inputData.id, inputData.checkInterval)
  if err != nil {
    internalServerError(w, err)
    return
  }
  makeRespond(w, http.StatusOK, jsonResult(result))
}

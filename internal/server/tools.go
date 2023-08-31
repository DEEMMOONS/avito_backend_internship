package server

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "io"
)

type segment struct {
  Name string `json:"slug"`
}

type data struct {
  Id int `json:"id"`
  Time  string `json:"period,omitempty"`
}

type user struct {
  Id int `json:"id"`
  AddSegs []string `json:"add,omitempty"`
  DelSegs []string `json:"del,omitempty"`
  DelTime string `json:"delete_time,omitempty"`
}

func getRequestData(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		internalServerError(w, err)
		return false
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		invalidData(w)
		return false
	}
	return true
}

func respondSuccess(w http.ResponseWriter) {
	makeRespond(w, http.StatusOK, jsonRespond("success"))
}

func internalServerError(w http.ResponseWriter, err error) {
	makeRespond(w, http.StatusInternalServerError, jsonRespond("internal server error"))
	log.Println(err)
}

func invalidData(w http.ResponseWriter) {
	makeRespond(w, http.StatusBadRequest, jsonRespond("invalid data"))
}

func jsonRespond(respond string) []byte {
	return []byte(fmt.Sprintf(`{"result": "%s"}`, respond))
}

func makeRespond(w http.ResponseWriter, code int, data []byte) {
  w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(data)
	if err != nil {
		log.Println(err)
	}
}

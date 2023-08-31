package server

import (
  "encoding/json"
  "fmt"
  "time"
  "log"
  "net/http"
)

type segment struct {
  name string `json:"slug"`
}

type data struct {
  id int `json:"id"`
  time  time.Time `json:"period,omitempty"`
}

type user struct {
  id int `json:"id"`
  addSegs []string `json:"add,omitempty"`
  delSegs []string `json:"del,omitempty"`
  deltime time.Time `json:"delete_time,omitempty"`
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
	makeJsonRespond(w, http.StatusBadRequest, jsonRespond("invalid data"))
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

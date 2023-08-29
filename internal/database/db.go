package database

import (
  "database/sql"
  _ "github.com/lib/pq"
)

func (serv *Server) checkSeg(segName string) (bool, error) {
  var status bool = false
  err := serv.db.QueryRow("SELECT name FROM segments WHERE name = $1", segName).Scan(segCheck)
  if err != nil {
    if err == sql.ErrNoRows {
      status = true
    } else {
      return status, err
    }
  }
  return status, nil
}

func (serv *Server) addSegment(segName string) error {
  check, err := serv.checkSeg(segName)
  if err != nil {
    return err
  }
  if (check) {
    _, err := serv.db.Exec("INSERT INTO segments VALUES ($1)", segName)
    if err != nil {
      return err
    }
  }
  return nil
}

func (serv *Server) removeSegment(segName string) error {
  check, err := serv.checkSeg(segName)
  if err != nil {
    return err
  }
  if (!check) {
    _, err := serv.db.Exec("DELETE FROM segments WHERE name = $1", segName)
    if err != nil {
      return err
    }
  }
  return nil
}

func (serv *Server) checkUser(id int, segName string) (bool, error) {
  var status bool = true
  rows, err := serv.db.Query("SELECT name FROM users WHERE id = $1 AND segment = $2 AND delete_at IS NULL", id, segName)
  defer rows.Close()
  if err != nil {
    return status, err
  }
  if rows.Next() {
    status = false
  }
  return status, nil
}

func (serv *Server) addUser(id int, addSegs []string, delSegs []string) error {
  for _, addSeg := range addSegs {
    if !checkUser(db, id, addSeg) && !checkSeg(addSeg) {
      _, err := serv.db.Exec("INSERT INTO users (id, segment, create_at, delete_at) VALUES ($1, $2, CURRENT_TIMESTAMP, $4)", id, addSeg, nil)
      if err != nil {
        return err
      }
    }
  }
  for _, delSeg := range delSegs {
    if checkUser(db, id, delSeg) && !checkSeg(delSeg) {
      _, err := serv.db.Exec("UPDATE users SET deleate_at = CURRENT_TIMESTAMP WHERE id = $2 AND segment = $3", id, delSeg)
      if err != nil {
        return err
      }
    }
  }
  return nil
}

func (serv *Server) getSegments(id int) (string, error) {
  rows, err := serv.db.Query("SELECT segment FROM users WHERE id = $1", id)
  if err != nil {
    return "", err
  }
  defer rows.Close()
  var segs string = ""
  for rows.Next() {
    var buf string
    err := rows.Scan(&buf)
    if err != nil {
      return "", err
    }
    segs += buf
  }
  return segs, nil
}

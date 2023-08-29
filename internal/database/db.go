package database

import (
  "database/sql"
  _ "github.com/lib/pq"
  "time"
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
  var status bool = false
  err := serv.db.QueryRow("SELECT name FROM users WHERE id = $1 AND segment = $2 AND (delete_at IS NULL OR delete_at > CURRENT_TIMESTAMP)", id, segName).Scan(segCheck)
  if err != nil {
    if err == sql.ErrNoRows {
      status = true
    } else {
      return status, err
    }
  }
  return status, nil
}

func (serv *Server) checkStatus (segName string, id int) (bool, bool, error) {
  segStat, errSeg := serv.checkSeg(segName)
  if errSeg != nil {
      return false, false, errSeg
  }
  userStat, errUser := serv.checkUser(id, segName)
  if errUser != nil {
      return false, false, errUser
  }
  return userStat, segStat, nil
}

func (serv *Server) addUser(id int, addSegs []string, delSegs []string, delTime Time) error {
  for _, addSeg := range addSegs {
    userStat, segStat, err := checkStatus(id, addSeg)
    if err != nil {
      return err
    }
    if userStat && !segStat(addSeg) {
      _, err := serv.db.Exec("INSERT INTO users (id, segment, create_at, delete_at) VALUES ($1, $2, CURRENT_TIMESTAMP, $3)", id, addSeg, delTime)
      if err != nil {
        return err
      }
    }
  }
  for _, delSeg := range delSegs {
    userStat, segStat, err := checkStatus(id, addSeg)
    if err != nil {
      return err
    }
    if !userStat && !segStat {
      _, err := serv.db.Exec("UPDATE users SET deleate_at = CURRENT_TIMESTAMP WHERE id = $1 AND segment = $2", id, delSeg)
      if err != nil {
        return err
      }
    }
  }
  return nil
}

func (serv *Server) getSegments(id int) (string, error) {
  rows, err := serv.db.Query("SELECT segment FROM users WHERE id = $1 AND (delete_at IS NULL OR delete_at > CURRENT_TIMESTAMP)", id)
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

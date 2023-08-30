package database

import (
  "database/sql"
  _ "github.com/lib/pq"
  "time"
  "strconv"
)

func checkSeg(db *sql.DB, segName string) (bool, error) {
  var status bool = false
  err := db.QueryRow("SELECT name FROM segments WHERE name = $1", segName).Scan(segCheck)
  if err != nil {
    if err == sql.ErrNoRows {
      status = true
    } else {
      return status, err
    }
  }
  return status, nil
}

func AddSegment(db *sql.DB, segName string) error {
  check, err := checkSeg(db, segName)
  if err != nil {
    return err
  }
  if (check) {
    _, err := db.Exec("INSERT INTO segments (name) VALUES ($1)", segName)
    if err != nil {
      return err
    }
  }
  return nil
}

func DelSegment(db *sql.Db, segName string) error {
  check, err := checkSeg(db, segName)
  if err != nil {
    return err
  }
  if (!check) {
    _, err := db.Exec("DELETE FROM segments WHERE name = $1", segName)
    if err != nil {
      return err
    }
  }
  return nil
}

func checkUser(db *sql.DB, id int, segName string) (bool, error) {
  var status bool = false
  err := db.QueryRow("SELECT name FROM users WHERE id = $1 AND segment = $2 AND (delete_at IS NULL OR delete_at > CURRENT_TIMESTAMP)", id, segName).Scan(segCheck)
  if err != nil {
    if err == sql.ErrNoRows {
      status = true
    } else {
      return status, err
    }
  }
  return status, nil
}

func checkStatus (db *sql.DB, segName string, id int) (bool, bool, error) {
  segStat, errSeg := checkSeg(db, segName)
  if errSeg != nil {
      return false, false, errSeg
  }
  userStat, errUser := checkUser(db, id, segName)
  if errUser != nil {
      return false, false, errUser
  }
  return userStat, segStat, nil
}

func AddUser(db *sql.DB, id int, addSegs []string, delSegs []string, delTime time.Time) error {
  for _, addSeg := range addSegs {
    userStat, segStat, err := checkStatus(db, id, addSeg)
    if err != nil {
      return err
    }
    if userStat && !segStat(addSeg) {
      _, err := db.Exec("INSERT INTO users (id, segment, create_at, delete_at) VALUES ($1, $2, CURRENT_TIMESTAMP, $3)", id, addSeg, delTime)
      if err != nil {
        return err
      }
    }
  }
  for _, delSeg := range delSegs {
    userStat, segStat, err := checkStatus(db, id, addSeg)
    if err != nil {
      return err
    }
    if !userStat && !segStat {
      _, err := db.Exec("UPDATE users SET deleate_at = CURRENT_TIMESTAMP WHERE id = $1 AND segment = $2", id, delSeg)
      if err != nil {
        return err
      }
    }
  }
  return nil
}

func GetSegments(db *sql.DB, id int) (string, error) {
  rows, err := db.Query("SELECT segment FROM users WHERE id = $1 AND (delete_at IS NULL OR delete_at > CURRENT_TIMESTAMP)", id)
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

func GetUserStat(db *sql.DB, id int, checkInterval time.Time) ([]string, error) {
  rows, err := db.Query("SELECT segment, create_at, delete_at FROM users WHERE id = $1 AND (create_at > $2 OR delete_at > $2)", id, checkInterval)
  if err != nil {
    return "", err
  }
  defer rows.Close()
  var stat []string
  for rows.Next() {
    var bufSeg string
    var bufAdd, bufDel time.Time
    err := rows.Scan(&bufSeg, &bufAdd, &bufDel)
    if bufAdd.After(checkInterval) {
      stat = append(stat, strconv.Itoa(id) + ";" + bufSeg + ";Add;" + bufAdd.Format("2006-01-02 15:04:05"))
    }
    if bufDel.After(checkInterval) {
      stat = append(stat, strconv.Itoa(id) + ";" + bufSeg + ";Delete;" + bufDel.Format("2006-01-02 15:04:05"))
    }
  }
  return stat, nil
}

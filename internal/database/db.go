package database

import (
  "database/sql"
  _ "github.com/lib/pq"
  "time"
  "strconv"
  "fmt"
)

func CheckSeg(db *sql.DB, segName string) (bool, error) {
  var status bool = false
  var segCheck string
  err := db.QueryRow("SELECT name FROM segments WHERE name = $1", segName).Scan(&segCheck)
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
  _, err := db.Exec("INSERT INTO segments (name) VALUES ($1)", segName)
  if err != nil {
    return err
  }
  return nil
}

func DelSegment(db *sql.DB, segName string) error {
  _, err := db.Exec("DELETE FROM segments WHERE name = $1", segName)
  if err != nil {
    return err
  }
  return nil
}

func CheckUser(db *sql.DB, id int, segName string) (bool, error) {
  var status bool = false
  var segCheck string
  err := db.QueryRow("SELECT segment FROM users WHERE id = $1 AND segment = $2 AND (delete_at IS NULL OR delete_at > CURRENT_TIMESTAMP)", id, segName).Scan(&segCheck)
  if err != nil {
    if err == sql.ErrNoRows {
      status = true
    } else {
      return status, err
    }
  }
  return status, nil
}

func CheckStatus (db *sql.DB, segName string, id int) (bool, bool, error) {
  segStat, errSeg := CheckSeg(db, segName)
  if errSeg != nil {
      return false, false, errSeg
  }
  userStat, errUser := CheckUser(db, id, segName)
  if errUser != nil {
      return false, false, errUser
  }
  return userStat, segStat, nil
}

func AddUserSegs(db *sql.DB, id int, addSegs []string, delTimeStr string) error {
  delTime, _ := time.Parse(time.RFC3339, delTimeStr)
  fmt.Println(delTime)
  for _, addSeg := range addSegs {
    if delTime.IsZero() {
      _, err := db.Exec("INSERT INTO users (id, segment, create_at, delete_at) VALUES ($1, $2, $3, $4)", id, addSeg, time.Now().Round(time.Second), nil)
      if err != nil {
        return err
      }
    } else {
      _, err := db.Exec("INSERT INTO users (id, segment, create_at, delete_at) VALUES ($1, $2, $3, $4)", id, addSeg, time.Now().Round(time.Second), delTime)
      if err != nil {
        return err
      }
    }
  }
  return nil
}
func DelUserSegs(db *sql.DB, id int, delSegs []string) error {
  for _, delSeg := range delSegs {
    _, err := db.Exec("UPDATE users SET delete_at = CURRENT_TIMESTAMP WHERE id = $1 AND segment = $2", id, delSeg)
    if err != nil {
      return err
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
  var segs string = "id: " + strconv.Itoa(id) + "; segments: "
  for rows.Next() {
    var buf string
    err := rows.Scan(&buf)
    if err != nil {
      return "", err
    }
    segs += buf + ","
  }
  return segs, nil
}

func GetUserStat(db *sql.DB, id int, checkIntervalStr string) ([]string, error) {
  checkInterval, _ := time.Parse(time.RFC3339, checkIntervalStr)
  rows, err := db.Query("SELECT segment, create_at, delete_at FROM users WHERE id = $1 AND (create_at > $2 OR delete_at > $2)", id, checkInterval)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  var stat []string
  for rows.Next() {
    var bufSeg, bufAddStr, bufDelStr string
    bufAdd, _ := time.Parse("2006-01-02 15:04:05", bufAddStr)
    bufDel, _ := time.Parse("2006-01-02 15:04:05", bufDelStr)
    err := rows.Scan(&bufSeg, &bufAdd, &bufDel)
    if err != nil {
      return nil, err
    }
    if bufAdd.After(checkInterval) {
      stat = append(stat, strconv.Itoa(id) + ";" + bufSeg + ";Add;" + bufAdd.Format("2006-01-02 15:04:05"))
    }
    if !bufDel.IsZero() && bufDel.After(checkInterval) {
      stat = append(stat, strconv.Itoa(id) + ";" + bufSeg + ";Delete;" + bufDel.Format("2006-01-02 15:04:05"))
    }
  }
  return stat, nil
}

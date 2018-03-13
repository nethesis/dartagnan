package models

import (
  "encoding/json"
  "database/sql/driver"
  "errors"
)

type NotificationMap map[string]interface{}

func (p NotificationMap) Value() (driver.Value, error) {
  j, err := json.Marshal(p)
  return j, err
}

func (p *NotificationMap) Scan(src interface{}) error {
  source, ok := src.([]byte)
  if !ok {
    return errors.New("Type assertion .([]byte) failed.")
  }

  var i interface{}
  if err := json.Unmarshal(source, &i); err != nil {
    return err
  }

  *p, ok = i.(map[string]interface{})
  if !ok {
    return errors.New("Type assertion .(map[string]interface{}) failed.")
  }

  return nil
}

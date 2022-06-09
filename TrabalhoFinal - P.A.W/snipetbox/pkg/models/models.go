package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record Found")

type Snippet struct{
  ID int
  Nome string
  Contato string
  Entrada time.Time
  Saida time.Time
}
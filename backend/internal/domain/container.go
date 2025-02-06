package domain

import (
	//"errors"
	"time"
)

var (
// ErrShortLinkNotFound = errors.New("short link not present in db")
)

type Container struct {
	IP          string    `json:"ip"`
	PingTime    time.Time `json:"pingtime"`
	SuccessDate time.Time `json:"successdate"`
}

type ListContainer []Container
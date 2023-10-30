package server

import "time"

type Options struct {
	Host, Port     string
	ConnectTimeout time.Duration
}

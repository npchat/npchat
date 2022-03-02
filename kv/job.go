package kv

import "github.com/intob/rocketkv/protocol"

type Job struct {
	Msg  *protocol.Msg
	Resp chan *protocol.Msg
}

func NewJob(msg *protocol.Msg) *Job {
	return &Job{
		Msg:  msg,
		Resp: make(chan *protocol.Msg),
	}
}

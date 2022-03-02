package kv

import (
	"errors"
	"time"

	"github.com/intob/rocketkv/protocol"
)

const TIMEOUT = time.Second * 5

// Helper to get a value with a timeout
func (st *Store) Get(key string) (*protocol.Msg, error) {
	j := NewJob(&protocol.Msg{
		Op:  protocol.OpGet,
		Key: key,
	})
	return st.doWithTimeout(j, TIMEOUT)
}

// Helper to set a value with a timeout
//
// ttl: seconds until key expires, if < 1, key will not expire
func (st *Store) Set(key string, value []byte, ttl int64) error {
	j := NewJob(&protocol.Msg{
		Op:    protocol.OpSetAck,
		Key:   key,
		Value: value,
	})
	if ttl > 0 {
		j.Msg.Expires = time.Now().Unix() + ttl
	}
	resp, err := st.doWithTimeout(j, TIMEOUT)
	if err != nil {
		return err
	}
	if resp.Status != protocol.StatusOk {
		return errors.New("did not recieve OK response")
	}
	return nil
}

func (st *Store) doWithTimeout(j *Job, timeout time.Duration) (*protocol.Msg, error) {
	st.StartJob(j)
	giveUp := make(chan bool)
	go func() {
		time.Sleep(timeout)
		giveUp <- true
	}()
	select {
	case resp := <-j.Resp:
		return resp, nil
	case <-giveUp:
		return nil, errors.New("request timed out")
	}
}

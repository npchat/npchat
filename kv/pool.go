package kv

import (
	"sync/atomic"
	"time"

	"github.com/intob/npchat/cfg"
	"github.com/spf13/viper"
)

// Maintains a pool of connections,
// handles jobs sent on the jobs channel
type Store struct {
	min   uint32
	count *uint32
	jobs  chan *Job
	cfg   RocketCfg
}

type RocketCfg struct {
	Network    string
	Address    string
	AuthSecret string
	CertFile   string
	KeyFile    string
}

func NewStore() *Store {
	count := uint32(0)
	st := &Store{
		count: &count,
		min:   viper.GetUint32(cfg.ROCKET_WORKERS_MIN),
		jobs:  make(chan *Job),
		cfg: RocketCfg{
			Network:    viper.GetString(cfg.ROCKET_NET),
			Address:    viper.GetString(cfg.ROCKET_ADDRESS),
			AuthSecret: viper.GetString(cfg.ROCKET_AUTH),
			CertFile:   viper.GetString(cfg.ROCKET_TLS_CERTFILE),
			KeyFile:    viper.GetString(cfg.ROCKET_TLS_KEYFILE),
		},
	}
	go st.startWorkers()
	return st
}

func (st *Store) StartJob(job *Job) {
	st.jobs <- job
}

func (st *Store) startWorkers() {
	for i := uint32(0); i < st.min; i++ {
		go st.startWorker()
	}

	// periodically check we have enough workers
	for {
		time.Sleep(time.Second)

		// atomaically load count
		c := atomic.LoadUint32(st.count)

		if c >= st.min {
			continue
		}

		go st.startWorker()
	}
}

package kv

import (
	"sync/atomic"
	"time"

	"github.com/intob/npchat/cfg"
	"github.com/spf13/viper"
)

// Maintains a pool of Agents
//
// If no agents are available & count < max,
// Pool spawns a new one.
type Pool struct {
	Min   uint32
	Count *uint32
	Jobs  chan Job
	Cfg   RocketCfg
}

type RocketCfg struct {
	Network    string
	Address    string
	AuthSecret string
	CertFile   string
	KeyFile    string
}

func NewPool() *Pool {
	count := uint32(0)
	p := &Pool{
		Count: &count,
		Min:   viper.GetUint32(cfg.ROCKET_WORKERS_MIN),
		Jobs:  make(chan Job),
		Cfg: RocketCfg{
			Network:    viper.GetString(cfg.ROCKET_NET),
			Address:    viper.GetString(cfg.ROCKET_ADDRESS),
			AuthSecret: viper.GetString(cfg.ROCKET_AUTH),
			CertFile:   viper.GetString(cfg.ROCKET_TLS_CERTFILE),
			KeyFile:    viper.GetString(cfg.ROCKET_TLS_KEYFILE),
		},
	}
	go p.StartWorkers()
	return p
}

func (p *Pool) StartWorkers() {
	for i := uint32(0); i < p.Min; i++ {
		go p.StartWorker()
	}

	// periodically check we have enough workers
	for {
		time.Sleep(time.Second)

		// atomaically load count
		c := atomic.LoadUint32(p.Count)

		if c >= p.Min {
			continue
		}

		go p.StartWorker()
	}
}

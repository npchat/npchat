package kv

import (
	"bufio"
	"fmt"
	"sync/atomic"

	"github.com/intob/rocketkv/protocol"
	"github.com/intob/rocketkv/util"
)

type Job struct {
	Msg  protocol.Msg
	Resp chan protocol.Msg
}

func (p *Pool) StartWorker() {

	// increment worker count
	atomic.AddUint32(p.Count, 1)

	// decrement worker count
	defer atomic.AddUint32(p.Count, ^uint32(0))

	defer recoverWorker()

	// get conn
	conn, err := util.GetConn(p.Cfg.Network, p.Cfg.Address,
		p.Cfg.CertFile, p.Cfg.KeyFile)
	if err != nil {
		panic(err)
	}
	fmt.Println("worker connected")

	scanner := bufio.NewScanner(conn)
	scanner.Split(protocol.SplitPlusEnd)

	// process jobs
	for job := range p.Jobs {
		// encode msg
		m, err := protocol.EncodeMsg(&job.Msg)
		p.checkJobError(err, job)

		// write to conn
		_, err = conn.Write(m)
		p.checkJobError(err, job)

		// read & dispatch response
		scanner.Scan()
		respBytes := scanner.Bytes()
		resp, err := protocol.DecodeMsg(respBytes)
		p.checkJobError(err, job)
		job.Resp <- *resp
	}

	fmt.Println("worker is done :)")
}

// If something goes wrong,
// send job back to queue & panic
func (p *Pool) checkJobError(err error, job Job) {
	if err != nil {
		p.Jobs <- job
		panic(err)
	}
}

func recoverWorker() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}

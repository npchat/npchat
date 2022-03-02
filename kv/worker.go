package kv

import (
	"bufio"
	"fmt"
	"sync/atomic"

	"github.com/intob/rocketkv/protocol"
	"github.com/intob/rocketkv/util"
)

func (st *Store) startWorker() {
	// increment worker count
	atomic.AddUint32(st.count, 1)

	// decrement worker count
	defer atomic.AddUint32(st.count, ^uint32(0))

	defer recoverWorker()

	// get conn
	conn, err := util.GetConn(st.cfg.Network, st.cfg.Address,
		st.cfg.CertFile, st.cfg.KeyFile)
	if err != nil {
		panic(err)
	}
	fmt.Println("worker connected")

	scanner := bufio.NewScanner(conn)
	scanner.Split(protocol.SplitPlusEnd)

	// process jobs
	for job := range st.jobs {
		// encode msg
		m, err := protocol.EncodeMsg(job.Msg)
		st.checkJobError(err, job)

		// write to conn
		_, err = conn.Write(m)
		st.checkJobError(err, job)

		// read & dispatch response
		scanner.Scan()
		respBytes := scanner.Bytes()
		resp, err := protocol.DecodeMsg(respBytes)
		st.checkJobError(err, job)
		job.Resp <- resp
	}

	fmt.Println("worker is done :)")
}

// If something goes wrong,
// send job back to queue & panic
func (st *Store) checkJobError(err error, job *Job) {
	if err != nil {
		st.jobs <- job
		panic(err)
	}
}

func recoverWorker() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}

package status

import (
	"fmt"
	"net/http"

	"github.com/intob/npchat/kv"
	"github.com/intob/rocketkv/protocol"
)

func HandleGetHealth(w http.ResponseWriter, r *http.Request, p *kv.Pool) {
	respChan := make(chan protocol.Msg)
	job := kv.Job{
		Msg: protocol.Msg{
			Op: protocol.OpPing,
		},
		Resp: respChan,
	}
	p.Jobs <- job

	resp := <-respChan
	fmt.Println("response:", protocol.MapStatus()[resp.Status])

	w.Write([]byte(protocol.MapStatus()[resp.Status]))
}

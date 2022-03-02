package status

import (
	"fmt"
	"net/http"

	"github.com/intob/npchat/kv"
	"github.com/intob/rocketkv/protocol"
)

func HandleGetHealth(w http.ResponseWriter, r *http.Request, st *kv.Store) {
	job := kv.NewJob(&protocol.Msg{
		Op: protocol.OpPing,
	})
	st.StartJob(job)

	resp := <-job.Resp
	fmt.Println("response:", protocol.MapStatus()[resp.Status])

	w.Write([]byte(protocol.MapStatus()[resp.Status]))
}

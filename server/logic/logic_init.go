package logic

import (
	"io"
	"log"
	pssmpb "pssm/gen"
)

type PSSM struct {
	pssmpb.UnimplementedPssmServer
}

// ServiceLinker implements pssmpb.PssmServer.
func (*PSSM) ServiceLinker(stream pssmpb.Pssm_ServiceLinkerServer) error {

	for {

		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Println(err, "%v closing", stream)
				return err
			}
			log.Println(err)
			continue
		}

		switch msg := req.Message.(type) {
		case *pssmpb.ServiceLinkerReq_Msg:
			log.Println(msg)
		case *pssmpb.ServiceLinkerReq_X:
			log.Println(msg)
		default:
			// 不支持的类型或未设置的oneof字段
			log.Printf("Unknown message type or oneof field not set %v", msg)
		}

		resp := &pssmpb.ServiceLinkerRes{}
		stream.Send(resp)

	}

	// panic("unimplemented")
}

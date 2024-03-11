package main

import (
	"fmt"
	"log"
	"net"
	pssmpb "pssm/gen"
	"pssm/server/logic"

	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 0)) //开启监听
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer() //新建一个grpc服务

	pssmpb.RegisterPssmServer(s, &logic.PSSM{}) // notify_logic 服务注册

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// rpc
	// user_id clear (resource)
	// guest_id clear (resource)
}

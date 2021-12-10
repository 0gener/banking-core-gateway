package client

import (
	"log"

	"github.com/0gener/banking-core-accounts/proto"
	"google.golang.org/grpc"
)

func NewAccountsClient() (proto.AccountsServiceClient, *grpc.ClientConn) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial("localhost:5000", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	return proto.NewAccountsServiceClient(conn), conn
}

package client

import (
	"log"

	"github.com/0gener/banking-core-accounts/proto"
	"google.golang.org/grpc"
)

type AccountsClientOptions struct {
	Url string
}

func NewAccountsClient(options AccountsClientOptions) (proto.AccountsServiceClient, *grpc.ClientConn) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(options.Url, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	return proto.NewAccountsServiceClient(conn), conn
}

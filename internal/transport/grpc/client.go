package grpcClient

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/whoami00911/gRPC-server/pkg/grpcPb"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	conn    *grpc.ClientConn
	gClient grpcPb.LogServiceClient
}

func NewClient() (*GrpcClient, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", viper.GetString("logdb.ip"), viper.GetString("logdb.port")))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return &GrpcClient{
		conn:    conn,
		gClient: grpcPb.NewLogServiceClient(conn),
	}, nil
}

func (g *GrpcClient) SendLogRequest() {

}

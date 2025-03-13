package grpcClient

import (
	"context"
	"fmt"
	"webPractice1/pkg/logger"

	"github.com/spf13/viper"
	"github.com/whoami00911/gRPC-server/pkg/grpcPb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcClient struct {
	conn    *grpc.ClientConn
	gClient grpcPb.LogServiceClient
	logger  *logger.Logger
}

func NewClient(logger *logger.Logger) *GrpcClient {
	var conn *grpc.ClientConn
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", viper.GetString("grpc.ip"), viper.GetString("grpc.port")),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(fmt.Sprintf("Error connect to gRPC: %s", err))
		return nil
	}
	//defer conn.Close()

	return &GrpcClient{
		conn:    conn,
		gClient: grpcPb.NewLogServiceClient(conn),
		logger:  logger,
	}
}
func (g *GrpcClient) ConnClose() error {
	return g.conn.Close()
}
func (g *GrpcClient) SendLogRequest(ctx context.Context, logRequest grpcPb.LogItem) error {
	entity, err := grpcPb.ToPbEntity(logRequest.Entity)
	if err != nil {
		g.logger.Error(fmt.Sprintf("Error parse entity to PB: %s", err))
		return err
	}
	action, err := grpcPb.ToPbAction(logRequest.Action)
	if err != nil {
		g.logger.Error(fmt.Sprintf("Error parse action to PB: %s", err))
		return err
	}
	_, err = g.gClient.Log(ctx, &grpcPb.LogRequest{
		Entity:    entity,
		Action:    action,
		EntityId:  &logRequest.EntityID,
		UserId:    logRequest.UserID,
		Timestamp: timestamppb.New(logRequest.Timestamp),
	})
	if err != nil {
		g.logger.Error(fmt.Sprintf("LogRequest failed: %s", err))
		return err
	}
	return nil
}

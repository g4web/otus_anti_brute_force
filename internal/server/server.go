package server

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"

	"github.com/g4web/otus_anti_brute_force/internal"
	"github.com/g4web/otus_anti_brute_force/internal/config"
	"github.com/g4web/otus_anti_brute_force/internal/proto"
	"google.golang.org/grpc"
)

type ABFServer struct {
	app        *app.App
	grpcServer *grpc.Server
	config     *config.Config
	proto.UnimplementedAntiBruteForceServer
}

func NewABFServer(app *app.App, config *config.Config) *ABFServer {
	grpcServer := grpc.NewServer()
	ABFServer := &ABFServer{
		app:        app,
		grpcServer: grpcServer,
		config:     config,
	}
	proto.RegisterAntiBruteForceServer(grpcServer, ABFServer)

	return ABFServer
}

func (a *ABFServer) Start(ctx context.Context) error {
	lsn, err := net.Listen("tcp", net.JoinHostPort(a.config.GrpcHost, a.config.GrpcPort))
	if err != nil {
		log.Fatalf("Fail start gprc server: %v", err)
		return err
	}

	if err := a.grpcServer.Serve(lsn); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Fail start gprc server: %v", err)
		return err
	}

	<-ctx.Done()

	return nil
}

func (a *ABFServer) Stop(ctx context.Context) error {
	a.grpcServer.GracefulStop()

	return nil
}

func (a *ABFServer) IsOk(ctx context.Context, request *proto.UserRequest) (*proto.UserResponse, error) {
	isOk, err := a.app.IsOk(request.GetIP(), request.GetLogin(), request.GetPassword())
	if err != nil {
		return nil, err
	}

	return &proto.UserResponse{IsOk: isOk}, nil
}

func (a *ABFServer) DeleteLoginStats(
	ctx context.Context,
	request *proto.DeleteLoginStatsRequest,
) (*proto.BaseResponse, error) {
	a.app.DeleteLoginStats(request.GetLogin())

	return &proto.BaseResponse{IsSuccess: true}, nil
}

func (a *ABFServer) DeleteIPStats(
	ctx context.Context,
	request *proto.DeleteIPStatsRequest,
) (*proto.BaseResponse, error) {
	a.app.DeleteIPStats(request.GetIP())

	return &proto.BaseResponse{IsSuccess: true}, nil
}

func (a *ABFServer) AddNetworkToWhiteList(
	ctx context.Context,
	request *proto.AddNetworkToWhiteListRequest,
) (*proto.BaseResponse, error) {
	err := a.app.AddNetworkToWhiteList(request.GetNetwork())
	if err != nil {
		return nil, err
	}

	return &proto.BaseResponse{IsSuccess: true}, nil
}

func (a *ABFServer) AddNetworkToBlackList(
	ctx context.Context,
	request *proto.AddNetworkToBlackListRequest,
) (*proto.BaseResponse, error) {
	err := a.app.AddNetworkToBlackList(request.GetNetwork())
	if err != nil {
		return nil, err
	}

	return &proto.BaseResponse{IsSuccess: true}, nil
}

func (a *ABFServer) RemoveNetworkFromWhiteList(
	ctx context.Context,
	request *proto.RemoveNetworkFromWhiteListRequest,
) (*proto.BaseResponse, error) {
	err := a.app.RemoveNetworkFromWhiteList(request.GetNetwork())
	if err != nil {
		return nil, err
	}

	return &proto.BaseResponse{IsSuccess: true}, nil
}

func (a *ABFServer) RemoveNetworkFromBlackList(
	ctx context.Context,
	request *proto.RemoveNetworkFromBlackListRequest,
) (*proto.BaseResponse, error) {
	err := a.app.RemoveNetworkFromBlackList(request.GetNetwork())
	if err != nil {
		return nil, err
	}

	return &proto.BaseResponse{IsSuccess: true}, nil
}

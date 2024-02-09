package notification

import (
	"context"
	pb "github.com/vaberof/notification-grpc/genproto/notification_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type serverAPI struct {
	pb.UnimplementedNotificationServiceServer
	notificationService NotificationService
}

func Register(gRPC *grpc.Server, notificationService NotificationService) {
	pb.RegisterNotificationServiceServer(gRPC, &serverAPI{notificationService: notificationService})
}

func (s *serverAPI) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*emptypb.Empty, error) {
	err := s.notificationService.SendEmail(ctx, req.To, req.Type, req.Subject, req.Body)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}

	return &emptypb.Empty{}, nil
}

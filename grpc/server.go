package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/sjengpho/tin/mail/gmail"
	"github.com/sjengpho/tin/os/network"
	"github.com/sjengpho/tin/os/packagemanager"
	"github.com/sjengpho/tin/os/temperature"
	"github.com/sjengpho/tin/proto/pb"
	"github.com/sjengpho/tin/tin"

	"google.golang.org/grpc"
)

// Server represents the GRPC server.
type Server struct {
	config                *tin.Config
	packageManagerService *tin.PackageManagerService
	temperatureService    *tin.TemperatureService
	networkService        *tin.NetworkService
	mailService           *tin.MailService
	gmail                 *gmail.Service
}

// logger returns a log.Logger with the given prefix.
func logger(prefix string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("%v ", prefix), log.Ldate|log.Ltime|log.Lshortfile)
}

// NewServer creates workers and initializes and returns a grpc.Server.
func NewServer(c tin.Config) Server {
	server := Server{
		config:                &c,
		gmail:                 gmail.NewService(c.GmailCredentials, c.GmailToken),
		mailService:           tin.NewMailService(gmail.NewService(c.GmailCredentials, c.GmailToken), logger("MailService")),
		networkService:        tin.NewNetworkService(network.NewNameLookup(), network.NewPublicIPLookup(), logger("NetworkService")),
		packageManagerService: tin.NewPackageManagerService(packagemanager.New(), logger("PackageManagerService")),
		temperatureService:    tin.NewTemperatureService(temperature.NewReader(), logger("TemperatureService")),
	}

	return server
}

// ListenAndServe starts the server.
func (s *Server) ListenAndServe(port int) error {
	address := fmt.Sprintf("0.0.0.0:%v", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTinServiceServer(grpcServer, s)

	log.Printf("Listening on %v", address)

	return grpcServer.Serve(listener)
}

// GmailAuthURL returns a pb.GmailAuthURLResponse.
func (s *Server) GmailAuthURL(c context.Context, r *pb.GmailAuthURLRequest) (*pb.GmailAuthURLResponse, error) {
	authURL, err := s.gmail.AuthURL()
	if err != nil {
		return nil, err
	}

	return &pb.GmailAuthURLResponse{AuthURL: authURL}, nil
}

// GmailAuthCode returns a pb.GmailAuthCodeResponse.
func (s *Server) GmailAuthCode(c context.Context, r *pb.GmailAuthCodeRequest) (*pb.GmailAuthCodeResponse, error) {
	err := s.gmail.ExchangeAuthCode(r.GetAuthCode())
	if err != nil {
		return nil, err
	}

	return &pb.GmailAuthCodeResponse{}, nil
}

// GmailUnread returns a pb.GmailUnreadResponse.
func (s *Server) GmailUnread(c context.Context, r *pb.GmailUnreadRequest) (*pb.GmailUnreadResponse, error) {
	m := s.mailService.UnreadMailCount()
	return &pb.GmailUnreadResponse{Value: int32(m)}, nil
}

// AvailableUpdates returns a pb.AvailableUpdatesResponse.
func (s *Server) AvailableUpdates(c context.Context, r *pb.AvailableUpdatesRequest) (*pb.AvailableUpdatesResponse, error) {
	u := s.packageManagerService.AvailableUpdatesCount()
	return &pb.AvailableUpdatesResponse{Value: int32(u)}, nil
}

// Temperature returns a pb.TemperatureResponse.
func (s *Server) Temperature(c context.Context, r *pb.TemperatureRequest) (*pb.TemperatureResponse, error) {
	t := s.temperatureService.Temperature()
	temperature := &pb.Temperature{
		Celsius:    int32(t.Celsius()),
		Fahrenheit: int32(t.Fahrenheit()),
}
	return &pb.TemperatureResponse{Temperature: temperature}, nil
}

// ESSID returns a pb.NetworkNameResponse.
func (s *Server) ESSID(c context.Context, r *pb.ESSIDRequest) (*pb.ESSIDResponse, error) {
	n := s.networkService.Name()
	return &pb.ESSIDResponse{Value: string(n)}, nil
}

// IPAddress returns a pb.IPAddressResponse.
func (s *Server) IPAddress(c context.Context, r *pb.IPAddressRequest) (*pb.IPAddressResponse, error) {
	v := s.networkService.IP()
	return &pb.IPAddressResponse{Value: v}, nil
}

// Config returns a pb.Config.
func (s *Server) Config(c context.Context, r *pb.ConfigRequest) (*pb.ConfigResponse, error) {
	resp := &pb.ConfigResponse{
		Config: &pb.Config{
		GmailCredentials: s.config.GmailCredentials,
		GmailToken:       s.config.GmailToken,
		},
	}
	return resp, nil
}

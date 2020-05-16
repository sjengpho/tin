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
	"github.com/sjengpho/tin/proto"
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
	proto.RegisterTinServer(grpcServer, s)

	log.Printf("Listening on %v", address)

	return grpcServer.Serve(listener)
}

// GmailAuthURL returns a proto.GmailAuthURLResponse.
func (s *Server) GmailAuthURL(c context.Context, r *proto.GmailAuthURLRequest) (*proto.GmailAuthURLResponse, error) {
	authURL, err := s.gmail.AuthURL()
	if err != nil {
		return nil, err
	}

	return &proto.GmailAuthURLResponse{AuthURL: authURL}, nil
}

// GmailAuthCode returns a proto.GmailAuthCodeResponse.
func (s *Server) GmailAuthCode(c context.Context, r *proto.GmailAuthCodeRequest) (*proto.GmailAuthCodeResponse, error) {
	err := s.gmail.ExchangeAuthCode(r.GetAuthCode())
	if err != nil {
		return nil, err
	}

	return &proto.GmailAuthCodeResponse{}, nil
}

// GmailUnread returns a proto.GmailUnreadResponse.
func (s *Server) GmailUnread(c context.Context, r *proto.GmailUnreadRequest) (*proto.GmailUnreadResponse, error) {
	m := s.mailService.UnreadMailCount()
	return &proto.GmailUnreadResponse{Value: int32(m)}, nil
}

// AvailableUpdates returns a proto.AvailableUpdatesResponse.
func (s *Server) AvailableUpdates(c context.Context, r *proto.AvailableUpdatesRequest) (*proto.AvailableUpdatesResponse, error) {
	u := s.packageManagerService.AvailableUpdatesCount()
	return &proto.AvailableUpdatesResponse{Value: int32(u)}, nil
}

// TemperatureCelsius returns a proto.TemperatureResponse.
func (s *Server) TemperatureCelsius(c context.Context, r *proto.TemperatureRequest) (*proto.TemperatureResponse, error) {
	t := s.temperatureService.Temperature()
	return &proto.TemperatureResponse{Value: int32(t.Celsius())}, nil
}

// TemperatureFahrenheit returns a proto.TemperatureResponse.
func (s *Server) TemperatureFahrenheit(c context.Context, r *proto.TemperatureRequest) (*proto.TemperatureResponse, error) {
	t := s.temperatureService.Temperature()
	return &proto.TemperatureResponse{Value: int32(t.Fahrenheit())}, nil
}

// ESSID returns a proto.NetworkNameResponse.
func (s *Server) ESSID(c context.Context, r *proto.ESSIDRequest) (*proto.ESSIDResponse, error) {
	n := s.networkService.Name()
	return &proto.ESSIDResponse{Value: string(n)}, nil
}

// IPAddress returns a proto.IPAddressResponse.
func (s *Server) IPAddress(c context.Context, r *proto.IPAddressRequest) (*proto.IPAddressResponse, error) {
	v := s.networkService.IP()
	return &proto.IPAddressResponse{Value: v}, nil
}

// Config returns a proto.Config.
func (s *Server) Config(c context.Context, r *proto.ConfigRequest) (*proto.ConfigResponse, error) {
	resp := &proto.ConfigResponse{
		GmailCredentials: s.config.GmailCredentials,
		GmailToken:       s.config.GmailToken,
	}
	return resp, nil
}

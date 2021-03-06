package grpc

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/sjengpho/tin/proto/pb"
	"google.golang.org/grpc"
)

// Client represents the GRPC client.
type Client struct {
	client pb.TinServiceClient
	conn   *grpc.ClientConn
}

// NewClient attempts to create a connection and returns a grpc.Client.
func NewClient(address string) (*Client, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failled connecting to %v: %w", address, err)
	}

	return &Client{conn: conn, client: pb.NewTinServiceClient(conn)}, nil
}

// AvailableUpdates returns a integer.
func (c *Client) AvailableUpdates() (int, error) {
	resp, err := c.client.AvailableUpdates(context.Background(), &pb.AvailableUpdatesRequest{})
	if err != nil {
		return 0, err
	}

	return int(resp.GetValue()), nil
}

// Temperature returns a pb.TemperatureResponse.
func (c *Client) Temperature() (*pb.TemperatureResponse, error) {
	resp, err := c.client.Temperature(context.Background(), &pb.TemperatureRequest{})
	if err != nil {
		return &pb.TemperatureResponse{}, err
	}

	return resp, nil
}

// GmailUnread returns a integer.
func (c *Client) GmailUnread() (int, error) {
	response, err := c.client.GmailUnread(context.Background(), &pb.GmailUnreadRequest{})
	if err != nil {
		return 0, err
	}

	return int(response.GetValue()), nil
}

// GmailAuthURL returns a string.
func (c *Client) GmailAuthURL() (string, error) {
	response, err := c.client.GmailAuthURL(context.Background(), &pb.GmailAuthURLRequest{})
	if err != nil {
		return "", err
	}

	return response.GetAuthURL(), nil
}

// GmailAuthCode returns a boolean.
func (c *Client) GmailAuthCode(code string) bool {
	request := &pb.GmailAuthCodeRequest{AuthCode: code}
	_, err := c.client.GmailAuthCode(context.Background(), request)

	return err == nil
}

// ESSID returns a string.
func (c *Client) ESSID() (string, error) {
	resp, err := c.client.ESSID(context.Background(), &pb.ESSIDRequest{})
	if err != nil {
		return "", err
	}

	return resp.Value, nil
}

// IPAddress returns a string.
func (c *Client) IPAddress() (string, error) {
	resp, err := c.client.IPAddress(context.Background(), &pb.IPAddressRequest{})
	if err != nil {
		return "", err
	}

	return resp.Value, nil
}

// Config returns a pb.Config.
func (c *Client) Config() (*pb.ConfigResponse, error) {
	resp, err := c.client.Config(context.Background(), &pb.ConfigRequest{})
	if err != nil {
		return &pb.ConfigResponse{}, err
	}

	return resp, nil
}

// InstalledPackages returns a pb.InstalledPackagesResponse.
func (c *Client) InstalledPackages() (*pb.InstalledPackagesResponse, error) {
	resp, err := c.client.InstalledPackages(context.Background(), &pb.InstalledPackagesRequest{})
	if err != nil {
		return &pb.InstalledPackagesResponse{}, err
	}

	return resp, nil
}

// InstalledPackagesSubscribe executes the process function when it receives a message.
func (c *Client) InstalledPackagesSubscribe(process func(r *pb.InstalledPackagesResponse)) error {
	stream, err := c.client.InstalledPackagesSubscribe(context.Background(), &pb.InstalledPackagesRequest{})
	if err != nil {
		return err
	}
	for {
		t, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		process(t)
	}
	return nil
}

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

// TemperatureCelsius returns a integer.
func (c *Client) TemperatureCelsius() (int, error) {
	resp, err := c.client.TemperatureCelsius(context.Background(), &pb.TemperatureRequest{})
	if err != nil {
		return 0, err
	}

	return int(resp.GetValue()), nil
}

// TemperatureFahrenheit returns a integer.
func (c *Client) TemperatureFahrenheit() (int, error) {
	resp, err := c.client.TemperatureFahrenheit(context.Background(), &pb.TemperatureRequest{})
	if err != nil {
		return 0, err
	}

	return int(resp.GetValue()), nil
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

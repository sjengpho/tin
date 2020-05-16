package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/sjengpho/tin/proto"
	"google.golang.org/grpc"
)

// Client represents the GRPC client.
type Client struct {
	client proto.TinClient
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

	return &Client{conn: conn, client: proto.NewTinClient(conn)}, nil
}

// AvailableUpdates returns a integer.
func (c *Client) AvailableUpdates() (int, error) {
	resp, err := c.client.AvailableUpdates(context.Background(), &proto.AvailableUpdatesRequest{})
	if err != nil {
		return 0, err
	}

	return int(resp.GetValue()), nil
}

// TemperatureCelsius returns a integer.
func (c *Client) TemperatureCelsius() (int, error) {
	resp, err := c.client.TemperatureCelsius(context.Background(), &proto.TemperatureRequest{})
	if err != nil {
		return 0, err
	}

	return int(resp.GetValue()), nil
}

// TemperatureFahrenheit returns a integer.
func (c *Client) TemperatureFahrenheit() (int, error) {
	resp, err := c.client.TemperatureFahrenheit(context.Background(), &proto.TemperatureRequest{})
	if err != nil {
		return 0, err
	}

	return int(resp.GetValue()), nil
}

// GmailUnread returns a integer.
func (c *Client) GmailUnread() (int, error) {
	response, err := c.client.GmailUnread(context.Background(), &proto.GmailUnreadRequest{})
	if err != nil {
		return 0, err
	}

	return int(response.GetValue()), nil
}

// GmailAuthURL returns a string.
func (c *Client) GmailAuthURL() (string, error) {
	response, err := c.client.GmailAuthURL(context.Background(), &proto.GmailAuthURLRequest{})
	if err != nil {
		return "", err
	}

	return response.GetAuthURL(), nil
}

// GmailAuthCode returns a boolean.
func (c *Client) GmailAuthCode(code string) bool {
	request := &proto.GmailAuthCodeRequest{AuthCode: code}
	_, err := c.client.GmailAuthCode(context.Background(), request)

	return err == nil
}

// ESSID returns a string.
func (c *Client) ESSID() (string, error) {
	resp, err := c.client.ESSID(context.Background(), &proto.ESSIDRequest{})
	if err != nil {
		return "", err
	}

	return resp.Value, nil
}

// IPAddress returns a string.
func (c *Client) IPAddress() (string, error) {
	resp, err := c.client.IPAddress(context.Background(), &proto.IPAddressRequest{})
	if err != nil {
		return "", err
	}

	return resp.Value, nil
}

// Config returns a proto.Config.
func (c *Client) Config() (*proto.ConfigResponse, error) {
	resp, err := c.client.Config(context.Background(), &proto.ConfigRequest{})
	if err != nil {
		return &proto.ConfigResponse{}, err
	}

	return resp, nil
}

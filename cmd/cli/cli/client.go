package cli

import (
	"fmt"
	"log"

	"github.com/sjengpho/tin/grpc"
)

// NewClient returns a grpc.Client.
func NewClient(port int) *grpc.Client {
	client, err := grpc.NewClient(fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal(err)
	}

	return client
}

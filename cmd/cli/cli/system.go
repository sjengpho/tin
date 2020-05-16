package cli

import (
	"fmt"
	"log"

	"github.com/sjengpho/tin/grpc"
)

// NewSystemCommander returns a cli.SystemCommander.
func NewSystemCommander() SystemCommander {
	return &systemCommander{}
}

// systemCommander implements cli.SystemCommander.
type systemCommander struct{}

// SystemUpdates outputs the available update count.
func (s *systemCommander) SystemUpdates(c *grpc.Client) {
	v, err := c.AvailableUpdates()
	if err != nil {
		log.Printf("failed getting the available updates: %v", err)
		return
	}
	fmt.Println(v)
}

// SystemTemperatureCelsius outputs the temperature in celsius format.
func (s *systemCommander) SystemTemperatureCelsius(c *grpc.Client) {
	v, err := c.TemperatureCelsius()
	if err != nil {
		log.Printf("failed getting the temperature: %v", err)
		return
	}
	fmt.Println(v)
}

// SystemTemperatureFahrenheit outputs the temperature in celsius format.
func (s *systemCommander) SystemTemperatureFahrenheit(c *grpc.Client) {
	v, err := c.TemperatureFahrenheit()
	if err != nil {
		log.Printf("failed getting the temperature: %v", err)
		return
	}
	fmt.Println(v)
}

package cli

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/sjengpho/tin/grpc"
	"github.com/sjengpho/tin/proto/pb"
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
	v, err := c.Temperature()
	if err != nil {
		log.Printf("failed getting the temperature: %v", err)
		return
	}
	fmt.Println(v.GetTemperature().GetCelsius())
}

// SystemTemperatureFahrenheit outputs the temperature in celsius format.
func (s *systemCommander) SystemTemperatureFahrenheit(c *grpc.Client) {
	v, err := c.Temperature()
	if err != nil {
		log.Printf("failed getting the temperature: %v", err)
		return
	}
	fmt.Println(v.GetTemperature().GetFahrenheit())
}

// SystemInstalled outputs or exports the installed packages.
func (s *systemCommander) SystemInstalled(c *grpc.Client, flags SystemInstalledFlags) {
	if flags.Subscribe && flags.Export {
		c.InstalledPackagesSubscribe(func(r *pb.InstalledPackagesResponse) {
			s.exportPackages(flags.ExportPath, r)
		})
	}

	if !flags.Subscribe && flags.Export {
		r, err := c.InstalledPackages()
		if err != nil {
			log.Printf("failed getting installed packages: %v", err)
			return
		}
		s.exportPackages(flags.ExportPath, r)
		return
	}

	if flags.Subscribe && !flags.Export {
		c.InstalledPackagesSubscribe(s.outputPackages)
	}

	r, err := c.InstalledPackages()
	if err != nil {
		log.Printf("failed getting installed packages: %v", err)
		return
	}
	s.outputPackages(r)
}

// outputPackages prints the packages to standard output.
func (s *systemCommander) outputPackages(r *pb.InstalledPackagesResponse) {
	for _, p := range r.Packages {
		fmt.Printf("%v %v\n", p.GetName(), p.GetVersion())
	}
}

// exportPackages exports the packages into a CSV file.
func (s *systemCommander) exportPackages(path string, r *pb.InstalledPackagesResponse) {
	name := path
	if name == "" {
		name = "installed_packages.csv"
	}

	file, err := os.Create(name)
	if err != nil {
		log.Printf("failed creating file: %v", err)
		return
	}

	writer := csv.NewWriter(file)
	if err := writer.Write([]string{"Name", "Version"}); err != nil {
		log.Printf("failed writing to file: %v", err)
		return
	}
	for _, p := range r.GetPackages() {
		if err := writer.Write([]string{p.Name, p.Version}); err != nil {
			log.Printf("failed writing to file: %v", err)
			return
		}
	}
	writer.Flush()
}

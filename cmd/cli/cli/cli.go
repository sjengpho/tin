package cli

import "github.com/sjengpho/tin/grpc"

// SystemCommander is the interace implemented by an object that can
// output system related info.
//
// SystemUpdates outputs the available update count.
// SystemInstalled outputs the installed packages.
// SystemTemperatureCelsius outputs the temperature in celsius format.
// SystemTemperatureFahrenheit outputs the temperature in fahrenheit format.
type SystemCommander interface {
	SystemUpdates(c *grpc.Client)
	SystemInstalled(c *grpc.Client, flags SystemInstalledFlags)
	SystemTemperatureCelsius(c *grpc.Client)
	SystemTemperatureFahrenheit(c *grpc.Client)
}

// SystemInstalledFlags represents the flags.
type SystemInstalledFlags struct {
	Subscribe  bool
	Export     bool
	ExportPath string
}

// NetworkCommander is the interace implemented by an object that can
// output network related info.
//
// ESSID outputs the network name.
// IP outputs the public IP address.
type NetworkCommander interface {
	ESSID(c *grpc.Client)
	IP(c *grpc.Client)
}

// GmailCommander is the interace implemented by an object that can
// output gmail related info.
//
// Login attempts to authorize the user at Gmail.
// Unread outputs the unread mail count.
type GmailCommander interface {
	Login(c *grpc.Client)
	Unread(c *grpc.Client)
}

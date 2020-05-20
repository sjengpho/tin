package main

import (
	"github.com/sjengpho/tin/cmd/cli/cli"
	"github.com/spf13/cobra"
)

// config holds general configuration values.
type config struct {
	port int
}

func main() {
	config := config{}

	c := &cobra.Command{Use: "tin"}
	c.PersistentFlags().IntVar(&config.port, "port", 8717, "The server port")
	c.AddCommand(NewCmdSystem(cli.NewSystemCommander(), &config))
	c.AddCommand(NewCmdNetwork(cli.NewNetworkCommander(), &config))
	c.AddCommand(NewCmdGmail(cli.NewGmailCommander(), &config))
	c.SetHelpCommand((&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	}))
	c.Execute()
}

// NewCmdSystem returns a cobra.Command.
func NewCmdSystem(s cli.SystemCommander, c *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "System info",
		Long:  `System info`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "updates",
		Short: "Available updates",
		Long:  `Available updates`,
		Run: func(cmd *cobra.Command, args []string) {
			s.SystemUpdates(cli.NewClient(c.port))
		},
	})

	systemInstalledFlags := cli.SystemInstalledFlags{}
	installedPackagesCmd := &cobra.Command{
		Use:   "installed",
		Short: "Installed packages",
		Long:  `Installed packages`,
		Run: func(cmd *cobra.Command, args []string) {

			s.SystemInstalled(cli.NewClient(c.port), systemInstalledFlags)
		},
	}
	installedPackagesCmd.PersistentFlags().BoolVar(&systemInstalledFlags.Subscribe, "subscribe", false, "Automatically process changes")
	installedPackagesCmd.PersistentFlags().BoolVar(&systemInstalledFlags.Export, "export", false, "Creates a CSV export")
	installedPackagesCmd.PersistentFlags().StringVar(&systemInstalledFlags.ExportPath, "exportPath", "", "CSV export path")
	cmd.AddCommand(installedPackagesCmd)

	cmd.AddCommand(&cobra.Command{
		Use:   "celsius",
		Short: "Temperature celsius",
		Long:  `Temperature celsius`,
		Run: func(cmd *cobra.Command, args []string) {
			s.SystemTemperatureCelsius(cli.NewClient(c.port))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "fahrenheit",
		Short: "Temperature fahrenheit",
		Long:  `Temperature fahrenheit`,
		Run: func(cmd *cobra.Command, args []string) {
			s.SystemTemperatureFahrenheit(cli.NewClient(c.port))
		},
	})

	return cmd
}

// NewCmdNetwork returns a cobra.Command.
func NewCmdNetwork(s cli.NetworkCommander, c *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Network info",
		Long:  `Network info`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "essid",
		Short: "Network name",
		Long:  `Network name`,
		Run: func(cmd *cobra.Command, args []string) {
			s.ESSID(cli.NewClient(c.port))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "ip",
		Short: "IP address",
		Long:  `IP address`,
		Run: func(cmd *cobra.Command, args []string) {
			s.IP(cli.NewClient(c.port))
		},
	})

	return cmd
}

// NewCmdGmail returns a cobra.Command.
func NewCmdGmail(s cli.GmailCommander, c *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gmail",
		Short: "Gmail info",
		Long:  `Gmail info`,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "login",
		Short: "Gmail authorization",
		Long:  `Gmail authorization`,
		Run: func(cmd *cobra.Command, args []string) {
			s.Login(cli.NewClient(c.port))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "unread",
		Short: "Unread mail count",
		Long:  `Unread mail count`,
		Run: func(cmd *cobra.Command, args []string) {
			s.Unread(cli.NewClient(c.port))
		},
	})

	return cmd
}

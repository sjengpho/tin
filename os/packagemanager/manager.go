package packagemanager

import (
	"os/exec"
	"strings"

	"github.com/sjengpho/tin/tin"
)

var execCommand = exec.Command
var lookPath = exec.LookPath

// New returns a tin.PackageManager.
//
// If a manager couldn't be resolved it will return nil.
func New() tin.PackageManager {
	if _, err := lookPath("xbps-install"); err == nil {
		return &XBPS{}
	}

	return nil
}

// XBPS implements tin.PackageManager.
type XBPS struct{}

// AvailableUpdates returns a slice of tin.Package.
//
// It executes xbps-install in dry-run mode by using -Mun as argument.
func (x *XBPS) AvailableUpdates() ([]tin.Package, error) {
	output, err := execCommand("xbps-install", "-Mun").Output()
	if err != nil {
		return []tin.Package{}, err
	}

	return x.parse(string(output)), nil
}

// Installed returns a slice of tin.Package.
//
// It executes xbps-query -m.
func (x *XBPS) Installed() ([]tin.Package, error) {
	output, err := execCommand("xbps-query", "-m").Output()
	if err != nil {
		return []tin.Package{}, err
	}

	return x.parse(string(output)), nil
}

// parse parses the string into a slice of tin.Package.
//
// It assumes that the output contains a multiline string of packages,
// separated by newlines. Blank lines are ignored.
// Example of a line: package-name-3.5.2_1
func (x *XBPS) parse(output string) []tin.Package {
	pp := []tin.Package{}
	for _, v := range strings.Split(output, "\n") {
		if v == "" {
			continue
		}

		p := strings.Split(v, " ")[0]  // Removing everything after the first white space.
		i := strings.LastIndex(p, "-") // Getting the index of the separator between the package name and version.
		pp = append(pp, tin.Package{
			Name:    p[:i],   // Extracting everything from the begin until the index of the separator.
			Version: p[i+1:], // Extracting everything after the index of the separator until the end.
		})
	}
	return pp
}

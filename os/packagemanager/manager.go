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

// AvailableUpdates executes xbps-install in dry-run mode by
// using -Mun as argument.
//
// It assumes that the output contains a multiline string of packages,
// separated by newlines. Blank lines are ignored.
func (x *XBPS) AvailableUpdates() ([]tin.Package, error) {
	output, err := execCommand("xbps-install", "-Mun").Output()
	if err != nil {
		return []tin.Package{}, err
	}

	pp := []tin.Package{}
	for _, v := range strings.Split(string(output), "\n") {
		if v == "" {
			continue
		}

		// Example value of v: package-name-3.5.2_1 update x86_64 https://alpha.de.repo.voidlinux.org/current 182572180 59477905
		s := strings.Split(v, " ")[0]  // Example: package-name-3.5.2_1
		i := strings.LastIndex(s, "-") // Example: 12
		pp = append(pp, tin.Package{
			Name:    s[:i],   // Example: package-name
			Version: s[i+1:], // Example: 3.5.2_1
		})
	}

	return pp, nil
}

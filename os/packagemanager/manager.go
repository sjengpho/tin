package packagemanager

import (
	"errors"
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

	if _, err := lookPath("checkupdates"); err == nil {
		return &Arch{
			Pacman: Pacman{},
			AUR:    &Yay{},
		}
	}

	return nil
}

// XBPS implements tin.PackageManager.
type XBPS struct{}

// AvailableUpdates returns a slice of tin.Package.
func (x *XBPS) AvailableUpdates() ([]tin.Package, error) {
	output, err := execCommand("xbps-install", "-Mun").Output()
	if err != nil {
		return []tin.Package{}, err
	}

	return x.parse(string(output)), nil
}

// Installed returns a slice of tin.Package.
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

// Arch implements tin.PackageManager.
type Arch struct {
	Pacman Pacman
	AUR    tin.PackageManager
}

// AvailableUpdates returns a slice of tin.Package.
func (a *Arch) AvailableUpdates() ([]tin.Package, error) {
	pacmanPackages, err := a.Pacman.AvailableUpdates()
	if err != nil {
		return []tin.Package{}, err
	}

	aurPackages, err := a.AUR.AvailableUpdates()
	if err != nil {
		return []tin.Package{}, err
	}

	return append(pacmanPackages, aurPackages...), nil
}

// Installed returns a slice of tin.Package.
func (a *Arch) Installed() ([]tin.Package, error) {
	p := Pacman{}
	packages, err := p.Installed()
	if err != nil {
		return []tin.Package{}, err
	}

	return packages, nil
}

// Pacman implements tin.PackageManager.
type Pacman struct{}

// AvailableUpdates returns a slice of tin.Package.
func (p *Pacman) AvailableUpdates() ([]tin.Package, error) {
	output, err := execCommand("checkupdates").Output()
	if err != nil {
		var e *exec.ExitError
		// Assuming exit code 2 means no updates.
		if errors.As(err, &e) && e.ExitCode() == 2 {
			return []tin.Package{}, nil
		}
		return []tin.Package{}, err
	}

	return p.parse(string(output)), nil
}

// Installed returns a slice of tin.Package.
func (p *Pacman) Installed() ([]tin.Package, error) {
	output, err := execCommand("pacman", "-Qe").Output()
	if err != nil {
		return []tin.Package{}, err
	}

	return p.parse(string(output)), nil
}

// parse parses the string into a slice of tin.Package.
//
// It assumes that the output contains a multiline string of packages,
// separated by newlines. Blank lines are ignored.
// Example of a line: package-name 1.2.0-1
func (p *Pacman) parse(output string) []tin.Package {
	pp := []tin.Package{}
	for _, v := range strings.Split(output, "\n") {
		if v == "" {
			continue
		}

		p := strings.Split(v, " ")
		pp = append(pp, tin.Package{
			Name:    p[0],
			Version: p[1],
		})
	}
	return pp
}

// Yay implements tin.PackageManager.
type Yay struct{}

// AvailableUpdates returns a slice of tin.Package.
func (y *Yay) AvailableUpdates() ([]tin.Package, error) {
	output, err := execCommand("yay", "-Qum").Output()
	if err != nil {
		return []tin.Package{}, err
	}

	return y.parse(string(output)), nil
}

// Installed returns a slice of tin.Package.
func (y *Yay) Installed() ([]tin.Package, error) {
	return []tin.Package{}, errors.New("unimplemented")
}

// parse parses the string into a slice of tin.Package.
//
// It assumes that the output contains a multiline string of packages,
// separated by newlines. Blank lines are ignored.
// Example of a line: package-name 1.2.0-1
func (y *Yay) parse(output string) []tin.Package {
	pp := []tin.Package{}
	for _, v := range strings.Split(output, "\n") {
		if v == "" {
			continue
		}

		p := strings.Split(v, " ")
		pp = append(pp, tin.Package{
			Name:    p[0],
			Version: p[1],
		})
	}
	return pp
}

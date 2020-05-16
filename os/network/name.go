package network

import (
	"os/exec"
	"strings"

	"github.com/sjengpho/tin/tin"
)

var execCommand = exec.Command
var lookPath = exec.LookPath

// NewNameLookup returns a tin.ESSIDLookup.
//
// If tin.ESSIDLookup couldn't be resolved it will return nil.
func NewNameLookup() tin.ESSIDLookup {
	if _, err := lookPath("iwgetid"); err == nil {
		return &iwgetid{}
	}

	return nil
}

// iwgetid implements tin.ESSIDLookup.
type iwgetid struct{}

// Lookup returns a tin.ESSID.
//
// It uses iwgetid to fetch the network name (ESSID).
func (i *iwgetid) Lookup() (tin.ESSID, error) {
	output, err := execCommand("iwgetid", "-r").Output()
	if err != nil {
		return "", err
	}

	return tin.ESSID(strings.TrimSpace(string(output))), nil
}

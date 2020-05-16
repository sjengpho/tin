package temperature

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sjengpho/tin/tin"
)

var execCommand = exec.Command
var osStat = os.Stat
var readFile = ioutil.ReadFile

// NewReader returns a tin.TemperatureReader.
//
// If a supported reader couldn't be resolved it will return nil.
func NewReader() tin.TemperatureReader {
	p := "/sys/class/thermal/thermal_zone2/temp"
	if _, err := osStat(p); err == nil {
		return &FileReader{p}
	}

	return nil
}

// FileReader implements tin.TemperatureReader.
type FileReader struct {
	path string
}

// Read reads the temperature from the file.
//
// It assumes that the file contains a millidegree temperature in Celsius format.
func (f *FileReader) Read() (tin.Temperature, error) {
	bytes, err := readFile(f.path)
	if err != nil {
		return tin.Temperature{}, err
	}

	v, err := strconv.Atoi(strings.Trim(string(bytes), "\n"))
	if err != nil {
		return tin.Temperature{}, err
	}

	return tin.Temperature{Value: v / 1000}, nil
}

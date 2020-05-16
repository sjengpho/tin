package network

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/sjengpho/tin/tin"
)

func fakeLookPathSuccess(file string) (string, error) {
	return "fake-path", nil
}

func fakeLookPathError(file string) (string, error) {
	return "", errors.New("executeable doesn't exists")
}

func fakeExecCommand(commandName string) func(name string, args ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		cs := []string{fmt.Sprintf("-test.run=%v", commandName), "--", name}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{"GO_TEST_PROCESS=1"}
		return cmd
	}
}

func TestNewNameLookupSuccess(t *testing.T) {
	lookPath = fakeLookPathSuccess
	defer func() { lookPath = exec.LookPath }()

	want := reflect.TypeOf(&iwgetid{})
	got := reflect.TypeOf(NewNameLookup())
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestNewNameLookupError(t *testing.T) {
	lookPath = fakeLookPathError
	defer func() { lookPath = exec.LookPath }()

	want := reflect.TypeOf(nil)
	got := reflect.TypeOf(NewNameLookup())
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestNameLookupSuccess(t *testing.T) {
	execCommand = fakeExecCommand("TestNameLookupCommandSuccess")
	defer func() { execCommand = exec.Command }()
	lookupper := iwgetid{}

	want := tin.ESSID("ESSID")
	got, _ := lookupper.Lookup()
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestNameLookupCommandSuccess(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}

	fmt.Println("ESSID")
	os.Exit(0)
}

func TestNameLookupError(t *testing.T) {
	execCommand = fakeExecCommand("TestNameLookupCommandError")
	defer func() { execCommand = exec.Command }()
	lookupper := iwgetid{}

	_, got := lookupper.Lookup()
	if got == nil {
		t.Errorf("want %v, got %v", "exit status 1", got)
	}
}

func TestNameLookupCommandError(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}

	os.Exit(1)
}

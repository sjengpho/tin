package packagemanager

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"
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

func TestNewSuccess(t *testing.T) {
	lookPath = fakeLookPathSuccess
	defer func() { lookPath = exec.LookPath }()

	want := reflect.TypeOf(&XBPS{})
	got := reflect.TypeOf(New())
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestNewError(t *testing.T) {
	lookPath = fakeLookPathError
	defer func() { lookPath = exec.LookPath }()

	want := reflect.TypeOf(nil)
	got := reflect.TypeOf(New())
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestPackageManagerAvailableUpdatesSuccess(t *testing.T) {
	execCommand = fakeExecCommand("TestPackageManagerAvailableUpdatesCommandSuccess")
	defer func() { execCommand = exec.Command }()
	pm := XBPS{}

	_, got := pm.AvailableUpdates()
	if got != nil {
		t.Errorf("want %v, got %v", nil, got)
	}
}

func TestPackageManagerAvailableUpdatesError(t *testing.T) {
	execCommand = fakeExecCommand("TestPackageManagerCommandError")
	defer func() { execCommand = exec.Command }()
	pm := XBPS{}

	_, got := pm.AvailableUpdates()
	if got == nil {
		t.Errorf("want %v, got %v", "exit status 1", got)
	}
}

func TestPackageManagerInstalledSuccess(t *testing.T) {
	execCommand = fakeExecCommand("TestPackageManagerInstalledCommandSuccess")
	defer func() { execCommand = exec.Command }()
	pm := XBPS{}

	_, got := pm.Installed()
	if got != nil {
		t.Errorf("want %v, got %v", nil, got)
	}
}

func TestPackageManagerInstalledError(t *testing.T) {
	execCommand = fakeExecCommand("TestPackageManagerCommandError")
	defer func() { execCommand = exec.Command }()
	pm := XBPS{}

	_, got := pm.Installed()
	if got == nil {
		t.Errorf("want %v, got %v", "exit status 1", got)
	}
}

func TestPackageManagerCommandError(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}

	os.Exit(1)
}

func TestPackageManagerAvailableUpdatesCommandSuccess(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}

	fmt.Println("package-name-3.5.2_1 update x86_64 https://alpha.de.repo.voidlinux.org/current 182572180 59477905")
	os.Exit(0)
}

func TestPackageManagerInstalledCommandSuccess(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}

	fmt.Println("package-name-3.5.2_1")
	os.Exit(0)
}

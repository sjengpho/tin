package packagemanager

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/sjengpho/tin/tin"
)

func fakeLookPathXBPS(file string) (string, error) {
	if file != "xbps-install" {
		return "", errors.New("executeable doesn't exists")
	}

	return "fake-path", nil
}

func fakeLookPathPacman(file string) (string, error) {
	if file != "checkupdates" {
		return "", errors.New("executeable doesn't exists")
	}

	return "fake-path", nil
}

func fakeLookPathError(file string) (string, error) {
	return "", errors.New("executeable doesn't exists")
}

type fakeYay struct{}

func (f fakeYay) Installed() ([]tin.Package, error) {
	return []tin.Package{}, errors.New("error")
}

func (f fakeYay) AvailableUpdates() ([]tin.Package, error) {
	return []tin.Package{}, errors.New("error")
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
	tests := []struct {
		want         tin.PackageManager
		fakeLookPath func(file string) (string, error)
	}{
		{
			want:         &XBPS{},
			fakeLookPath: fakeLookPathXBPS,
		},
		{
			want:         &Arch{},
			fakeLookPath: fakeLookPathPacman,
		},
	}

	for _, tt := range tests {
		lookPath = tt.fakeLookPath

		got := reflect.TypeOf(New())
		want := reflect.TypeOf(tt.want)
		if got != want {
			t.Errorf("want %v, got %v", want, got)
		}

		lookPath = exec.LookPath
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

func TestAvailableUpdatesSuccess(t *testing.T) {
	tests := []struct {
		pm              tin.PackageManager
		fakeExecCommand func(name string, args ...string) *exec.Cmd
	}{
		{
			pm:              &XBPS{},
			fakeExecCommand: fakeExecCommand("TestXBPSAvailableUpdatesCommandSuccess"),
		},
		{
			pm:              &Arch{Pacman: Pacman{}, AUR: &Yay{}},
			fakeExecCommand: fakeExecCommand("TestArchCommandSuccess"),
		},
		{
			pm:              &Pacman{},
			fakeExecCommand: fakeExecCommand("TestArchCommandSuccess"),
		},
		{
			pm:              &Yay{},
			fakeExecCommand: fakeExecCommand("TestArchCommandSuccess"),
		},
	}

	for _, tt := range tests {
		execCommand = tt.fakeExecCommand

		_, got := tt.pm.AvailableUpdates()
		if got != nil {
			t.Errorf("want %v, got %v", nil, got)
		}

		execCommand = exec.Command
	}
}

func TestAvailableUpdatesError(t *testing.T) {
	tests := []struct {
		pm              tin.PackageManager
		fakeExecCommand func(name string, args ...string) *exec.Cmd
	}{
		{
			pm:              &XBPS{},
			fakeExecCommand: fakeExecCommand("TestCommandError"),
		},
		{
			pm:              &Arch{Pacman: Pacman{}, AUR: &Yay{}},
			fakeExecCommand: fakeExecCommand("TestCommandError"),
		},
		{
			pm:              &Arch{Pacman: Pacman{}, AUR: fakeYay{}},
			fakeExecCommand: fakeExecCommand("TestArchCommandSuccess"),
		},
		{
			pm:              &Pacman{},
			fakeExecCommand: fakeExecCommand("TestCommandError"),
		},
		{
			pm:              &Yay{},
			fakeExecCommand: fakeExecCommand("TestCommandError"),
		},
	}

	for _, tt := range tests {
		execCommand = tt.fakeExecCommand

		_, got := tt.pm.AvailableUpdates()
		if got == nil {
			t.Errorf("want %v, got %v", nil, got)
		}

		execCommand = exec.Command
	}
}

func TestInstalledSuccess(t *testing.T) {
	tests := []struct {
		pm              tin.PackageManager
		fakeExecCommand func(name string, args ...string) *exec.Cmd
	}{
		{
			pm:              &XBPS{},
			fakeExecCommand: fakeExecCommand("TestXBPSInstalledCommandSuccess"),
		},
		{
			pm:              &Arch{Pacman: Pacman{}, AUR: &Yay{}},
			fakeExecCommand: fakeExecCommand("TestArchCommandSuccess"),
		},
		{
			pm:              &Pacman{},
			fakeExecCommand: fakeExecCommand("TestArchCommandSuccess"),
		},
	}

	for _, tt := range tests {
		execCommand = tt.fakeExecCommand

		_, got := tt.pm.Installed()
		if got != nil {
			t.Errorf("want %v, got %v", nil, got)
		}

		execCommand = exec.Command
	}
}

func TestInstalledError(t *testing.T) {
	tests := []struct {
		pm              tin.PackageManager
		fakeExecCommand func(name string, args ...string) *exec.Cmd
	}{
		{
			pm:              &XBPS{},
			fakeExecCommand: fakeExecCommand("TestCommandError"),
		},
		{
			pm:              &Arch{Pacman: Pacman{}, AUR: &Yay{}},
			fakeExecCommand: fakeExecCommand("TestCommandError"),
		},
		{
			pm:              &Pacman{},
			fakeExecCommand: fakeExecCommand("TestCommandError"),
		},
		{
			pm:              &Yay{},
			fakeExecCommand: fakeExecCommand("TestCommandError"),
		},
	}

	for _, tt := range tests {
		execCommand = tt.fakeExecCommand

		_, got := tt.pm.Installed()
		if got == nil {
			t.Errorf("want %v, got %v", nil, got)
		}

		execCommand = exec.Command
	}
}

func TestCommandError(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}

	os.Exit(1)
}

func TestXBPSAvailableUpdatesCommandSuccess(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}

	fmt.Println("package-name-3.5.2_1 update x86_64 https://alpha.de.repo.voidlinux.org/current 182572180 59477905")
	os.Exit(0)
}

func TestXBPSInstalledCommandSuccess(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}

	fmt.Println("package-name-3.5.2_1")
	os.Exit(0)
}

func TestArchCommandSuccess(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}

	fmt.Println("package-name 3.5.2-1")
	os.Exit(0)
}

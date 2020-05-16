package temperature

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

type fakeFile struct{}

func (f fakeFile) Name() string {
	return "file name"
}

func (f fakeFile) Size() int64 {
	return 1
}

func (f fakeFile) Mode() os.FileMode {
	return os.ModeTemporary
}

func (f fakeFile) ModTime() time.Time {
	return time.Now()
}

func (f fakeFile) IsDir() bool {
	return f.Mode().IsDir()
}

func (f fakeFile) Sys() interface{} {
	return nil
}

func fakeOsStatSuccess(name string) (os.FileInfo, error) {
	return fakeFile{}, nil
}

func fakeOsStatError(name string) (os.FileInfo, error) {
	return nil, errors.New("executeable doesn't exists")
}

func fakeReadFileSuccess(filename string) ([]byte, error) {
	return []byte("30000"), nil
}

func fakeReadFileError(filename string) ([]byte, error) {
	return nil, errors.New("read file error")
}

func fakeReadFileInvalidContent(filename string) ([]byte, error) {
	return []byte("this should be a millidegree in celsius format"), nil
}

func TestNewReaderSuccess(t *testing.T) {
	osStat = fakeOsStatSuccess
	defer func() { osStat = os.Stat }()

	want := reflect.TypeOf(&FileReader{})
	got := reflect.TypeOf(NewReader())
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestNewReaderError(t *testing.T) {
	osStat = fakeOsStatError
	defer func() { osStat = os.Stat }()

	want := reflect.TypeOf(nil)
	got := reflect.TypeOf(NewReader())
	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestReaderReadSuccess(t *testing.T) {
	readFile = fakeReadFileSuccess
	defer func() { readFile = ioutil.ReadFile }()

	reader := FileReader{path: ""}
	want := 30
	got, err := reader.Read()
	if err != nil {
		t.Errorf("want %v, got %v", nil, err)
	}

	if got.Value != want {
		t.Errorf("want %v, got %v", want, got.Value)
	}
}

func TestReaderReadReadFileError(t *testing.T) {
	readFile = fakeReadFileError
	defer func() { readFile = ioutil.ReadFile }()

	reader := FileReader{path: ""}
	_, got := reader.Read()
	if got == nil {
		t.Errorf("want %v, got %v", "error", got)
	}
}

func TestReaderReadParseError(t *testing.T) {
	readFile = fakeReadFileInvalidContent
	defer func() { readFile = ioutil.ReadFile }()

	reader := FileReader{path: ""}
	_, got := reader.Read()
	if got == nil {
		t.Errorf("want %v, got %v", "error", got)
	}
}

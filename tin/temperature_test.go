package tin

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

type temperatureReaderMock struct{ returnError bool }

func (r temperatureReaderMock) Read() (Temperature, error) {
	if r.returnError {
		return Temperature{}, errors.New("error")
	}

	return Temperature{Value: 17}, nil
}

func TestNewTemperatureService(t *testing.T) {
	tt := []struct {
		want *TemperatureService
		got  *TemperatureService
	}{
		{
			want: &TemperatureService{},
			got:  NewTemperatureService(temperatureReaderMock{returnError: false}, log.New(os.Stdout, "", log.Flags())),
		},
		{
			want: &TemperatureService{},
			got:  NewTemperatureService(temperatureReaderMock{returnError: true}, log.New(os.Stdout, "", log.Flags())),
		},
		{
			want: &TemperatureService{},
			got:  NewTemperatureService(nil, log.New(os.Stdout, "", log.Flags())),
		},
	}

	for _, tc := range tt {
		if reflect.TypeOf(tc.got) != reflect.TypeOf(tc.want) {
			t.Errorf("want %v, got %v", reflect.TypeOf(tc.want), reflect.TypeOf(tc.got))
		}
	}
}

func TestTemperatureTemperature(t *testing.T) {
	withState := NewTemperatureService(nil, log.New(ioutil.Discard, "", log.Flags()))
	withState.SetTemperature(Temperature{Value: 17})

	tt := []struct {
		service *TemperatureService
		want    Temperature
	}{
		{
			service: withState,
			want:    Temperature{Value: 17},
		},
		{
			service: NewTemperatureService(nil, log.New(ioutil.Discard, "", log.Flags())),
			want:    Temperature{},
		},
	}

	for _, tc := range tt {
		got := tc.service.Temperature()

		if got != tc.want {
			t.Errorf("want %v, got %v", tc.want, got)
		}
	}
}

func TestTemperatureCelsius(t *testing.T) {
	temperature := Temperature{Value: 17}

	want := 17
	got := temperature.Celsius()

	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestTemperatureFahrenheit(t *testing.T) {
	temperature := Temperature{Value: 17}

	want := 62
	got := temperature.Fahrenheit()

	if got != want {
		t.Errorf("want %v, got %v", want, got)
	}
}

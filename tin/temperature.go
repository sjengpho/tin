package tin

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// TemperatureReader is the interface implemented by an object that can
// return the current temperature of the hardware.
type TemperatureReader interface {
	Read() (Temperature, error)
}

// Temperature represents the temperature.
type Temperature struct {
	Value int // Celsius
}

// Equal implements tin.Comparable.
func (t Temperature) Equal(v interface{}) bool {
	if b, ok := v.(Temperature); ok {
		return t.Value == b.Value
	}
	return false
}

// Celsius returns the temperature in Celsius format.
func (t *Temperature) Celsius() int {
	return t.Value
}

// Fahrenheit returns the temperature in Fahrenheit format.
func (t *Temperature) Fahrenheit() int {
	return (t.Value * 9 / 5) + 32
}

// Temp represents a StateKey.
const Temp StateKey = "Temperature"

// TemperatureService provides access to the temperature.
type TemperatureService struct {
	Reader TemperatureReader
	Worker *Worker
	state  *State
	logger *log.Logger
}

// NewTemperatureService returns a tin.TemperatureService.
func NewTemperatureService(r TemperatureReader, l *log.Logger) *TemperatureService {
	s := &TemperatureService{
		Reader: r,
		state:  NewState(),
		logger: l,
	}

	// Worker that reads the temperature on intervals and updates the state.
	if r == nil {
		s.logger.Println(errors.New("failed initializing worker"))
	} else {
		s.Worker = NewWorker(10*time.Second, func() {
			t, err := s.Reader.Read()
			if err != nil {
				s.logger.Println(fmt.Errorf("worker failed: %w", err))
			} else {
				s.SetTemperature(t)
			}
		})
	}

	return s
}

// Temperature returns a tin.Temperature.
func (s *TemperatureService) Temperature() Temperature {
	v, err := s.state.Get(Temp)
	if err != nil {
		return Temperature{Value: 0}
	}

	return v.(Temperature)
}

// SetTemperature updates the state.
func (s *TemperatureService) SetTemperature(t Temperature) {
	s.state.Set(Temp, t)
}

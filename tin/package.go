package tin

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// PackageManager is the interface implemented by an object that can
// fetch info about system packages.
//
// AvailableUpdates returns the packages that are updateable.
// Installed returns the packages that are currently installed.
type PackageManager interface {
	AvailableUpdates() ([]Package, error)
	Installed() ([]Package, error)
}

// PackageManagerServiceState represents the state.
type PackageManagerServiceState struct {
	sync.RWMutex
	AvailableUpdates PackageCount
}

// Packages represents a slice of tin.Package.
type Packages []Package

// Equal implements tin.Comparable.
func (a Packages) Equal(t interface{}) bool {
	b, ok := t.(Packages)
	if !ok || len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if !v.Equal(b[i]) {
			return false
		}
	}
	return true
}

// Package represents a package from a package manager.
type Package struct {
	Name    string
	Version string
}

// Equal implements tin.Comparable.
func (a Package) Equal(t interface{}) bool {
	if b, ok := t.(Package); ok {
		return a == b
	}
	return false
}

// PackageCount represents the amount of packages.
type PackageCount int

// Equal implements tin.Comparable.
func (a PackageCount) Equal(t interface{}) bool {
	if b, ok := t.(PackageCount); ok {
		return a == b
	}
	return false
}

// Represents a tin.StateKey.
const (
	AvailableUpdates StateKey = "AvailableUpdates"
	Installed                 = "Installed"
)

// PackageManagerService provides access to data from package managers.
type PackageManagerService struct {
	manager PackageManager
	state   *State
	worker  *Worker
	logger  *log.Logger
}

// NewPackageManagerService returns a tin.PackageManagerService.
func NewPackageManagerService(m PackageManager, l *log.Logger) *PackageManagerService {
	s := &PackageManagerService{
		manager: m,
		state:   NewState(),
		logger:  l,
	}

	// Worker that fetches available package updates on intervals and updates the state.
	if m == nil {
		s.logger.Println(errors.New("failed initializing worker"))
	} else {
		s.worker = NewWorker(time.Minute, func() {
			packages, err := s.manager.AvailableUpdates()
			if err != nil {
				s.logger.Println(fmt.Errorf("worker failed: %w", err))
			} else {
				s.SetAvailableUpdates(PackageCount(len(packages)))
			}

		})
	}

	// Worker that fetches available package updates on intervals and updates the state.
	if m == nil {
		s.logger.Println(errors.New("failed initializing worker"))
	} else {
		s.worker = NewWorker(time.Minute, func() {
			packages, err := s.manager.Installed()
			if err != nil {
				s.logger.Println(fmt.Errorf("worker failed: %w", err))
			} else {
				s.SetInstalled(Packages(packages))
			}

		})
	}

	return s
}

// Subscribe returns a tin.StateSubscription.
func (s *PackageManagerService) Subscribe() StateSubscription {
	return s.state.Subscribe()
}

// SetAvailableUpdates updates the state.
func (s *PackageManagerService) SetAvailableUpdates(c PackageCount) {
	s.state.Set(AvailableUpdates, c)
}

// AvailableUpdatesCount returns a tin.PackageCount.
func (s *PackageManagerService) AvailableUpdatesCount() PackageCount {
	v, err := s.state.Get(AvailableUpdates)
	if err != nil {
		return PackageCount(0)
	}

	return v.(PackageCount)
}

// SetInstalled updates the state.
func (s *PackageManagerService) SetInstalled(p Packages) {
	s.state.Set(Installed, p)
}

// Installed returns a slice of tin.Package.
func (s *PackageManagerService) Installed() Packages {
	v, err := s.state.Get(Installed)
	if err != nil {
		return Packages{}
	}

	return v.(Packages)
}

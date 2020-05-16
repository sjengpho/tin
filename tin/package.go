package tin

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// PackageManager is the interface implemented by an object that can
// return the amount of available updates.
type PackageManager interface {
	AvailableUpdates() ([]Package, error)
}

// PackageManagerServiceState represents the state.
type PackageManagerServiceState struct {
	sync.RWMutex
	AvailableUpdates PackageCount
}

// Package represents a package from a package manager.
type Package struct {
	Name    string
	Version string
}

// PackageCount represents the amount of packages.
type PackageCount = int

// AvailableUpdates represents a tin.StateKey.
const AvailableUpdates StateKey = "AvailableUpdates"

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

	return s
}

// Subscribe returns a read-only channel.
func (s *PackageManagerService) Subscribe() <-chan interface{} {
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

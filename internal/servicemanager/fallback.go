//go:build !linux && !windows

package servicemanager

func New() Manager {
	return NewMockManager()
}

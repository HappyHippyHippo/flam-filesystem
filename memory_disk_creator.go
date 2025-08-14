package filesystem

import (
	"github.com/spf13/afero"

	flam "github.com/happyhippyhippo/flam"
)

type memoryDiskCreator struct{}

func newMemoryDiskCreator() DiskCreator {
	return &memoryDiskCreator{}
}

func (memoryDiskCreator) Accept(
	config flam.Bag,
) bool {
	return config.String("driver") == DiskDriverMemory
}

func (memoryDiskCreator) Create(
	_ flam.Bag,
) (Disk, error) {
	return afero.NewMemMapFs(), nil
}

package filesystem

import (
	flam "github.com/happyhippyhippo/flam"
)

type Facade interface {
	HasDisk(id string) bool
	ListDisks() []string
	GetDisk(id string) (Disk, error)
	AddDisk(id string, disk Disk) error
}

type facade struct {
	diskFactory diskFactory
}

func newFacade(
	diskFactory diskFactory,
) Facade {
	return &facade{
		diskFactory: diskFactory,
	}
}

func (facade facade) HasDisk(
	id string,
) bool {
	return facade.diskFactory.Has(id)
}

func (facade facade) ListDisks() []string {
	return facade.diskFactory.List()
}

func (facade facade) GetDisk(
	id string,
) (Disk, error) {
	return facade.diskFactory.Get(id)
}

func (facade facade) AddDisk(
	id string,
	disk Disk,
) error {
	if disk == nil {
		return flam.ErrNilReference
	}
	return facade.diskFactory.Add(id, disk)
}

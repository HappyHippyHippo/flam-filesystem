package filesystem

import (
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
)

type provider struct{}

func NewProvider() flam.Provider {
	return &provider{}
}

func (*provider) Id() string {
	return providerId
}

func (*provider) Register(
	container *dig.Container,
) error {
	if container == nil {
		return newErrNilReference("container")
	}

	registerer := flam.NewRegisterer()
	registerer.Queue(newOsDiskCreator, dig.Group(DiskCreatorGroup))
	registerer.Queue(newMemoryDiskCreator, dig.Group(DiskCreatorGroup))
	registerer.Queue(newDiskFactory)
	registerer.Queue(newFacade)

	return registerer.Run(container)
}

func (provider *provider) Close(
	container *dig.Container,
) error {
	if container == nil {
		return newErrNilReference("container")
	}

	executor := flam.NewExecutor()
	executor.Queue(provider.closeDiskFactory)

	return executor.Run(container)
}

func (*provider) closeDiskFactory(
	diskFactory diskFactory,
) error {
	return diskFactory.Close()
}

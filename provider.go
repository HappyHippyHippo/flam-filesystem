package filesystem

import (
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
)

type provider struct{}

func NewProvider() flam.Provider {
	return &provider{}
}

func (provider) Id() string {
	return providerId
}

func (provider) Register(
	container *dig.Container,
) error {
	if container == nil {
		return newErrNilReference("container")
	}

	var e error
	provide := func(constructor any, opts ...dig.ProvideOption) bool {
		e = container.Provide(constructor, opts...)
		return e == nil
	}

	_ = provide(newOsDiskCreator, dig.Group(DiskCreatorGroup)) &&
		provide(newMemoryDiskCreator, dig.Group(DiskCreatorGroup)) &&
		provide(newDiskFactory) &&
		provide(newFacade)

	return e
}

func (provider) Close(
	container *dig.Container,
) error {
	if container == nil {
		return newErrNilReference("container")
	}

	return container.Invoke(func(
		diskFactory diskFactory,
	) error {
		return diskFactory.Close()
	})
}

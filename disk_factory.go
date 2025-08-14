package filesystem

import (
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
)

type diskFactory = flam.Factory[Disk]

type diskFactoryArgs struct {
	dig.In

	Creators      []DiskCreator `group:"flam.filesystem.disks.creator"`
	FactoryConfig flam.FactoryConfig
}

func newDiskFactory(
	args diskFactoryArgs,
) (diskFactory, error) {
	var creators []flam.ResourceCreator[Disk]
	for _, creator := range args.Creators {
		creators = append(creators, creator)
	}

	return flam.NewFactory(
		creators,
		PathDisks,
		args.FactoryConfig,
		nil)
}

package filesystem

import (
	"github.com/spf13/afero"

	flam "github.com/happyhippyhippo/flam"
)

type osDiskCreator struct{}

func newOsDiskCreator() DiskCreator {
	return &osDiskCreator{}
}

func (osDiskCreator) Accept(
	config flam.Bag,
) bool {
	return config.String("driver") == DiskDriverOS
}

func (osDiskCreator) Create(
	_ flam.Bag,
) (Disk, error) {
	return afero.NewOsFs(), nil
}

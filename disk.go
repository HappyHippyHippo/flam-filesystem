package filesystem

import (
	"github.com/spf13/afero"
)

type Disk interface {
	afero.Fs
}

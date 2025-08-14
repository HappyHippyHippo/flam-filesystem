package filesystem

import (
	"github.com/spf13/afero"
)

type File interface {
	afero.File
}

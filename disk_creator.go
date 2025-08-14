package filesystem

import (
	flam "github.com/happyhippyhippo/flam"
)

type DiskCreator interface {
	Accept(config flam.Bag) bool
	Create(config flam.Bag) (Disk, error)
}

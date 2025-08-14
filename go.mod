module github.com/happyhippyhippo/flam-filesystem

go 1.24.0

replace github.com/happyhippyhippo/flam => ../flam

require (
	github.com/golang/mock v1.6.0
	github.com/happyhippyhippo/flam v0.0.0-00010101000000-000000000000
	github.com/spf13/afero v1.14.0
	go.uber.org/dig v1.19.0
)

require (
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.8.4
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

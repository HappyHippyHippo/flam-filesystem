package tests

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
	filesystem "github.com/happyhippyhippo/flam-filesystem"
	mocks "github.com/happyhippyhippo/flam-filesystem/tests/mocks"
)

func Test_NewProvider(t *testing.T) {
	assert.NotNil(t, filesystem.NewProvider())
}

func Test_Provider_Id(t *testing.T) {
	assert.Equal(t, "flam.filesystem.provider", filesystem.NewProvider().Id())
}

func Test_Provider_Register(t *testing.T) {
	t.Run("should return error if nil container is passed", func(t *testing.T) {
		assert.ErrorIs(t, filesystem.NewProvider().Register(nil), flam.ErrNilReference)
	})

	t.Run("should successfully provides Facade", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))

		factoryConfig := mocks.NewFactoryConfig(ctrl)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade filesystem.Facade) {
			assert.NotNil(t, facade)
		}))
	})
}

func Test_Provider_Close(t *testing.T) {
	t.Run("should return error if nil container is passed", func(t *testing.T) {
		assert.ErrorIs(
			t,
			filesystem.NewProvider().(flam.ClosableProvider).Close(nil),
			flam.ErrNilReference)
	})

	t.Run("should return instantiated disk closing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))

		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.
			EXPECT().
			Get(filesystem.PathDisks).
			Return(flam.Bag{"mock": flam.Bag{}}).
			Times(1)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfig
		}))

		expectedErr := errors.New("mock error")
		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Close().Return(expectedErr).Times(1)

		diskCreatorConfig := flam.Bag{"id": "mock"}
		diskCreator := mocks.NewDiskCreator(ctrl)
		diskCreator.EXPECT().Accept(diskCreatorConfig).Return(true).Times(1)
		diskCreator.EXPECT().Create(diskCreatorConfig).Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.DiskCreator {
			return diskCreator
		}, dig.Group(filesystem.DiskCreatorGroup)))

		assert.NoError(t, container.Invoke(func(facade filesystem.Facade) error {
			got, e := facade.GetDisk("mock")
			assert.NotNil(t, got)
			assert.NoError(t, e)

			return e
		}))

		assert.ErrorIs(
			t,
			filesystem.NewProvider().(flam.ClosableProvider).Close(container),
			expectedErr)
	})

	t.Run("should successfully close disk", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, filesystem.NewProvider().Register(container))

		factoryConfig := mocks.NewFactoryConfig(ctrl)
		factoryConfig.
			EXPECT().
			Get(filesystem.PathDisks).
			Return(flam.Bag{"mock": flam.Bag{}}).
			Times(1)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfig
		}))

		disk := mocks.NewDisk(ctrl)
		disk.EXPECT().Close().Return(nil).Times(1)

		diskCreatorConfig := flam.Bag{"id": "mock"}
		diskCreator := mocks.NewDiskCreator(ctrl)
		diskCreator.EXPECT().Accept(diskCreatorConfig).Return(true).Times(1)
		diskCreator.EXPECT().Create(diskCreatorConfig).Return(disk, nil).Times(1)
		require.NoError(t, container.Provide(func() filesystem.DiskCreator {
			return diskCreator
		}, dig.Group(filesystem.DiskCreatorGroup)))

		assert.NoError(t, container.Invoke(func(facade filesystem.Facade) error {
			got, e := facade.GetDisk("mock")
			assert.NotNil(t, got)
			assert.NoError(t, e)

			return e
		}))

		assert.NoError(t, filesystem.NewProvider().(flam.ClosableProvider).Close(container))
	})
}

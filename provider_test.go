package filesystem

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
)

func Test_NewProvider(t *testing.T) {
	assert.NotNil(t, NewProvider())
}

func Test_Provider_Id(t *testing.T) {
	assert.Equal(t, "flam.filesystem.provider", NewProvider().Id())
}

func Test_Provider_Register(t *testing.T) {
	t.Run("should return error if nil container is passed", func(t *testing.T) {
		assert.ErrorIs(t, NewProvider().Register(nil), flam.ErrNilReference)
	})

	t.Run("should successfully provides Facade", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		factoryConfigMock := NewFactoryConfigMock(ctrl)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfigMock
		}))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.NotNil(t, facade)
		}))
	})
}

func Test_Provider_Close(t *testing.T) {
	t.Run("should return error if nil container is passed", func(t *testing.T) {
		assert.ErrorIs(
			t,
			NewProvider().(flam.ClosableProvider).Close(nil),
			flam.ErrNilReference)
	})

	t.Run("should return instantiated disk closing error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		factoryConfigMock := NewFactoryConfigMock(ctrl)
		factoryConfigMock.
			EXPECT().
			Get(PathDisks).
			Return(flam.Bag{"mock": flam.Bag{}}).
			Times(1)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfigMock
		}))

		expectedErr := errors.New("mock error")
		diskMock := NewDiskMock(ctrl)
		diskMock.EXPECT().Close().Return(expectedErr).Times(1)

		diskCreatorConfig := flam.Bag{"id": "mock"}
		diskCreatorMock := NewDiskCreatorMock(ctrl)
		diskCreatorMock.EXPECT().Accept(diskCreatorConfig).Return(true).Times(1)
		diskCreatorMock.EXPECT().Create(diskCreatorConfig).Return(diskMock, nil).Times(1)
		require.NoError(t, container.Provide(func() DiskCreator {
			return diskCreatorMock
		}, dig.Group(DiskCreatorGroup)))

		assert.NoError(t, container.Invoke(func(facade Facade) error {
			got, e := facade.GetDisk("mock")
			assert.NotNil(t, got)
			assert.NoError(t, e)

			return e
		}))

		assert.ErrorIs(
			t,
			NewProvider().(flam.ClosableProvider).Close(container),
			expectedErr)
	})

	t.Run("should successfully close disk", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		factoryConfigMock := NewFactoryConfigMock(ctrl)
		factoryConfigMock.
			EXPECT().
			Get(PathDisks).
			Return(flam.Bag{"mock": flam.Bag{}}).
			Times(1)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfigMock
		}))

		diskMock := NewDiskMock(ctrl)
		diskMock.EXPECT().Close().Return(nil).Times(1)

		diskCreatorConfig := flam.Bag{"id": "mock"}
		diskCreatorMock := NewDiskCreatorMock(ctrl)
		diskCreatorMock.EXPECT().Accept(diskCreatorConfig).Return(true).Times(1)
		diskCreatorMock.EXPECT().Create(diskCreatorConfig).Return(diskMock, nil).Times(1)
		require.NoError(t, container.Provide(func() DiskCreator {
			return diskCreatorMock
		}, dig.Group(DiskCreatorGroup)))

		assert.NoError(t, container.Invoke(func(facade Facade) error {
			got, e := facade.GetDisk("mock")
			assert.NotNil(t, got)
			assert.NoError(t, e)

			return e
		}))

		assert.NoError(t, NewProvider().(flam.ClosableProvider).Close(container))
	})
}

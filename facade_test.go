package filesystem

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"

	flam "github.com/happyhippyhippo/flam"
)

func Test_Facade_HasDisk(t *testing.T) {
	t.Run("should return false on unknown disk", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		config := flam.Bag{}
		factoryConfig := NewFactoryConfigMock(ctrl)
		factoryConfig.EXPECT().Get(PathDisks).Return(config).Times(1)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.False(t, facade.HasDisk("mock"))
		}))
	})

	t.Run("should return true on known disk", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		config := flam.Bag{"mock": flam.Bag{}}
		factoryConfig := NewFactoryConfigMock(ctrl)
		factoryConfig.EXPECT().Get(PathDisks).Return(config).Times(1)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.True(t, facade.HasDisk("mock"))
		}))
	})

	t.Run("should return true on added disk", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		config := flam.Bag{}
		factoryConfig := NewFactoryConfigMock(ctrl)
		factoryConfig.EXPECT().Get(PathDisks).Return(config).Times(1)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfig
		}))

		disk := NewDiskMock(ctrl)

		assert.NoError(t, container.Invoke(func(facade Facade) {
			require.NoError(t, facade.AddDisk("mock", disk))

			assert.True(t, facade.HasDisk("mock"))
		}))
	})
}

func Test_Facade_ListDisks(t *testing.T) {
	t.Run("should return empty list on empty config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		config := flam.Bag{}
		factoryConfig := NewFactoryConfigMock(ctrl)
		factoryConfig.EXPECT().Get(PathDisks).Return(config).Times(1)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.Empty(t, facade.ListDisks())
		}))
	})

	t.Run("should return a sorted list of disks", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		config := flam.Bag{
			"gamma": flam.Bag{},
			"alpha": flam.Bag{},
			"beta":  flam.Bag{}}
		factoryConfig := NewFactoryConfigMock(ctrl)
		factoryConfig.EXPECT().Get(PathDisks).Return(config).Times(1)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.ElementsMatch(
				t,
				[]string{"alpha", "beta", "gamma"},
				facade.ListDisks())
		}))
	})

	t.Run("should return a sorted list of disks (with added disks)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		config := flam.Bag{
			"gamma": flam.Bag{},
			"alpha": flam.Bag{},
			"beta":  flam.Bag{},
		}
		factoryConfig := NewFactoryConfigMock(ctrl)
		factoryConfig.EXPECT().Get(PathDisks).Return(config).Times(2)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfig
		}))

		disk := NewDiskMock(ctrl)

		assert.NoError(t, container.Invoke(func(facade Facade) {
			require.NoError(t, facade.AddDisk("delta", disk))

			assert.ElementsMatch(
				t,
				[]string{"alpha", "beta", "delta", "gamma"},
				facade.ListDisks())
		}))
	})
}

func Test_Facade_GetDisk(t *testing.T) {
	t.Run("should return error on unknown disk", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		config := flam.Bag{}
		factoryConfig := NewFactoryConfigMock(ctrl)
		factoryConfig.EXPECT().Get(PathDisks).Return(config).Times(1)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.GetDisk("mock")
			assert.Nil(t, got)
			assert.ErrorIs(t, e, flam.ErrUnknownResource)
		}))
	})

	t.Run("should return 'os' driver disk", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		config := flam.Bag{"disk": flam.Bag{"driver": DiskDriverOS}}
		factoryConfig := NewFactoryConfigMock(ctrl)
		factoryConfig.EXPECT().Get(PathDisks).Return(config).Times(1)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.GetDisk("disk")
			require.NotNil(t, got)
			require.NoError(t, e)

			assert.IsType(t, afero.NewOsFs(), got)
		}))
	})

	t.Run("should return 'memory' driver disk", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		config := flam.Bag{"disk": flam.Bag{"driver": DiskDriverMemory}}
		factoryConfig := NewFactoryConfigMock(ctrl)
		factoryConfig.EXPECT().Get(PathDisks).Return(config).Times(1)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			got, e := facade.GetDisk("disk")
			require.NotNil(t, got)
			require.NoError(t, e)

			assert.IsType(t, afero.NewMemMapFs(), got)
		}))
	})

	t.Run("should return added disk", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		config := flam.Bag{}
		factoryConfig := NewFactoryConfigMock(ctrl)
		factoryConfig.EXPECT().Get(PathDisks).Return(config).Times(1)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfig
		}))

		disk := afero.NewMemMapFs()

		assert.NoError(t, container.Invoke(func(facade Facade) {
			require.NoError(t, facade.AddDisk("disk", disk))

			got, e := facade.GetDisk("disk")
			assert.Same(t, disk, got)
			assert.NoError(t, e)
		}))
	})
}

func Test_Facade_AddDisk(t *testing.T) {
	t.Run("should return error on nil disk", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		factoryConfig := NewFactoryConfigMock(ctrl)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfig
		}))

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.ErrorIs(t, facade.AddDisk("disk", nil), flam.ErrNilReference)
		}))
	})

	t.Run("should return error on existing disk", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		config := flam.Bag{"disk": flam.Bag{}}
		factoryConfig := NewFactoryConfigMock(ctrl)
		factoryConfig.EXPECT().Get(PathDisks).Return(config).Times(1)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfig
		}))

		disk := NewDiskMock(ctrl)

		assert.NoError(t, container.Invoke(func(facade Facade) {
			assert.ErrorIs(t, facade.AddDisk("disk", disk), flam.ErrDuplicateResource)
		}))
	})

	t.Run("should correctly add the disk", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		container := dig.New()
		require.NoError(t, NewProvider().Register(container))

		config := flam.Bag{}
		factoryConfig := NewFactoryConfigMock(ctrl)
		factoryConfig.EXPECT().Get(PathDisks).Return(config).Times(1)
		require.NoError(t, container.Provide(func() flam.FactoryConfig {
			return factoryConfig
		}))

		disk := NewDiskMock(ctrl)

		assert.NoError(t, container.Invoke(func(facade Facade) {
			require.NoError(t, facade.AddDisk("disk", disk))

			got, e := facade.GetDisk("disk")
			assert.Same(t, got, disk)
			assert.NoError(t, e)
		}))
	})
}

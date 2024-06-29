package gobs_test

import (
	"context"
	"fmt"

	"github.com/stretchr/testify/require"
	"github.com/traphamxuan/gobs"
	"github.com/traphamxuan/gobs/common"
)

func (s *BootstrapSuit) TestAsyncPerformance() {
	t := s.T()
	numOfDependencies := 10
	level := 3
	// shared := 5
	ctx := context.TODO()
	service := NewSampleAsyncService(numOfDependencies, level)
	bs := gobs.NewBootstrap(gobs.Config{
		IsConcurrent: true,
	})
	require.NoError(t, bs.AddDefault(service), "AddDefault expected no error")
	require.NoError(t, bs.Init(ctx), "Init expected no error")
	require.NoError(t, bs.Setup(ctx), "Setup expected no error")
	require.NoError(t, bs.Start(ctx), "Setup expected no error")
	require.NoError(t, bs.Stop(ctx), "Stop expected no error")
}

type SampleAsyncService struct {
	level     int
	id        int
	numOfDeps int
	// dependencies []gobs.IService
}

// Init implements gobs.IService.
func (s *SampleAsyncService) Init(ctx context.Context, c *gobs.Service) error {
	numOfServices++
	if s.level > 0 {
		deps := make([]SampleAsyncService, s.numOfDeps)
		c.ExtraDeps = make([]gobs.CustomService, s.numOfDeps)
		newLevel := s.level - 1
		for i := 0; i < s.numOfDeps; i++ {
			deps[i].level = newLevel
			deps[i].numOfDeps = s.numOfDeps
			deps[i].id = numOfServices
			c.ExtraDeps[i] = gobs.CustomService{
				Name:     fmt.Sprintf("Sample-%d-%d-%d", newLevel, i, numOfServices),
				Instance: &deps[i],
			}
		}
	}
	c.AsyncMode[common.StatusSetup] = true

	return nil
}

var _ gobs.IService = (*SampleAsyncService)(nil)

func NewSampleAsyncService(numOfDeps, level int) gobs.IService {
	return &SampleAsyncService{
		level:     level,
		id:        0,
		numOfDeps: numOfDeps,
	}
}

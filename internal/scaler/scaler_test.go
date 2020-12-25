package scaler_test

import (
	"errors"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/tdc-yamada-ya/aks-scheduled-poolscaler/internal/mock_scaler"
	"github.com/tdc-yamada-ya/aks-scheduled-poolscaler/internal/scaler"
)

type poolUpdaterUpdateCall struct {
	inError              error
	outResourceGroupName scaler.ResourceGroupName
	outResourceName      scaler.ResourceName
	outAgentPoolName     scaler.AgentPoolName
	outParameters        *scaler.Parameters
}

func TestScalerScale(t *testing.T) {
	bpf := func(b bool) *bool { return &b }
	ipf := func(i int32) *int32 { return &i }
	p1 := &scaler.Parameters{bpf(false), ipf(0), ipf(0), ipf(0)}
	p2 := &scaler.Parameters{bpf(true), ipf(1), ipf(1), ipf(1)}
	configuration := &scaler.Configuration{
		scaler.ParametersMap{
			"pn1": p1,
			"pn2": p2,
		},
		scaler.Resources{
			{
				"rgn1",
				"rn1",
				"apn1",
				scaler.Rules{
					{"* * * * 2020 *", "pn1"},
					{"* * * * 2021 *", "pn2"},
				},
			},
			{
				"rgn2",
				"rn2",
				"apn2",
				scaler.Rules{
					{"* * * * 2020 *", "pn2"},
					{"* * * * 2021 *", "pn1"},
				},
			},
		},
	}
	logger := log.New(ioutil.Discard, "", log.LstdFlags)

	tests := []struct {
		inTime   time.Time
		call1    poolUpdaterUpdateCall
		call2    poolUpdaterUpdateCall
		outError bool
	}{
		{
			time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			poolUpdaterUpdateCall{nil, "rgn1", "rn1", "apn1", p1},
			poolUpdaterUpdateCall{nil, "rgn2", "rn2", "apn2", p2},
			false,
		},
		{
			time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
			poolUpdaterUpdateCall{nil, "rgn1", "rn1", "apn1", p2},
			poolUpdaterUpdateCall{errors.New("error"), "rgn2", "rn2", "apn2", p1},
			true,
		},
	}

	for _, tt := range tests {
		(func() {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			updater := mock_scaler.NewMockPoolUpdater(ctrl)
			gomock.InOrder(
				updater.
					EXPECT().
					UpdatePool(
						gomock.Any(),
						gomock.Eq(tt.call1.outResourceGroupName),
						gomock.Eq(tt.call1.outResourceName),
						gomock.Eq(tt.call1.outAgentPoolName),
						gomock.Eq(tt.call1.outParameters),
					).
					Return(tt.call1.inError),
				updater.
					EXPECT().
					UpdatePool(
						gomock.Any(),
						gomock.Eq(tt.call2.outResourceGroupName),
						gomock.Eq(tt.call2.outResourceName),
						gomock.Eq(tt.call2.outAgentPoolName),
						gomock.Eq(tt.call2.outParameters),
					).
					Return(tt.call2.inError),
			)

			scaler := &scaler.Scaler{
				Logger:        logger,
				PoolUpdater:   updater,
				Configuration: configuration,
			}
			err := scaler.Scale(tt.inTime)
			if (err != nil) != tt.outError {
				t.Errorf("got %v, want %v", err, tt.outError)
			}
		})()
	}
}

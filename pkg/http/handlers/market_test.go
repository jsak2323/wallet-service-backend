package handlers

import (
	"testing"

	"github.com/btcid/wallet-services-backend-go/pkg/domain/market"
	"github.com/golang/mock/gomock"
)

func TestMarketService_ConvertCoinToIdr(t *testing.T) {
	ctrl := gomock.NewController(t)

	marketMock := market.NewMockRepository(ctrl)
	
	type fields struct {
		marketRepo market.Repository
		LastPrices map[string]string
	}
	type args struct {
		amount string
		symbol string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		mock 	   func()
		wantResult string
		wantErr    bool
	}{
		{
			name: "ok",
			fields: fields{ marketRepo: marketMock },
			args: args{ amount: "1.2", symbol: "btc" },
			mock: func() {
				marketMock.EXPECT().LastPriceBySymbol("btc","idr").Return("10", nil)
			},
			wantResult: "12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MarketService{
				marketRepo: tt.fields.marketRepo,
				LastPrices: tt.fields.LastPrices,
			}

			tt.mock()

			gotResult, err := s.ConvertCoinToIdr(tt.args.amount, tt.args.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarketService.ConvertCoinToIdr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("MarketService.ConvertCoinToIdr() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

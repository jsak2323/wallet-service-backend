package util

import (
	"testing"
)

func TestRawToCoin(t *testing.T) {
	type args struct {
		value   string
		decimal int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{value: "2000", decimal: 8},
			want: "0.00002000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RawToCoin(tt.args.value, tt.args.decimal)
			if (err != nil) != tt.wantErr {
				t.Errorf("RawToCoin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RawToCoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatCurrency(t *testing.T) {
	type args struct {
		value  string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			args: args{value: "12345678"},
			wantResult: "12,345,678 IDR",
		},
		{
			args: args{value: "1234567"},
			wantResult: "1,234,567 IDR",
		},
		{
			args: args{value: "123456"},
			wantResult: "123,456 IDR",
		},
		{
			args: args{value: "1234"},
			wantResult: "1,234 BTC",
		},
		{
			args: args{value: "1234.222"},
			wantResult: "1,234.222 BTC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := FormatCurrency(tt.args.value); gotResult != tt.wantResult {
				t.Errorf("FormatCurrency() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

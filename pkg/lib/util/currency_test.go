package util

import (
	"testing"
)

func TestCoinToRaw(t *testing.T) {
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
			args: args{value: "0.00002000", decimal: 8},
			want: "2000",
		},
		{
			name: "ok",
			args: args{value: "96092252354.64214000", decimal: 8},
			want: "9609225235464214000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CoinToRaw(tt.args.value, tt.args.decimal)
			if got != tt.want {
				t.Errorf("CoinToRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			want: "0.00002",
		},
		{
			name: "ok",
			args: args{value: "9609225235464214001", decimal: 8},
			want: "96092252354.64214001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RawToCoin(tt.args.value, tt.args.decimal)
			if got != tt.want {
				t.Errorf("RawToCoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatCurrency(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			args:       args{value: "12345678"},
			wantResult: "12,345,678",
		},
		{
			args:       args{value: "1234567"},
			wantResult: "1,234,567",
		},
		{
			args:       args{value: "123456"},
			wantResult: "123,456",
		},
		{
			args:       args{value: "1234"},
			wantResult: "1,234",
		},
		{
			args:       args{value: "1234.222"},
			wantResult: "1,234.222",
		},
		{
			args:       args{value: "96092252354.64214000"},
			wantResult: "96,092,252,354.64214000",
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

func TestAddCurency(t *testing.T) {
	type args struct {
		value1 string
		value2 string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{value1: "100000000000", value2: "0.0001"},
			want: "100000000000.0001",
		},
		{
			name: "ok",
			args: args{value1: "0.00002000", value2: "0.00002000"},
			want: "0.00004",
		},
		{
			name: "ok",
			args: args{value1: "21", value2: "2"},
			want: "23",
		},
		{
			name: "ok",
			args: args{value1: "0.00002001", value2: "0.00002"},
			want: "0.00004001",
		},
		{
			name: "ok",
			args: args{value1: "0.10002001", value2: "0.00002"},
			want: "0.10004001",
		},
		{
			name: "ok",
			args: args{value1: "96092252354.642140011", value2: "96092252354.64214000"},
			want: "192184504709.28428001",
		},
		{
			name: "ok",
			args: args{value1: "96092252354.64214000", value2: "96092252354.64214000"},
			want: "192184504709.28428",
		},
		{
			name: "ok",
			args: args{value1: "1111111111111111111111111111111111111111111111", value2: "1111111111111111111111111111111111111111111111"},
			want: "2222222222222222222222222222222222222222222222",
		},
		{
			name: "ok",
			args: args{value1: "1111111111111111111111111111111111111111111111.1", value2: "1111111111111111111111111111111111111111111111"},
			want: "2222222222222222222222222222222222222222222222.1",
		},
		{
			name: "ok",
			args: args{value1: "1111111111111111111111111111111111111111111111", value2: "1111111111111111111111111111111111111111111111.1"},
			want: "2222222222222222222222222222222222222222222222.1",
		},
		{
			name: "ok",
			args: args{value1: "1111111111111111111111111111111111111111111111.11111111111", value2: "1111111111111111111111111111111111111111111111.11111111111"},
			want: "2222222222222222222222222222222222222222222222.22222222",
		},
		{
			name: "ok",
			args: args{value1: "09609225235464214000000", value2: "9609225235464214000000"},
			want: "19218450470928428000000",
		},
		{
			name: "ok",
			args: args{value1: "9609225235464214000", value2: "9609225235464214000"},
			want: "19218450470928428000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddCurrency(tt.args.value1, tt.args.value2)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCurrency() Error = %v", err)
			}
			if got != tt.want {
				t.Errorf("AddCurrency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubCurrency(t *testing.T) {
	type args struct {
		value1 string
		value2 string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{value1: "200000001", value2: "200000000"},
			want: "1",
		},
		{
			name: "ok",
			args: args{value1: "09609225235464214000001", value2: "9609225235464214000000"},
			want: "1",
		},
		{
			name: "ok",
			args: args{value1: "9609225235464214001", value2: "9609225235464214000"},
			want: "1",
		},
		{
			name: "ok",
			args: args{value1: "100000000000", value2: "0.0001"},
			want: "99999999999.9999",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SubCurrency(tt.args.value1, tt.args.value2)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubCurrency() Error = %v", err)
			}
			if got != tt.want {
				t.Errorf("SubCurrency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPercentCurrency(t *testing.T) {
	type args struct {
		value1 string
		value2 string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{value1: "100", value2: "3"},
			want: "3333.33",
		},
		{
			name: "ok",
			args: args{value1: "9609225235464214001", value2: "9609225235464214001"},
			want: "100.00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PercentCurrency(tt.args.value1, tt.args.value2)
			if (err != nil) != tt.wantErr {
				t.Errorf("PercentCurrency() Error = %v", err)
			}
			if got != tt.want {
				t.Errorf("PercentCurrency() = %v, want %v", got, tt.want)
			}
		})
	}
}

package fireblocks

import (
	"testing"
	"fmt"
	"time"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
)

func Test_auth(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(time.Now())
			config.CONF.FireblocksServeruser = "fireblocks"
			config.CONF.FireblocksServerhashkey = "sr8|1o6B~li3W_."
			config.CONF.FireblocksServerpass = "ds9FntYP55"
			if got := auth(); got != tt.want {
				t.Errorf("auth() = %v, want %v", got, tt.want)
			}
		})
	}
}

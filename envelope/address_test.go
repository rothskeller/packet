package envelope

import (
	"reflect"
	"testing"
)

func TestParseAddressList(t *testing.T) {
	tests := []struct {
		name      string
		arg       string
		wantAddrs []*Address
		wantErr   bool
	}{
		{
			name: "basic internet address",
			arg:  "steve@rothskeller.net",
			wantAddrs: []*Address{
				{Address: "steve@rothskeller.net"},
			},
		},
		{
			name: "address with full name",
			arg:  "Steve Roth <steve@rothskeller.net>",
			wantAddrs: []*Address{
				{Name: "Steve Roth", Address: "steve@rothskeller.net"},
			},
		},
		{
			name: "bare word address",
			arg:  "kc6rsc",
			wantAddrs: []*Address{
				{Address: "kc6rsc"},
			},
		},
		{
			name: "comma separated list",
			arg:  "alpha, beta",
			wantAddrs: []*Address{
				{Address: "alpha"},
				{Address: "beta"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddrs, err := ParseAddressList(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAddressList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAddrs, tt.wantAddrs) {
				t.Errorf("ParseAddressList() = %v, want %v", gotAddrs, tt.wantAddrs)
			}
		})
	}
}

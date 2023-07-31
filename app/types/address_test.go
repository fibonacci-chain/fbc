package types

import (
	"strings"
	"testing"

	"github.com/fibonacci-chain/fbc/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// test compate 0x address to send and query
func TestAccAddressFromBech32(t *testing.T) {
	config := types.GetConfig()
	SetBech32Prefixes(config)

	//make data
	tests := []struct {
		addrStr    string
		expectPass bool
	}{
		{"fb1e34x06m7hqwt4n98yuaxuxljr0juumktxdav63", true},
		{"0x0073F2E28ef8F117e53d858094086Defaf1837D5", true},
		{"2CF4ea7dF75b513509d95946B43062E26bD88035", true},
		{strings.ToLower("2CF4ea7dF75b513509d95946B43062E26bD88035"), true},
		{strings.ToUpper("2CF4ea7dF75b513509d95946B43062E26bD88035"), true},
		{"fb1e34x06m7hqwt4n98yuaxuxljr0juumktxdav63_", false},
		{"0x0073F2E28ef8F117e53d858094086Defaf1837D5_", false},
		{"0073F2E28ef8F117e53d858094086Defaf1837D5_", false},
	}

	//test run
	for _, tc := range tests {
		addr, err := types.AccAddressFromBech32(tc.addrStr)
		if tc.expectPass {
			require.NotNil(t, addr, "test: %v", tc.addrStr)
			require.Nil(t, err, "test: %v", tc.addrStr)
		} else {
			require.Nil(t, addr, "test: %v", tc.addrStr)
			require.NotNil(t, err, "test: %v", tc.addrStr)
		}
	}

}

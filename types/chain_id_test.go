package types

import (
	"errors"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseChainID(t *testing.T) {
	testCases := []struct {
		name      string
		chainID   string
		parseFunc func(chainID string) (*big.Int, error)
		validFunc func(chainId string) bool
		expError  bool
		expInt    *big.Int
	}{
		{
			"valid chain-id, single digit", "ethermint_1-1", nil, nil, false, big.NewInt(1),
		},
		{
			"valid chain-id, multiple digits", "aragonchain_256-1", nil, nil, false, big.NewInt(256),
		},
		{
			"invalid chain-id, double dash", "aragonchain-1-1", nil, nil, true, nil,
		},
		{
			"invalid chain-id, double underscore", "aragonchain_1_1", nil, nil, true, nil,
		},
		{
			"invalid chain-id, dash only", "-", nil, nil, true, nil,
		},
		{
			"invalid chain-id, undefined identifier and EIP155", "-1", nil, nil, true, nil,
		},
		{
			"invalid chain-id, undefined identifier", "_1-1", nil, nil, true, nil,
		},
		{
			"invalid chain-id, uppercases", "ETHERMINT_1-1", nil, nil, true, nil,
		},
		{
			"invalid chain-id, mixed cases", "Ethermint_1-1", nil, nil, true, nil,
		},
		{
			"invalid chain-id, special chars", "$&*#!_1-1", nil, nil, true, nil,
		},
		{
			"invalid eip155 chain-id, cannot start with 0", "ethermint_001-1", nil, nil, true, nil,
		},
		{
			"invalid eip155 chain-id, cannot invalid base", "ethermint_0x212-1", nil, nil, true, nil,
		},
		{
			"invalid eip155 chain-id, non-integer", "ethermint_ethermint_9000-1", nil, nil, true, nil,
		},
		{
			"invalid epoch, undefined", "ethermint_-", nil, nil, true, nil,
		},
		{
			"blank chain ID", " ", nil, nil, true, nil,
		},
		{
			"empty chain ID", "", nil, nil, true, nil,
		},
		{
			"empty content for chain id, eip155 and epoch numbers", "_-", nil, nil, true, nil,
		},
		{
			"long chain-id", "ethermint_" + strings.Repeat("1", 40) + "-1", nil, nil, true, nil,
		},
		{
			"overwritten function",
			"ethermint",
			func(chainID string) (*big.Int, error) { return big.NewInt(90001), nil },
			func(chainId string) bool { return !(len(chainId) > 48) },
			false,
			big.NewInt(90001),
		},
		{
			"overwritten function, too long",
			"ethermint-test1-test2-test3-test4-test5-test6-test7",
			func(chainID string) (*big.Int, error) {
				chainID = strings.TrimSpace(chainID)
				if len(chainID) > 48 {
					return nil, errors.New("invalid chain id")
				}
				return big.NewInt(90001), nil
			},
			func(chainId string) bool { return !(len(chainId) > 48) },
			true,
			nil,
		},
	}

	for _, tc := range testCases {
		if tc.parseFunc != nil {
			SetParseChainIDFunc(tc.parseFunc)
		}
		if tc.validFunc != nil {
			SetValidChainIDFunc(tc.validFunc)
		}
		chainIDEpoch, err := ParseChainID(tc.chainID)
		if tc.expError {
			require.Error(t, err, tc.name)
			require.Nil(t, chainIDEpoch)

			require.False(t, IsValidChainID(tc.chainID), tc.name)
		} else {
			require.NoError(t, err, tc.name)
			require.Equal(t, tc.expInt, chainIDEpoch, tc.name)
			require.True(t, IsValidChainID(tc.chainID))
		}
	}
}

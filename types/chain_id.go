package types

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	regexChainID         = `[a-z]{1,}`
	regexEIP155Separator = `_{1}`
	regexEIP155          = `[1-9][0-9]*`
	regexEpochSeparator  = `-{1}`
	regexEpoch           = `[1-9][0-9]*`
	ethermintChainID     = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)%s(%s)$`, regexChainID, regexEIP155Separator, regexEIP155, regexEpochSeparator, regexEpoch))
)

// DefaultValidChainID returns the default valid function
func DefaultValidChainID(chainID string) bool {
	if len(chainID) > 48 {
		return false
	}

	return ethermintChainID.MatchString(chainID)
}

// validChainIDFunc returns the current function and can be overwritten for custom valid function
var validChainIDFunc = DefaultValidChainID

// SetValidChainIDFunc allows for chain-id's custom validation by overriding the valid
// function used for chain-id validation.
func SetValidChainIDFunc(fn func(chainId string) bool) {
	validChainIDFunc = fn
}

// IsValidChainID returns false if the given chain identifier is incorrectly formatted.
func IsValidChainID(chainID string) bool {
	return validChainIDFunc(chainID)
}

// DefaultParseChainID returns the default parse function
func DefaultParseChainID(chainID string) (*big.Int, error) {
	chainID = strings.TrimSpace(chainID)
	if len(chainID) > 48 {
		return nil, sdkerrors.Wrapf(ErrInvalidChainID, "chain-id '%s' cannot exceed 48 chars", chainID)
	}

	matches := ethermintChainID.FindStringSubmatch(chainID)
	if matches == nil || len(matches) != 4 || matches[1] == "" {
		return nil, sdkerrors.Wrapf(ErrInvalidChainID, "%s: %v", chainID, matches)
	}

	// verify that the chain-id entered is a base 10 integer
	chainIDInt, ok := new(big.Int).SetString(matches[2], 10)
	if !ok {
		return nil, sdkerrors.Wrapf(ErrInvalidChainID, "epoch %s must be base-10 integer format", matches[2])
	}

	return chainIDInt, nil
}

// parseChainIDFunc returns the current function and can be overwritten for custom parse function
var parseChainIDFunc = DefaultParseChainID

// SetParseChainIDFunc allows for chain-id's custom parse function by overriding the
// function used for chain-id parse.
func SetParseChainIDFunc(fn func(chainID string) (*big.Int, error)) {
	parseChainIDFunc = fn
}

// ParseChainID parses a string chain identifier's epoch to an Ethereum-compatible
// chain-id in *big.Int format. The function returns an error if the chain-id has an invalid format
func ParseChainID(chainID string) (*big.Int, error) {
	return parseChainIDFunc(chainID)
}

/// tokens.go is responsible for mapping token names to token addresses

package tokens

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type Network string

const (
	Mainnet Network = "mainnet"
	Devnet          = "devnet"
)

// Token is a struct that contains the token name and address
type Token struct {
	Address    string
	ChainId    int
	Decimals   int
	Name       string
	Symbol     string
	LogoURI    string
	Tags       []string
	Extensions map[string]string
}
type TokenCache map[string]Token

var tokenCache TokenCache

func memoize(f func(string, Network) (Token, error)) func(string, Network) Token {

	return func(address string, network Network) Token {
		if _, ok := tokenCache[address]; !ok {
			token, err := f(address, network)
			if err != nil {
				return Token{}
			}
			tokenCache[address] = token
		}
		return tokenCache[address]
	}
}

// / findTokenInJupagents finds the token in the jupagents api
func findTokenInJupagents(address string, network Network) (Token, error) {
	var token Token
	resp, err := http.Get("https://token.jup.ag/strict")
	if err != nil {
		return token, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return token, err
	}
	var tokens []Token
	if err := json.Unmarshal(body, &tokens); err != nil {
		return token, err
	}

	for _, token := range tokens {
		if token.Address == address {
			return token, nil
		}
	}
	return token, errors.New("token not found")
}

// findTokenSymbolInJupagents finds the token in the jupagents api based on the
// token symbol
func findTokenSymbolInJupagents(symbol string, network Network) (Token, error) {
	var token Token
	resp, err := http.Get("https://token.jup.ag/strict")
	if err != nil {
		return token, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return token, err
	}
	var tokens []Token
	if err := json.Unmarshal(body, &tokens); err != nil {
		return token, err
	}

	for _, token := range tokens {
		if strings.ToLower(symbol) == strings.ToLower(token.Symbol) {
			return token, nil
		}
	}
	return token, errors.New("token not found")
}

func GetTokenFromSymbol(symbol string, network Network) (Token, error) {
	tokenCache = make(TokenCache)
	// setup redis connection
	if network == Mainnet {
		token := memoize(findTokenSymbolInJupagents)(symbol, network)
		if token.Address == "" {
			return token, errors.New("token not found")
		}
		return token, nil
	}

	if network == Devnet {
		return Token{
			Address:    "sandphoQsRiNd85VgRrdSXdhS56d58Xa9iDKwdnKfWR",
			ChainId:    1,
			Decimals:   6,
			Name:       "sand",
			Symbol:     "SAND",
			LogoURI:    "https://jup.io/images/jup-logo.png",
			Tags:       []string{"sand", "jupiter", "jupiter token"},
			Extensions: map[string]string{},
		}, nil
	}

	return Token{}, errors.New("network not supported")

}

// Tokens is a map of token names to token addresses
func GetTokenFromAddress(address string, network Network) (Token, error) {
	tokenCache = make(TokenCache)
	// setup redis connection
	if network == Mainnet {
		token := memoize(findTokenInJupagents)(address, network)
		if token.Address == "" {
			return token, errors.New("token not found")
		}
		return token, nil
	}

	if network == Devnet {
		return Token{
			Address:    "0x000000",
			ChainId:    0,
			Decimals:   6,
			Name:       "Jupiter",
			Symbol:     "JUP",
			LogoURI:    "https://jup.io/images/jup-logo.png",
			Tags:       []string{"jup", "jupiter", "jupiter token"},
			Extensions: map[string]string{},
		}, nil
	}

	return Token{}, errors.New("network not supported")

}

func IsValidAccount(ctx context.Context, address, rpcUrl string) bool {
	cluster := rpc.New(rpcUrl)
	pk, err := solana.PublicKeyFromBase58(address)
	if err != nil {
		return false
	}
	if _, err = cluster.GetAccountInfo(ctx, pk); err != nil {
		return false
	}

	return true
}

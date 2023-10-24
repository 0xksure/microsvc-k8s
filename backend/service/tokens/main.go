/// tokens.go is responsible for mapping token names to token addresses

package tokens

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Network string

const (
	Mainnet Network = "mainnet"
	Devnet          = "devnet"
)

// Token is a struct that contains the token name and address
type Token struct {
	address    string
	chainId    int
	decimals   int
	name       string
	symbol     string
	logoURI    string
	tags       []string
	extensions map[string]string
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
		if token.address == address {
			return token, nil
		}
	}
	return token, errors.New("token not found")
}

// Tokens is a map of token names to token addresses
func getTokenFromAddress(address string, network Network) (Token, error) {
	tokenCache = make(TokenCache)
	// setup redis connection
	if network == Mainnet {
		token := memoize(findTokenInJupagents)(address, network)
		if token.address == "" {
			return token, errors.New("token not found")
		}
		return token, nil
	}

	if network == Devnet {
		return Token{
			address:    "0x000000",
			chainId:    0,
			decimals:   0,
			name:       "Jupiter",
			symbol:     "JUP",
			logoURI:    "https://jup.io/images/jup-logo.png",
			tags:       []string{"jup", "jupiter", "jupiter token"},
			extensions: map[string]string{},
		}, nil
	}

	return Token{}, errors.New("network not supported")

}

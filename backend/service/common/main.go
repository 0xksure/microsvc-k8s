package common

import (
	"context"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/err/identity"
	"github.com/err/tokens"
	"github.com/gagliardetto/solana-go"
	"github.com/google/go-github/v55/github"
	"github.com/pkg/errors"
)

// CreateSigningLink creates a signing link which is used to redirect the user to
// a signing page. The arguments are used to link the signing page to the bounty
func CreateSigningLink(bountytId, installationId int64, tokenAddress, bountyUIAmount, creatorAddress, issueUrl, organization, team, domainType string) string {
	signingUrl := os.Getenv("SIGNING_URL")
	return fmt.Sprintf("%s/bounty?bountyId=%d&tokenAddress=%s&bountyUIAmount=%s&creatorAddress=%s&installationId=%d&referrer=%s&platform=%s&organization=%s&team=%s&domainType=%s", signingUrl, bountytId, tokenAddress, bountyUIAmount, creatorAddress, installationId, issueUrl, "github", organization, team, domainType)
}

// FindAtUsers finds all the users in a string
func FindAtUsers(text string) []string {
	r := regexp.MustCompile(`@(\w+)`)
	usernamesWithAt := r.FindAllString(text, -1)
	var usernames []string
	for _, userName := range usernamesWithAt {
		usernames = append(usernames, userName[1:])
	}
	return usernames
}

// FindUserIdsFromNames finds the user ids from the user names
func FindUserIdsFromNames(names []string, client *github.Client) ([]uint64, error) {
	var userIds []uint64
	errors := []error{}
	for _, name := range names {
		user, _, err := client.Users.Get(context.Background(), name)
		if err != nil {
			errors = append(errors, err)
		}
		userIds = append(userIds, uint64(user.GetID()))
	}
	if len(errors) > 0 {
		// combine errors
		var errStr string
		for _, err := range errors {
			errStr += err.Error()
		}
		return userIds, fmt.Errorf(errStr)
	}
	return userIds, nil
}

func FindIdentitiesFromAtName(text, rpcUrl string, client *github.Client) (identity.Identities, error) {
	usernames := FindAtUsers(text)
	if len(usernames) == 0 {
		return nil, errors.New("No usernames found in text")
	}
	userIds, err := FindUserIdsFromNames(usernames, client)
	if err != nil {
		return nil, err
	}
	if len(userIds) == 0 {
		return nil, errors.Errorf("No user ids was retrieved from usernames %v found in text", usernames)
	}
	identities, err := identity.GetIdentities(rpcUrl, "github", userIds)
	if err != nil {
		return nil, err
	}
	if len(identities) == 0 {
		return nil, errors.Errorf("No identities was retrieved from user ids %v found in text", userIds)
	}
	return identities, nil
}

func FindWalletsFromAtName(text, rpcUrl string, client *github.Client) ([]solana.PublicKey, error) {
	identities, err := FindIdentitiesFromAtName(text, rpcUrl, client)
	if err != nil {
		return nil, err
	}
	var wallets []solana.PublicKey
	for _, identity := range identities {
		wallets = append(wallets, identity.Address)
	}
	return wallets, nil
}

type ParsedGithubBounty struct {
	AmountUI    string
	Amount      uint64
	TokenSymbol string
	Token       tokens.Token
}

// ParseBountyMessage parses the bounty message from the issue text
func ParseBountyMessage(text string, network tokens.Network) (ParsedGithubBounty, error) {
	var parsedGithubBounty ParsedGithubBounty
	r := regexp.MustCompile(`\$(\w+:\d+)\$`)
	bounty := r.FindString(text)
	if bounty == "" {
		return parsedGithubBounty, errors.New("No bounty found in issueText")
	}
	bountyParts := strings.Split(strings.Trim(bounty, "$"), ":")
	if len(bountyParts) != 2 {
		return parsedGithubBounty, errors.Errorf("Expected bounty to be two values. Got %v", bounty)
	}
	// token is a string literal e.g. USDC
	tokenSymbol := bountyParts[0]
	// assume amount is in decimals e.g. 100.00 so prettyAmount
	amountProto := bountyParts[1]
	amountUi, err := strconv.ParseFloat(amountProto, 64)
	if err != nil {
		return parsedGithubBounty, err
	}

	token, err := tokens.GetTokenFromSymbol(tokenSymbol, network)
	if err != nil {
		return parsedGithubBounty, err
	}

	amount := uint64(math.Floor(amountUi * (float64(10 ^ token.Decimals))))

	return ParsedGithubBounty{
		AmountUI:    amountProto,
		Amount:      amount,
		TokenSymbol: tokenSymbol,
		Token:       token,
	}, nil
}

type SandblizzardError struct{}

func GetExplorerLink(network tokens.Network, signature solana.Signature) string {
	if network == tokens.Devnet {
		return fmt.Sprintf("https://explorer.solana.com/tx/%s?cluster=devnet", signature)
	} else if network == tokens.Testnet {
		return fmt.Sprintf("https://explorer.solana.com/tx/%s?cluster=testnet", signature)
	} else if network == tokens.Localnet {
		return fmt.Sprintf("https://explorer.solana.com/tx/%s?cluster=localnet", signature)
	}
	return fmt.Sprintf("https://explorer.solana.com/tx/%s?cluster=mainnet", signature)
}

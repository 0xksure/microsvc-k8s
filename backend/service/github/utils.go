package github_bounty

import (
	"fmt"
	"os"
)

// CreateSigningLink creates a signing link which is used to redirect the user to
// a signing page. The arguments are used to link the signing page to the bounty
func CreateSigningLink(bountytId, installationId int64, tokenAddress, bountyUIAmount, creatorAddress, issueUrl, organization, team, domainType string) string {
	signingUrl := os.Getenv("SIGNING_URL")
	return fmt.Sprintf("%s/bounty?bountyId=%d&tokenAddress=%s&bountyUIAmount=%s&creatorAddress=%s&installationId=%d&referrer=%s&platform=%s&organization=%s&team=%s&domainType=%s", signingUrl, bountytId, tokenAddress, bountyUIAmount, creatorAddress, installationId, issueUrl, "github", organization, team, domainType)
}

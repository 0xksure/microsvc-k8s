package kafka

import "encoding/json"

type LinkerMessage struct {
	Username      string
	UserId        int
	WalletAddress string
}

func (l *LinkerMessage) Decode(message KafkaMessage) error {
	return json.Unmarshal(message.Msg, &l)
}

type BountyMessage struct {
	BountySignStatus string
	BountyId         int
	BountyUIAmount   string
	TokenAddress     string
	CreatorAddress   string
}

// decode decodes the kafka message into a bounty message
// todo: validate the bounty message
func (l *BountyMessage) Decode(message KafkaMessage) error {
	return json.Unmarshal(message.Msg, &l)
}

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

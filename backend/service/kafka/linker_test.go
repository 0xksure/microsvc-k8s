package kafka

import (
	"encoding/json"
	"log"
	"testing"
)

func TestLinker(t *testing.T) {
	jsonData := "{\"username\":\"0xksure\",\"userId\":47750504,\"walletAddress\":\"HGk4BG9mQHWnX6GdUjS2gdskXUhQx1VzuHqnbdkKuwhP\"}"
	data := []byte(jsonData)
	t.Run("TestLinker", func(t *testing.T) {
		t.Parallel()

		var linkerMessage LinkerMessage
		kafkaMessage := KafkaMessage{Msg: data}
		err := linkerMessage.Decode(kafkaMessage)
		if err != nil {
			log.Fatal(err)
		}
		if linkerMessage.Username != "0xksure" {
			t.Errorf("username not decoded: %s", linkerMessage.Username)
		}
		if linkerMessage.UserId != 47750504 {
			t.Errorf("userId not decoded: %d", linkerMessage.UserId)
		}
		if linkerMessage.WalletAddress != "HGk4BG9mQHWnX6GdUjS2gdskXUhQx1VzuHqnbdkKuwhP" {
			t.Errorf("walletAddress not decoded: %s", linkerMessage.WalletAddress)
		}
	})
	t.Run("TestGob", func(t *testing.T) {
		var linkerMessage LinkerMessage
		err := json.Unmarshal(data, &linkerMessage)
		if err != nil {
			log.Fatal(err)
		}
		if linkerMessage.Username != "0xksure" {
			t.Errorf("username not decoded: %s", linkerMessage.Username)
		}
		if linkerMessage.UserId != 47750504 {
			t.Errorf("userId not decoded: %d", linkerMessage.UserId)
		}
		if linkerMessage.WalletAddress != "HGk4BG9mQHWnX6GdUjS2gdskXUhQx1VzuHqnbdkKuwhP" {
			t.Errorf("walletAddress not decoded: %s", linkerMessage.WalletAddress)
		}
	})
}

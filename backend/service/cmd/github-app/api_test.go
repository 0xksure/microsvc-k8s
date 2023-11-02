package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/err/db"
	"github.com/err/kafka"
	"github.com/err/protoc/bounty"
	"github.com/golang/protobuf/proto"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog"
)

type MExternal interface{}

type MKafkaClient struct{}

func (m MKafkaClient) GenerateKafkaConsumer(ctx context.Context, topic string, kf chan kafka.KafkaMessage) {
}

type MBountyOrm struct{}

func (m MBountyOrm) Close() {}

func (m MBountyOrm) CreateBountyCreator(entityId int, username, entityType string) error {
	return nil
}
func (m MBountyOrm) CreateBounty(ctx context.Context, bountyInput db.BountyInput, entityName string) (int, error) {
	return 0, nil
}
func (m MBountyOrm) GetBountyOnIssueId(ctx context.Context, issueId int) (db.Bounty, error) {
	return db.Bounty{
		BountyInput: db.BountyInput{
			Id:          0,
			EntityId:    0,
			Url:         "",
			IssueId:     issueId,
			IssueNumber: 0,
			RepoId:      0,
			RepoName:    "",
			RepoOwner:   "",
			OwnerId:     0,
			Status:      "",
		},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}, nil
}
func (m MBountyOrm) UpdateBountyStatus(ctx context.Context, issueId int, status string) error {
	return nil
}

func TestApi(t *testing.T) {
	t.Log("Testing api")

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log("New request: ", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 1}`))
	}))
	defer func() { testServer.Close() }()
	apiUrl, err := url.Parse(testServer.URL + "/")
	if err != nil {
		panic(err)
	}

	githubAppConfig := githubapp.Config{
		WebURL:   apiUrl.String(),
		V3APIURL: apiUrl.String(),
		V4APIURL: apiUrl.String(),
	}
	githubAppConfig.App.IntegrationID = 4444
	githubAppConfig.App.WebhookSecret = "foobar"
	githubAppConfig.App.PrivateKey = "-----BEGIN PRIVATE KEY-----\n" +
		"MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDKlvdvjS7EEkYm" +
		"RcMYBemMNvIbMRoJPBuwtRHph8pehi/BUXmTUsGUSDelFAl6tH2eMpHD6FpMJkiu" +
		"PAm875yJYo9nxgnH8SNzS1FSR05qoFkJGSvkQwLd++V4Z17kqg9LyzQf/XGSiGGm" +
		"VbvPS8rkJUdrFxgTjqhw8KIfNXGG4zLNcB6i0YgW2lZFgaFB16J3/O4KDm21bgdi" +
		"qI6F5xKel/YSQl+AfF2NlFt5xI3wuZ4YhD7wrp3WWgiS/JIxu0ESDmJTYExdb3K0" +
		"Dh6bmQl+wYJbkBpO2csNRDehtM2YEU83N8mmlSbIbSSinZYT3KwT9JQbd8QqJZCP" +
		"/mOPFacrAgMBAAECggEAMqGWR3fWd0RF6ezHfGqF2vgke+1Cn4o5NWmbh2zbg9Iv" +
		"fzYYl1w4axG9bnFaiSMwvefPjFG2t49d3MW+fUy5J5DNXFcfPKwkev0Y3uJZU8at" +
		"WdvDn3Gr9sSsrfHPwoBKAFxRs6kIyGFzXjnRDVbY5zn15mrIJqMhr9BEBF6798TA" +
		"qTG0FYTqkGK0D+FVfaWXvQ0u9Jw0KootS4kKHNwDmbZK2xYI2Ilt1ikeN9MMt8ZD" +
		"tXV4shTnQaYPty5Atr9Dzh052FTnlwsclVo33XHF2N2dfe7TaJYTaf5uXuh7Vkj3" +
		"99bKuvA8iLmEVe+i86L9K9LD5QEywO/sNcqQkNU8sQKBgQDkycMSv1dfSAgSZywU" +
		"hcJqFWJ3TpAP5aoWqw8Svill3Qs6/2zatU4XC4tzRX9KW437M4ORNee/XrCi9z4L" +
		"VOeOtR+gp/zn9DebBzaEdTfMlof+znPfx9fVIpBkpjezDgiFNkeeMnKTQrT/kULl" +
		"Zso29pCfgO/57L9Vi3pjiqKt0wKBgQDir4EDVesZTGcBstSRLIGrUQFmvxICPaMm" +
		"0PogvQhUv8UFMxx6nBl/ZB7AZWQH0TRlpVL7iqS2clO1dQOgrtgwgUV4m3Ml9Ivr" +
		"vI2fgCWzFudPst3oSu9Udc9Pq995Pl2nnwV612p3P6p3wgCaFjbouyyJZxNmHQj/" +
		"JPtv65nSSQKBgQDVcRnJorLLlHLbYF9yYfunZn3vWl7yRcvxy/KLBNewTZENoIAY" +
		"Zm8M9ttJVjvTziheg4ep8EVdduSJlOnQPoysyXNROYeril5aBlepKYY+Gu2THV5j" +
		"FpjYIZ/eFmf+Zwgx5xrXjq7vjZs4lnd3dvcOYec4t1yqqGE0WKR8uzjbuwKBgAYL" +
		"0FEadY7TLtwovOqyWTMMkhD/f6d3pWZfpIxC/nnkM4kT9+p9R2DSds+C5MwglFkx" +
		"s6jp5cLIAduRJ2udvj5s9EFnRAb7ItBC0zQx4s+ICNtjVe/gL8n86m6hkvBU7YKP" +
		"B0JjhH9xv0Y6cnGprgU/GM0BZs8ObzL+9YXirtOhAoGACk1TGT6I5oMHLDB4yOOe" +
		"a7dT+kAY6a2W2Sq2e0VN70EkVmXGC6ODpRzPcH7nojcMN5jk8QHHosWNK8DECwAj" +
		"uCGYn8G/0yAlhkddzE1+y1f5nVm+GCQTrdMMuqOwJnosifdoNDbWg4oGiBRt1uwI" +
		"aoxbWonlDAZaLC+8Bxe1Hss=" +
		"\n-----END PRIVATE KEY-----"

	clientCreator := githubapp.NewClientCreator(apiUrl.String(), apiUrl.String(), githubAppConfig.App.IntegrationID, []byte(githubAppConfig.App.PrivateKey))
	t.Logf("Testing api %v", githubAppConfig)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	kafkaClient := kafka.BountyKafkaClient{}
	bountyOrm := MBountyOrm{}
	api := API{
		serverAddr:    "localhost",
		port:          "8080",
		logger:        logger,
		kafkaClient:   kafkaClient,
		clientCreator: clientCreator,
		ghConfig:      &githubAppConfig,
		rpcUrl:        apiUrl.String(),
		network:       "devnet",
		bountyOrm:     bountyOrm,
	}

	t.Run("Test bountyKafka Handler", func(t *testing.T) {
		t.Log("Testing bountyKafka Handler")
		bountyMessage := bounty.BountyMessage{
			BountySignStatus: 1,
			Bountyid:         1,
			BountyUIAmount:   "1000",
			TokenAddress:     "0x000",
			CreatorAddress:   "0x000",
			InstallationId:   4444,
			Platform:         "github",
			Organization:     "sandblizzard",
			Team:             "sandblizzard",
			DomainType:       "issues",
		}
		bountyMessageBytes, err := proto.Marshal(&bountyMessage)
		if err != nil {
			t.Error(err)
		}
		kafkaMessage := kafka.KafkaMessage{
			Topic: "bounty",
			Key:   "bounty",
			Msg:   bountyMessageBytes,
		}
		if err := api.BountyKafkaHandler(&kafkaMessage); err != nil {
			t.Error(err)
		}
	})

}

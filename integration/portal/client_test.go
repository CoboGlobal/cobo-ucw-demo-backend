package portal

import (
	"context"
	"testing"

	"cobo-ucw-backend/internal/conf"

	"github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/go-kratos/kratos/v2/log"
)

func TestNewClient(t *testing.T) {
	cobo_waas2.Transaction{}
	t.Run("enable_tokens", func(t *testing.T) {
		client := NewClient(&conf.UCW{
			ProjectId: "",
			Threshold: 0,
			CoboPortal: &conf.CoboPortal{
				Apikey: "dca8011baccc336f600c94ec3f0ae73c64a3336286d232b4551b6baad669fccb",
				Env:    conf.CoboPortal_SANDBOX,
				Debug:  false,
			}}, log.DefaultLogger)

		res, _, err := client.WalletsAPI.ListEnabledTokens(client.WithContext(context.Background())).Execute()
		t.Error(err)
		t.Log(res)
	})

	t.Run("supported_tokens", func(t *testing.T) {
		client := NewClient(&conf.UCW{
			ProjectId: "",
			Threshold: 0,
			CoboPortal: &conf.CoboPortal{
				Apikey: "dca8011baccc336f600c94ec3f0ae73c64a3336286d232b4551b6baad669fccb",
				Env:    conf.CoboPortal_SANDBOX,
				Debug:  false,
			}}, log.DefaultLogger)

		res, _, err := client.WalletsAPI.ListSupportedTokens(client.WithContext(context.Background())).Execute()
		t.Error(err)
		t.Log(res.GetData())
		for _, each := range res.GetData() {

			if each.GetAssetId() == "SETH" || each.GetAssetId() == "XTN" {
				t.Log(each)
			}
		}
	})

	t.Run("get token", func(t *testing.T) {
		client := NewClient(&conf.UCW{
			ProjectId: "",
			Threshold: 0,
			CoboPortal: &conf.CoboPortal{
				Apikey: "dca8011baccc336f600c94ec3f0ae73c64a3336286d232b4551b6baad669fccb",
				Env:    conf.CoboPortal_SANDBOX,
				Debug:  false,
			}}, log.DefaultLogger)

		res, _, err := client.WalletsAPI.GetTokenById(client.WithContext(context.Background()), "XTN").Execute()
		t.Error(err)
		t.Log(res)
	})

	t.Run("list wallet tokens", func(t *testing.T) {
		client := NewClient(&conf.UCW{
			ProjectId: "",
			Threshold: 0,
			CoboPortal: &conf.CoboPortal{
				Apikey: "dca8011baccc336f600c94ec3f0ae73c64a3336286d232b4551b6baad669fccb",
				Env:    conf.CoboPortal_SANDBOX,
				Debug:  false,
			}}, log.DefaultLogger)

		res, _, err := client.WalletsAPI.ListTokenBalancesForWallet(client.WithContext(context.Background()), "34518c3f-ebbd-4406-8208-1bd5c08180a4").Limit(2).Execute()
		t.Error(err)
		t.Log(res)
	})
}

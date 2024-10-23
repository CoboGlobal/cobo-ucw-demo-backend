package portal

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"sync"
	"sync/atomic"

	"cobo-ucw-backend/internal/conf"

	CoboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2/crypto"
	"github.com/go-kratos/kratos/v2/log"
)

type Client struct {
	*CoboWaas2.APIClient
	config *conf.CoboPortal
	signer crypto.ApiSigner

	tokenCache sync.Map

	coboNodeID atomic.Value
}

type logProvider struct {
	logger *log.Helper
}

func (p *logProvider) Printf(ctx context.Context, format string, v ...any) {
	p.logger.WithContext(ctx).Debugf(format, v...)
}

func NewClient(ucw *conf.UCW, logger log.Logger) *Client {
	c := ucw.CoboPortal
	configuration := CoboWaas2.NewConfiguration()
	configuration.Log = &logProvider{
		logger: log.NewHelper(log.With(logger, "mgr.name", "portal")),
	}
	configuration.Servers = append(configuration.Servers, CoboWaas2.ServerConfiguration{
		URL:         "https://api.sandbox.cobo.com/v2",
		Description: "sandbox server",
	})
	configuration.Debug = c.Debug

	if c.Apikey == "" {
		log.Fatalf("cobo portal apikey is null")
	}

	key, err := hex.DecodeString(c.Apikey)
	if err != nil {
		log.Fatalf("cobo portal apikey decode err: %v", err)
	}

	if len(key) != ed25519.SeedSize {
		log.Fatalf("cobo portal apikey len not equal %d", ed25519.SeedSize)
	}

	signer := crypto.Ed25519Signer{
		Secret: c.Apikey,
	}

	client := CoboWaas2.NewAPIClient(configuration)
	return &Client{
		APIClient: client,
		config:    ucw.CoboPortal,
		signer:    signer,
	}
}

func (c *Client) WithContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, CoboWaas2.ContextPortalSigner, c.signer)
	ctx = c.withEnv(ctx)
	return ctx
}

func (c *Client) withEnv(ctx context.Context) context.Context {
	return context.WithValue(ctx, CoboWaas2.ContextEnv, int(c.config.Env))
}

func (c *Client) GetCoboNodeID(ctx context.Context) (string, error) {
	coboNodeID, ok := c.coboNodeID.Load().(string)
	if ok && coboNodeID != "" {
		return coboNodeID, nil
	}
	list, _, err := c.WalletsMPCWalletsAPI.ListCoboKeyHolders(c.WithContext(ctx)).Execute()
	if err != nil {
		return "", err
	}
	for _, each := range list {
		if each.GetStatus() != CoboWaas2.KEYSHAREHOLDERSTATUS_VALID {
			continue
		}
		if each.GetType() != CoboWaas2.KEYSHAREHOLDERTYPE_COBO {
			continue
		}
		c.coboNodeID.Store(each.GetTssNodeId())
		return each.GetTssNodeId(), nil
	}

	return "", nil
}

func (c *Client) GetToken(ctx context.Context, tokenID string) (*CoboWaas2.ExtendedTokenInfo, error) {
	token, ok := c.tokenCache.Load(tokenID)
	if !ok {
		execute, _, err := c.WalletsAPI.GetTokenById(c.WithContext(ctx), tokenID).Execute()
		if err != nil {
			return nil, err
		}
		c.tokenCache.Store(tokenID, execute)
		return execute, nil
	}

	return token.(*CoboWaas2.ExtendedTokenInfo), nil
}

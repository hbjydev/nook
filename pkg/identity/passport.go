package identity

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/bluesky-social/indigo/util"
	redis "github.com/redis/go-redis/v9"
)

const HandleCacheDuration = 5 * time.Minute
const DidDocCacheDuration = 15 * time.Minute

type Passport struct {
	cli *http.Client
	red *redis.Client
	logger *slog.Logger
}

func NewPassport(
	cli *http.Client,
	red *redis.Client,
	logger *slog.Logger,
) *Passport {
	if cli == nil {
		cli = util.RobustHTTPClient()
	}

	return &Passport {
		cli: cli,
		red: red,
		logger: logger,
	}
}

func (p *Passport) ResolveHandle(ctx context.Context, handle string) (string, error) {
	key := handleKey(handle)

	var did string

	res := p.red.Get(ctx, key)
	if res.Err() == nil {
		if err := res.Scan(&did); err != nil {
			return "", err
		}
		return did, nil
	}

	// Cache lookup failed, run query
	did, err := ResolveHandle(ctx, p.cli, handle)
	if err != nil {
		return "", err
	}

	if err := p.red.Set(ctx, key, did, HandleCacheDuration).Err(); err != nil {
		p.logger.Error("failed to save cached handle resolution", "error", err)
	}

	return did, nil
}

func (p *Passport) FetchDidDoc(ctx context.Context, did string) (*DidDoc, error) {
	key := didDocKey(did)

	res := p.red.Get(ctx, key)
	if res.Err() == nil {
		var didDoc DidDoc
		if err := res.Scan(&didDoc); err != nil {
			return nil, err
		}
		return &didDoc, nil
	}

	// Cache lookup failed, run query
	didDoc, err := FetchDidDoc(ctx, p.cli, did)
	if err != nil {
		return nil, err
	}

	if err := p.red.Set(ctx, key, didDoc, DidDocCacheDuration).Err(); err != nil {
		p.logger.Error("failed to save cached did doc", "error", err)
	}

	return didDoc, nil
}

func handleKey(handle string) string {
	return fmt.Sprintf("identity:handle:%s", handle)
}

func didDocKey(did string) string {
	return fmt.Sprintf("identity:did-doc:%s", did)
}

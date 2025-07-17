package identity

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/bluesky-social/indigo/util"
)

func ResolveHandle(ctx context.Context, cli *http.Client, handle string) (string, error) {
	if cli == nil {
		cli = util.RobustHTTPClient()
	}

	var did string

	_, err := syntax.ParseHandle(handle)
	if err != nil {
		return "", err
	}

	recs, err := net.LookupTXT(fmt.Sprintf("_atproto.%s", handle))
	if err == nil {
		for _, rec := range recs {
			if strings.HasPrefix(rec, "did=") {
				did = strings.Split(rec, "did=")[1]
				break
			}
		}
	} else {
		fmt.Printf("error getting txt records: %v\n", err)
	}

	// DNS check fail, fall back to web
	if did == "" {
		req, err := http.NewRequestWithContext(
			ctx,
			"GET",
			fmt.Sprintf("https://%s/.well-known/atproto-did", handle),
			nil,
		)
		if err != nil {
			return "", err
		}

		resp, err := cli.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			io.Copy(io.Discard, resp.Body)
			return "", fmt.Errorf("unable to resolve handle")
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		maybeDid := string(b)

		if _, err := syntax.ParseDID(maybeDid); err != nil {
			return "", fmt.Errorf("unable to resolve handle")
		}

		did = maybeDid
	}

	return did, nil
}

func FetchDidDoc(ctx context.Context, cli *http.Client, did string) (*DidDoc, error) {
	if cli == nil {
		cli = util.RobustHTTPClient()
	}

	var ustr string
	if strings.HasPrefix(did, "did:plc:") {
		ustr = fmt.Sprintf("https://plc.directory/%s", did)
	} else if strings.HasPrefix(did, "did:web:") {
		ustr = fmt.Sprintf("https://%s/.well-known/did.json", strings.TrimPrefix(did, "did:web:"))
	} else {
		return nil, fmt.Errorf("did was not a supported type, must be: did:web or did:plc")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", ustr, nil)
	if err != nil {
		return nil, err
	}

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		io.Copy(io.Discard, resp.Body)
		return nil, fmt.Errorf("could not find did document")
	}

	var didDoc DidDoc
	if err := json.NewDecoder(resp.Body).Decode(&didDoc); err != nil {
		return nil, err
	}

	return &didDoc, nil
}

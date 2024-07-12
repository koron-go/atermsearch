/*
Package atermsearch provides functions to search for Aterm devices.
*/
package atermsearch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Mode struct {
	ID   int
	Name string

	NameJA string
}

var (
	ModeBridge              = Mode{0, "Bridge", "ブリッジ"}
	ModePPPoERouter         = Mode{1, "PPPoE Router", "PPPoEルーター"}
	ModeLocalRouter         = Mode{2, "Local Router", "ローカルルーター"}
	ModeWirelessLANAdapter  = Mode{3, "Wireless LAN Adapter", "無線LAN子機"}
	ModeWirelessLANRepeater = Mode{4, "Wireless LAN Repeater", "無線LAN中継器"}
	ModeMapE                = Mode{5, "MAP-E", "MAP-E"}
	Mode464XLAT             = Mode{6, "464XLAT", "464XLAT"}
	ModeDSLite              = Mode{7, "DS-Lite", "DS-Lite"}
	ModeFixedIPs            = Mode{8, "Fixed IP 1", "固定IP1"}
	ModeMultipleFixedIPs    = Mode{9, "Multiple Fixed IPs", "複数固定IP"}
	ModeMeshRepeater        = Mode{10, "Mesh Repeater", "メッシュ中継器"}
)

var modes = []Mode{
	ModeBridge,
	ModePPPoERouter,
	ModeLocalRouter,
	ModeWirelessLANAdapter,
	ModeWirelessLANRepeater,
	ModeMapE,
	Mode464XLAT,
	ModeDSLite,
	ModeFixedIPs,
	ModeMultipleFixedIPs,
	ModeMeshRepeater,
}

func getParam(ctx context.Context, addr string, reqID string, wantName string) (string, error) {
	u := "http://" + addr + "/aterm_httpif.cgi/getparamcmd_no_auth"
	data := url.Values{"REQ_ID": {reqID}}
	req, err := http.NewRequestWithContext(ctx, "POST", u, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return "", fmt.Errorf("the request to %s for addr %s returns %d, expected 200", u, addr, resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	r := strings.SplitN(string(b), "=", 2)
	if r[0] != wantName {
		return "", fmt.Errorf("unexpected param name, want %q: got body %q", wantName, string(b))
	}
	return strings.TrimSpace(r[1]), nil
}

// ProductName tries to get the product name of addr as an Aterm device.
func ProductName(ctx context.Context, addr string) (string, error) {
	return getParam(ctx, addr, "PRODUCT_NAME_GET", "PRODUCT_NAME")
}

// SystemMode tries to get the system mode of addr as an Aterm device.
func SystemMode(ctx context.Context, addr string) (Mode, error) {
	s, err := getParam(ctx, addr, "SYS_MODE_GET", "SYSTEM_MODE")
	if err != nil {
		return Mode{}, err
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return Mode{}, fmt.Errorf("invalid system mode format: %w", err)
	}
	if n < 0 || n >= len(modes) {
		return Mode{}, fmt.Errorf("unsupported system mode %d", n)
	}
	return modes[n], nil
}

type Device struct {
	Address     string
	ProductName string
	SystemMode  Mode
}

// Search tries to get the device information of addr as an Aterm device.
func Search(ctx context.Context, addr string) (*Device, error) {
	n, err := ProductName(ctx, addr)
	if err != nil {
		return nil, err
	}
	m, err := SystemMode(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &Device{
		Address:     addr,
		ProductName: n,
		SystemMode:  m,
	}, nil
}

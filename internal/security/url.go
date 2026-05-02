package security

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net"
	"net/url"
	"strings"
	"time"
)

var (
	ErrInvalidURL = errors.New("invalid url")
	ErrBlockedIP  = errors.New("blocked private or reserved network")
)

func Normalize(raw string) (string, *url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil, ErrInvalidURL
	}
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil || u.Hostname() == "" {
		return "", nil, ErrInvalidURL
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", nil, ErrInvalidURL
	}
	u.Fragment = ""
	u.Host = strings.ToLower(u.Host)
	return u.String(), u, nil
}

func ValidatePublicHost(ctx context.Context, hostname string) error {
	if hostname == "" {
		return ErrInvalidURL
	}
	lower := strings.ToLower(hostname)
	if lower == "localhost" || strings.HasSuffix(lower, ".localhost") {
		return ErrBlockedIP
	}
	if ip := net.ParseIP(hostname); ip != nil {
		if IsBlockedIP(ip) {
			return ErrBlockedIP
		}
		return nil
	}
	resolver := net.Resolver{}
	lookupCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	ips, err := resolver.LookupIPAddr(lookupCtx, hostname)
	if err != nil || len(ips) == 0 {
		return err
	}
	for _, item := range ips {
		if IsBlockedIP(item.IP) {
			return ErrBlockedIP
		}
	}
	return nil
}

func IsBlockedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast() || ip.IsUnspecified() {
		return true
	}
	privateCIDRs := []string{
		"10.0.0.0/8",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"0.0.0.0/8",
		"100.64.0.0/10",
		"192.0.0.0/24",
		"192.0.2.0/24",
		"198.18.0.0/15",
		"198.51.100.0/24",
		"203.0.113.0/24",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
	}
	for _, cidr := range privateCIDRs {
		_, network, _ := net.ParseCIDR(cidr)
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

func Hash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func HashWithSalt(value, salt string) string {
	sum := sha256.Sum256([]byte(salt + ":" + value))
	return hex.EncodeToString(sum[:])
}

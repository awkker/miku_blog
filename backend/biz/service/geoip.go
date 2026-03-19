package service

import (
	"net"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

type GeoIPResolver struct {
	reader *geoip2.Reader
}

func NewGeoIPResolver(dbPath string) (*GeoIPResolver, error) {
	trimmed := strings.TrimSpace(dbPath)
	if trimmed == "" {
		return nil, nil
	}

	reader, err := geoip2.Open(trimmed)
	if err != nil {
		return nil, err
	}
	return &GeoIPResolver{reader: reader}, nil
}

func (r *GeoIPResolver) Lookup(ipStr string) (countryCode, region, city string, ok bool) {
	if r == nil || r.reader == nil {
		return "", "", "", false
	}

	ip := parseClientIP(ipStr)
	if ip == nil {
		return "", "", "", false
	}
	if !isPublicIP(ip) {
		return "", "", "", false
	}

	record, err := r.reader.City(ip)
	if err != nil || record == nil {
		return "", "", "", false
	}

	countryCode = strings.ToUpper(strings.TrimSpace(record.Country.IsoCode))
	if len(countryCode) != 2 {
		countryCode = ""
	}

	if len(record.Subdivisions) > 0 {
		region = strings.TrimSpace(record.Subdivisions[0].Names["en"])
		if region == "" {
			region = strings.TrimSpace(record.Subdivisions[0].IsoCode)
		}
	}
	city = strings.TrimSpace(record.City.Names["en"])

	if countryCode == "" && region == "" && city == "" {
		return "", "", "", false
	}
	return countryCode, region, city, true
}

func parseClientIP(raw string) net.IP {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}
	if host, _, err := net.SplitHostPort(trimmed); err == nil {
		trimmed = host
	}
	if idx := strings.Index(trimmed, ","); idx >= 0 {
		trimmed = strings.TrimSpace(trimmed[:idx])
	}
	return net.ParseIP(trimmed)
}

func isPublicIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() || ip.IsMulticast() || ip.IsUnspecified() {
		return false
	}
	return true
}

package shared

import (
	"regexp"
	"strings"
)

var mexicoPostalCodePattern = regexp.MustCompile(`^[0-9]{1,5}$`)

type Config struct {
	DefaultCountryCode string
	LimitDefault       int32
	LimitMax           int32
}

type Normalizer struct {
	defaultCountryCode string
	limitDefault       int32
	limitMax           int32
}

func NewNormalizer(cfg Config) Normalizer {
	cfg.DefaultCountryCode = NormalizeCountryCode(cfg.DefaultCountryCode)
	if cfg.DefaultCountryCode == "" {
		cfg.DefaultCountryCode = "MX"
	}
	if cfg.LimitDefault <= 0 {
		cfg.LimitDefault = 20
	}
	if cfg.LimitMax < cfg.LimitDefault {
		cfg.LimitMax = cfg.LimitDefault
	}

	return Normalizer{
		defaultCountryCode: cfg.DefaultCountryCode,
		limitDefault:       cfg.LimitDefault,
		limitMax:           cfg.LimitMax,
	}
}

func (n Normalizer) DefaultCountryCode() string {
	return n.defaultCountryCode
}

func (n Normalizer) NormalizeCountry(value string) string {
	value = NormalizeCountryCode(value)
	if value == "" {
		return n.defaultCountryCode
	}
	if len(value) != 2 {
		return ""
	}
	return value
}

func (n Normalizer) NormalizeLimit(limit int32) int32 {
	if limit <= 0 {
		limit = n.limitDefault
	}
	if limit > n.limitMax {
		limit = n.limitMax
	}
	return limit
}

func (n Normalizer) ValidPostalCode(countryCode, value string) bool {
	if countryCode == "MX" {
		return mexicoPostalCodePattern.MatchString(value)
	}
	return len(value) <= 12
}

func NormalizeCountryCode(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func NormalizeText(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func NormalizeCode(value string) string {
	return strings.TrimSpace(value)
}

func NormalizePostalCode(countryCode, value string) string {
	value = strings.TrimSpace(value)
	if countryCode != "MX" {
		return strings.ToUpper(strings.Join(strings.Fields(value), ""))
	}

	var b strings.Builder
	for _, r := range value {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

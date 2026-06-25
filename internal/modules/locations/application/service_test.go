package application

import (
	"context"
	"errors"
	"testing"

	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/domain"
	locationerrors "github.com/zchelalo/neuraclinic-location/internal/shared/locationerrors"
)

func TestSearchPostalCodesNormalizesMexicoPostalCode(t *testing.T) {
	repo := &memoryRepo{}
	service := NewService(Config{DefaultCountryCode: "MX", LimitDefault: 20, LimitMax: 50}, repo)

	_, err := service.SearchPostalCodes(context.Background(), domain.PostalCodeFilter{
		CountryCode:      "mx",
		PostalCodePrefix: " 832-0 ",
		Limit:            200,
	})
	if err != nil {
		t.Fatalf("SearchPostalCodes returned error: %v", err)
	}

	if repo.postalFilter.CountryCode != "MX" {
		t.Fatalf("expected normalized country MX, got %q", repo.postalFilter.CountryCode)
	}
	if repo.postalFilter.PostalCodePrefix != "8320" {
		t.Fatalf("expected normalized postal prefix 8320, got %q", repo.postalFilter.PostalCodePrefix)
	}
	if repo.postalFilter.Limit != 50 {
		t.Fatalf("expected max limit 50, got %d", repo.postalFilter.Limit)
	}
}

func TestSearchPostalCodesRejectsInvalidMexicoPostalCode(t *testing.T) {
	service := NewService(Config{DefaultCountryCode: "MX", LimitDefault: 20, LimitMax: 50}, &memoryRepo{})

	_, err := service.SearchPostalCodes(context.Background(), domain.PostalCodeFilter{
		CountryCode:      "MX",
		PostalCodePrefix: "123456",
	})
	if !errors.Is(err, locationerrors.ErrInvalidInput) {
		t.Fatalf("expected invalid input, got %v", err)
	}
}

func TestSuggestAddressAllowsNoMatches(t *testing.T) {
	repo := &memoryRepo{suggestions: []domain.AddressSuggestion{}}
	service := NewService(Config{DefaultCountryCode: "MX", LimitDefault: 20, LimitMax: 50}, repo)

	suggestions, err := service.SuggestAddress(context.Background(), domain.AddressSuggestionFilter{
		Query: "colonia que no existe",
	})
	if err != nil {
		t.Fatalf("SuggestAddress returned error: %v", err)
	}
	if len(suggestions) != 0 {
		t.Fatalf("expected no suggestions, got %d", len(suggestions))
	}
	if repo.suggestionFilter.CountryCode != "MX" {
		t.Fatalf("expected default country MX, got %q", repo.suggestionFilter.CountryCode)
	}
}

func TestListAdminAreasTrimsQueryAndDefaultsCountry(t *testing.T) {
	repo := &memoryRepo{}
	service := NewService(Config{DefaultCountryCode: "MX", LimitDefault: 20, LimitMax: 50}, repo)

	_, err := service.ListAdminAreas(context.Background(), domain.AdminAreaFilter{
		Query: "  Sonora   Norte ",
	})
	if err != nil {
		t.Fatalf("ListAdminAreas returned error: %v", err)
	}
	if repo.adminAreaFilter.CountryCode != "MX" {
		t.Fatalf("expected default country MX, got %q", repo.adminAreaFilter.CountryCode)
	}
	if repo.adminAreaFilter.Query != "Sonora Norte" {
		t.Fatalf("expected compact query, got %q", repo.adminAreaFilter.Query)
	}
}

type memoryRepo struct {
	adminAreaFilter  domain.AdminAreaFilter
	postalFilter     domain.PostalCodeFilter
	suggestionFilter domain.AddressSuggestionFilter
	suggestions      []domain.AddressSuggestion
}

func (r *memoryRepo) ListCountries(_ context.Context, filter domain.CountryFilter) ([]domain.Country, error) {
	return nil, nil
}

func (r *memoryRepo) ListAdminAreas(_ context.Context, filter domain.AdminAreaFilter) ([]domain.AdminArea, error) {
	r.adminAreaFilter = filter
	return nil, nil
}

func (r *memoryRepo) ListLocalities(_ context.Context, filter domain.LocalityFilter) ([]domain.Locality, error) {
	return nil, nil
}

func (r *memoryRepo) ListSettlements(_ context.Context, filter domain.SettlementFilter) ([]domain.Settlement, error) {
	return nil, nil
}

func (r *memoryRepo) SearchPostalCodes(_ context.Context, filter domain.PostalCodeFilter) ([]domain.PostalCodeMatch, error) {
	r.postalFilter = filter
	return nil, nil
}

func (r *memoryRepo) SuggestAddress(_ context.Context, filter domain.AddressSuggestionFilter) ([]domain.AddressSuggestion, error) {
	r.suggestionFilter = filter
	return r.suggestions, nil
}

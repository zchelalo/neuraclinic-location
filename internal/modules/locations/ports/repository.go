package ports

import (
	"context"

	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/domain"
)

type Repository interface {
	ListCountries(ctx context.Context, filter domain.CountryFilter) ([]domain.Country, error)
	ListAdminAreas(ctx context.Context, filter domain.AdminAreaFilter) ([]domain.AdminArea, error)
	ListLocalities(ctx context.Context, filter domain.LocalityFilter) ([]domain.Locality, error)
	ListSettlements(ctx context.Context, filter domain.SettlementFilter) ([]domain.Settlement, error)
	SearchPostalCodes(ctx context.Context, filter domain.PostalCodeFilter) ([]domain.PostalCodeMatch, error)
	SuggestAddress(ctx context.Context, filter domain.AddressSuggestionFilter) ([]domain.AddressSuggestion, error)
}

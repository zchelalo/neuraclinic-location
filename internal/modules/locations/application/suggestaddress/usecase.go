package suggestaddress

import (
	"context"

	appshared "github.com/zchelalo/neuraclinic-location/internal/modules/locations/application/shared"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/domain"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/ports"
	locationerrors "github.com/zchelalo/neuraclinic-location/internal/shared/locationerrors"
)

type UseCase struct {
	normalizer appshared.Normalizer
	repo       ports.Repository
}

func New(cfg appshared.Config, repo ports.Repository) *UseCase {
	return &UseCase{
		normalizer: appshared.NewNormalizer(cfg),
		repo:       repo,
	}
}

func (uc *UseCase) Execute(ctx context.Context, filter domain.AddressSuggestionFilter) ([]domain.AddressSuggestion, error) {
	filter.CountryCode = uc.normalizer.NormalizeCountry(filter.CountryCode)
	if filter.CountryCode == "" {
		return nil, locationerrors.ErrInvalidInput
	}
	filter.Query = appshared.NormalizeText(filter.Query)
	filter.PostalCode = appshared.NormalizePostalCode(filter.CountryCode, filter.PostalCode)
	if filter.PostalCode != "" && !uc.normalizer.ValidPostalCode(filter.CountryCode, filter.PostalCode) {
		return nil, locationerrors.ErrInvalidInput
	}
	if filter.Query == "" && filter.PostalCode == "" {
		return nil, locationerrors.ErrInvalidInput
	}
	filter.Limit = uc.normalizer.NormalizeLimit(filter.Limit)
	return uc.repo.SuggestAddress(ctx, filter)
}

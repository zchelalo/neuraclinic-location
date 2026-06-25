package listlocalities

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

func (uc *UseCase) Execute(ctx context.Context, filter domain.LocalityFilter) ([]domain.Locality, error) {
	filter.CountryCode = uc.normalizer.NormalizeCountry(filter.CountryCode)
	if filter.CountryCode == "" {
		return nil, locationerrors.ErrInvalidInput
	}
	filter.AdminAreaCode = appshared.NormalizeCode(filter.AdminAreaCode)
	filter.Query = appshared.NormalizeText(filter.Query)
	filter.Limit = uc.normalizer.NormalizeLimit(filter.Limit)
	return uc.repo.ListLocalities(ctx, filter)
}

package listcountries

import (
	"context"

	appshared "github.com/zchelalo/neuraclinic-location/internal/modules/locations/application/shared"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/domain"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/ports"
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

func (uc *UseCase) Execute(ctx context.Context, filter domain.CountryFilter) ([]domain.Country, error) {
	filter.Query = appshared.NormalizeText(filter.Query)
	filter.Limit = uc.normalizer.NormalizeLimit(filter.Limit)
	return uc.repo.ListCountries(ctx, filter)
}

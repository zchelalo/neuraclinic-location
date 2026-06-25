package importinegiageeml

import (
	"context"

	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/domain"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/ports"
)

type Command struct {
	Source         domain.DataSource
	Entities       []domain.INEGIEntityRow
	Municipalities []domain.INEGIMunicipalityRow
}

type UseCase struct {
	repo ports.ImportRepository
}

func New(repo ports.ImportRepository) UseCase {
	return UseCase{repo: repo}
}

func (uc UseCase) Execute(ctx context.Context, cmd Command) (domain.ImportSummary, error) {
	return uc.repo.ImportINEGIAgeeml(ctx, cmd.Source, cmd.Entities, cmd.Municipalities)
}

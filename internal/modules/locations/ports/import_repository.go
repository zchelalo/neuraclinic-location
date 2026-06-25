package ports

import (
	"context"

	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/domain"
)

type ImportRepository interface {
	ImportSEPOMEX(ctx context.Context, source domain.DataSource, rows []domain.SEPOMEXRow) (domain.ImportSummary, error)
	ImportINEGIAgeeml(ctx context.Context, source domain.DataSource, entities []domain.INEGIEntityRow, municipalities []domain.INEGIMunicipalityRow) (domain.ImportSummary, error)
}

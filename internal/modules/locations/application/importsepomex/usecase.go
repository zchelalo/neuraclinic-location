package importsepomex

import (
	"context"
	"strings"

	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/domain"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/ports"
	locationerrors "github.com/zchelalo/neuraclinic-location/internal/shared/locationerrors"
)

type Command struct {
	Source domain.DataSource
	Rows   []domain.SEPOMEXRow
}

type UseCase struct {
	repo ports.ImportRepository
}

func New(repo ports.ImportRepository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (domain.ImportSummary, error) {
	cmd.Source.Key = strings.TrimSpace(cmd.Source.Key)
	cmd.Source.Name = strings.TrimSpace(cmd.Source.Name)
	cmd.Source.Version = strings.TrimSpace(cmd.Source.Version)
	cmd.Source.License = strings.TrimSpace(cmd.Source.License)
	cmd.Source.Attribution = strings.TrimSpace(cmd.Source.Attribution)
	cmd.Source.URL = strings.TrimSpace(cmd.Source.URL)
	cmd.Source.CountryCode = strings.ToUpper(strings.TrimSpace(cmd.Source.CountryCode))
	cmd.Source.CountryName = strings.TrimSpace(cmd.Source.CountryName)

	if cmd.Source.Key == "" || cmd.Source.Name == "" || cmd.Source.Version == "" ||
		cmd.Source.License == "" || cmd.Source.Attribution == "" ||
		cmd.Source.CountryCode == "" || cmd.Source.CountryName == "" ||
		len(cmd.Source.CountryCode) != 2 {
		return domain.ImportSummary{}, locationerrors.ErrInvalidInput
	}
	if len(cmd.Rows) == 0 {
		return domain.ImportSummary{}, locationerrors.ErrInvalidInput
	}

	return uc.repo.ImportSEPOMEX(ctx, cmd.Source, cmd.Rows)
}

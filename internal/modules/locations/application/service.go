package application

import (
	"context"

	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/application/listadminareas"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/application/listcountries"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/application/listlocalities"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/application/listsettlements"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/application/searchpostalcodes"
	appshared "github.com/zchelalo/neuraclinic-location/internal/modules/locations/application/shared"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/application/suggestaddress"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/domain"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/ports"
)

type Config = appshared.Config

type Service struct {
	listCountries    *listcountries.UseCase
	listAdminAreas   *listadminareas.UseCase
	listLocalities   *listlocalities.UseCase
	listSettlements  *listsettlements.UseCase
	searchPostalCode *searchpostalcodes.UseCase
	suggestAddress   *suggestaddress.UseCase
}

func NewService(cfg Config, repo ports.Repository) *Service {
	return &Service{
		listCountries:    listcountries.New(cfg, repo),
		listAdminAreas:   listadminareas.New(cfg, repo),
		listLocalities:   listlocalities.New(cfg, repo),
		listSettlements:  listsettlements.New(cfg, repo),
		searchPostalCode: searchpostalcodes.New(cfg, repo),
		suggestAddress:   suggestaddress.New(cfg, repo),
	}
}

func (s *Service) ListCountries(ctx context.Context, filter domain.CountryFilter) ([]domain.Country, error) {
	return s.listCountries.Execute(ctx, filter)
}

func (s *Service) ListAdminAreas(ctx context.Context, filter domain.AdminAreaFilter) ([]domain.AdminArea, error) {
	return s.listAdminAreas.Execute(ctx, filter)
}

func (s *Service) ListLocalities(ctx context.Context, filter domain.LocalityFilter) ([]domain.Locality, error) {
	return s.listLocalities.Execute(ctx, filter)
}

func (s *Service) ListSettlements(ctx context.Context, filter domain.SettlementFilter) ([]domain.Settlement, error) {
	return s.listSettlements.Execute(ctx, filter)
}

func (s *Service) SearchPostalCodes(ctx context.Context, filter domain.PostalCodeFilter) ([]domain.PostalCodeMatch, error) {
	return s.searchPostalCode.Execute(ctx, filter)
}

func (s *Service) SuggestAddress(ctx context.Context, filter domain.AddressSuggestionFilter) ([]domain.AddressSuggestion, error) {
	return s.suggestAddress.Execute(ctx, filter)
}

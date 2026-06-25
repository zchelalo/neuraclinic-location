package grpc

import (
	"context"

	locationv1 "github.com/zchelalo/neuraclinic-location/gen/go/location/v1"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/domain"
	locationerrors "github.com/zchelalo/neuraclinic-location/internal/shared/locationerrors"
	locationgrpc "github.com/zchelalo/neuraclinic-location/internal/shared/locationgrpc"
)

func (s *LocationService) ListCountries(ctx context.Context, req *locationv1.ListCountriesRequest) (*locationv1.ListCountriesResponse, error) {
	countries, err := s.app.ListCountries(ctx, domain.CountryFilter{
		Query: optionalValue(req.Query),
		Limit: req.GetLimit(),
	})
	if err != nil {
		return nil, locationgrpc.MapError(err)
	}

	resp := &locationv1.ListCountriesResponse{Countries: make([]*locationv1.Country, 0, len(countries))}
	for _, country := range countries {
		resp.Countries = append(resp.Countries, countryToProto(country))
	}
	return resp, nil
}

func (s *LocationService) ListAdminAreas(ctx context.Context, req *locationv1.ListAdminAreasRequest) (*locationv1.ListAdminAreasResponse, error) {
	typeFilter, err := adminAreaTypeFilterFromProto(req)
	if err != nil {
		return nil, locationgrpc.MapError(err)
	}

	adminAreas, err := s.app.ListAdminAreas(ctx, domain.AdminAreaFilter{
		CountryCode: req.GetCountryCode(),
		ParentCode:  optionalValue(req.ParentCode),
		Type:        typeFilter,
		Query:       optionalValue(req.Query),
		Limit:       req.GetLimit(),
	})
	if err != nil {
		return nil, locationgrpc.MapError(err)
	}

	resp := &locationv1.ListAdminAreasResponse{AdminAreas: make([]*locationv1.AdminArea, 0, len(adminAreas))}
	for _, adminArea := range adminAreas {
		resp.AdminAreas = append(resp.AdminAreas, adminAreaToProto(adminArea))
	}
	return resp, nil
}

func (s *LocationService) ListLocalities(ctx context.Context, req *locationv1.ListLocalitiesRequest) (*locationv1.ListLocalitiesResponse, error) {
	localities, err := s.app.ListLocalities(ctx, domain.LocalityFilter{
		CountryCode:   req.GetCountryCode(),
		AdminAreaCode: optionalValue(req.AdminAreaCode),
		Query:         optionalValue(req.Query),
		Limit:         req.GetLimit(),
	})
	if err != nil {
		return nil, locationgrpc.MapError(err)
	}

	resp := &locationv1.ListLocalitiesResponse{Localities: make([]*locationv1.Locality, 0, len(localities))}
	for _, locality := range localities {
		resp.Localities = append(resp.Localities, localityToProto(locality))
	}
	return resp, nil
}

func (s *LocationService) ListSettlements(ctx context.Context, req *locationv1.ListSettlementsRequest) (*locationv1.ListSettlementsResponse, error) {
	settlements, err := s.app.ListSettlements(ctx, domain.SettlementFilter{
		CountryCode:   req.GetCountryCode(),
		AdminAreaCode: optionalValue(req.AdminAreaCode),
		LocalityCode:  optionalValue(req.LocalityCode),
		PostalCode:    optionalValue(req.PostalCode),
		Query:         optionalValue(req.Query),
		Limit:         req.GetLimit(),
	})
	if err != nil {
		return nil, locationgrpc.MapError(err)
	}

	resp := &locationv1.ListSettlementsResponse{Settlements: make([]*locationv1.Settlement, 0, len(settlements))}
	for _, settlement := range settlements {
		resp.Settlements = append(resp.Settlements, settlementToProto(settlement))
	}
	return resp, nil
}

func (s *LocationService) SearchPostalCodes(ctx context.Context, req *locationv1.SearchPostalCodesRequest) (*locationv1.SearchPostalCodesResponse, error) {
	postalCodes, err := s.app.SearchPostalCodes(ctx, domain.PostalCodeFilter{
		CountryCode:      req.GetCountryCode(),
		PostalCodePrefix: req.GetPostalCode(),
		Limit:            req.GetLimit(),
	})
	if err != nil {
		return nil, locationgrpc.MapError(err)
	}

	resp := &locationv1.SearchPostalCodesResponse{PostalCodes: make([]*locationv1.PostalCodeMatch, 0, len(postalCodes))}
	for _, postalCode := range postalCodes {
		resp.PostalCodes = append(resp.PostalCodes, postalCodeToProto(postalCode))
	}
	return resp, nil
}

func (s *LocationService) SuggestAddress(ctx context.Context, req *locationv1.SuggestAddressRequest) (*locationv1.SuggestAddressResponse, error) {
	suggestions, err := s.app.SuggestAddress(ctx, domain.AddressSuggestionFilter{
		CountryCode: req.GetCountryCode(),
		Query:       req.GetQuery(),
		PostalCode:  optionalValue(req.PostalCode),
		Limit:       req.GetLimit(),
	})
	if err != nil {
		return nil, locationgrpc.MapError(err)
	}

	resp := &locationv1.SuggestAddressResponse{Suggestions: make([]*locationv1.AddressSuggestion, 0, len(suggestions))}
	for _, suggestion := range suggestions {
		resp.Suggestions = append(resp.Suggestions, suggestionToProto(suggestion))
	}
	return resp, nil
}

func countryToProto(country domain.Country) *locationv1.Country {
	return &locationv1.Country{
		CountryCode:   country.CountryCode,
		Name:          country.Name,
		Label:         country.Label,
		Source:        country.Source,
		SourceVersion: country.SourceVersion,
	}
}

func adminAreaToProto(adminArea domain.AdminArea) *locationv1.AdminArea {
	return &locationv1.AdminArea{
		Id:            adminArea.ID,
		CountryCode:   adminArea.CountryCode,
		Code:          adminArea.Code,
		Name:          adminArea.Name,
		ParentCode:    adminArea.ParentCode,
		Label:         adminArea.Label,
		Source:        adminArea.Source,
		SourceVersion: adminArea.SourceVersion,
		AdminAreaType: adminAreaTypeToProto(adminArea.Type),
	}
}

func localityToProto(locality domain.Locality) *locationv1.Locality {
	return &locationv1.Locality{
		Id:            locality.ID,
		CountryCode:   locality.CountryCode,
		AdminAreaCode: locality.AdminAreaCode,
		Code:          locality.Code,
		Name:          locality.Name,
		Type:          locality.Type,
		Label:         locality.Label,
		Source:        locality.Source,
		SourceVersion: locality.SourceVersion,
	}
}

func settlementToProto(settlement domain.Settlement) *locationv1.Settlement {
	return &locationv1.Settlement{
		Id:            settlement.ID,
		CountryCode:   settlement.CountryCode,
		AdminAreaCode: settlement.AdminAreaCode,
		LocalityCode:  settlement.LocalityCode,
		PostalCode:    settlement.PostalCode,
		Name:          settlement.Name,
		Type:          settlement.Type,
		Label:         settlement.Label,
		Source:        settlement.Source,
		SourceVersion: settlement.SourceVersion,
	}
}

func postalCodeToProto(postalCode domain.PostalCodeMatch) *locationv1.PostalCodeMatch {
	return &locationv1.PostalCodeMatch{
		PostalCode:    postalCode.PostalCode,
		Label:         postalCode.Label,
		Components:    componentsToProto(postalCode.Components),
		Source:        postalCode.Source,
		SourceVersion: postalCode.SourceVersion,
		Score:         postalCode.Score,
	}
}

func suggestionToProto(suggestion domain.AddressSuggestion) *locationv1.AddressSuggestion {
	return &locationv1.AddressSuggestion{
		Label:         suggestion.Label,
		Components:    componentsToProto(suggestion.Components),
		Source:        suggestion.Source,
		SourceVersion: suggestion.SourceVersion,
		Score:         suggestion.Score,
	}
}

func componentsToProto(components domain.Components) *locationv1.LocationComponents {
	return &locationv1.LocationComponents{
		CountryCode:    components.CountryCode,
		CountryName:    components.CountryName,
		AdminAreaCode:  components.AdminAreaCode,
		AdminAreaName:  components.AdminAreaName,
		LocalityCode:   components.LocalityCode,
		LocalityName:   components.LocalityName,
		PostalCode:     components.PostalCode,
		SettlementName: components.SettlementName,
		SettlementType: components.SettlementType,
		StreetName:     components.StreetName,
	}
}

func optionalValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func adminAreaTypeFilterFromProto(req *locationv1.ListAdminAreasRequest) (string, error) {
	if req == nil || req.AdminAreaType == nil {
		return "", nil
	}
	return adminAreaTypeFromProto(req.GetAdminAreaType())
}

func adminAreaTypeFromProto(value locationv1.AdminAreaType) (string, error) {
	switch value {
	case locationv1.AdminAreaType_ADMIN_AREA_TYPE_UNSPECIFIED:
		return "", nil
	case locationv1.AdminAreaType_ADMIN_AREA_TYPE_STATE:
		return domain.AdminAreaTypeState, nil
	case locationv1.AdminAreaType_ADMIN_AREA_TYPE_MUNICIPALITY:
		return domain.AdminAreaTypeMunicipality, nil
	default:
		return "", locationerrors.ErrInvalidInput
	}
}

func adminAreaTypeToProto(value string) locationv1.AdminAreaType {
	switch value {
	case domain.AdminAreaTypeState:
		return locationv1.AdminAreaType_ADMIN_AREA_TYPE_STATE
	case domain.AdminAreaTypeMunicipality:
		return locationv1.AdminAreaType_ADMIN_AREA_TYPE_MUNICIPALITY
	default:
		return locationv1.AdminAreaType_ADMIN_AREA_TYPE_UNSPECIFIED
	}
}

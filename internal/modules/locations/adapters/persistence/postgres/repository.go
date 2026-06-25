package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	locationsdb "github.com/zchelalo/neuraclinic-location/internal/db/sqlc/locations"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/domain"
	pgutil "github.com/zchelalo/neuraclinic-location/internal/shared/postgresutil"
)

type Repository struct {
	db *pgxpool.Pool
	q  *locationsdb.Queries
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db, q: locationsdb.New(db)}
}

func (r *Repository) ListCountries(ctx context.Context, filter domain.CountryFilter) ([]domain.Country, error) {
	rows, err := r.q.ListCountries(ctx, locationsdb.ListCountriesParams{
		Query:      filter.Query,
		LimitCount: filter.Limit,
	})
	if err != nil {
		return nil, err
	}

	countries := make([]domain.Country, 0, len(rows))
	for _, row := range rows {
		countries = append(countries, domain.Country{
			CountryCode:   row.CountryCode,
			Name:          row.Name,
			Label:         row.Label,
			Source:        row.Source,
			SourceVersion: row.SourceVersion,
			Score:         row.Score,
		})
	}
	return countries, nil
}

func (r *Repository) ListAdminAreas(ctx context.Context, filter domain.AdminAreaFilter) ([]domain.AdminArea, error) {
	rows, err := r.q.ListAdminAreas(ctx, locationsdb.ListAdminAreasParams{
		CountryCode: filter.CountryCode,
		ParentCode:  filter.ParentCode,
		TypeFilter:  filter.Type,
		Query:       filter.Query,
		LimitCount:  filter.Limit,
	})
	if err != nil {
		return nil, err
	}

	adminAreas := make([]domain.AdminArea, 0, len(rows))
	for _, row := range rows {
		adminAreas = append(adminAreas, domain.AdminArea{
			ID:            row.ID,
			CountryCode:   row.CountryCode,
			Code:          row.Code,
			Name:          row.Name,
			Type:          row.Type,
			ParentCode:    pgutil.TextPtr(row.ParentCode),
			Label:         row.Label,
			Source:        row.Source,
			SourceVersion: row.SourceVersion,
			Score:         row.Score,
		})
	}
	return adminAreas, nil
}

func (r *Repository) ListLocalities(ctx context.Context, filter domain.LocalityFilter) ([]domain.Locality, error) {
	rows, err := r.q.ListLocalities(ctx, locationsdb.ListLocalitiesParams{
		CountryCode:   filter.CountryCode,
		AdminAreaCode: filter.AdminAreaCode,
		Query:         filter.Query,
		LimitCount:    filter.Limit,
	})
	if err != nil {
		return nil, err
	}

	localities := make([]domain.Locality, 0, len(rows))
	for _, row := range rows {
		localities = append(localities, domain.Locality{
			ID:            row.ID,
			CountryCode:   row.CountryCode,
			AdminAreaCode: row.AdminAreaCode,
			Code:          row.Code,
			Name:          row.Name,
			Type:          row.Type,
			Label:         row.Label,
			Source:        row.Source,
			SourceVersion: row.SourceVersion,
			Score:         row.Score,
		})
	}
	return localities, nil
}

func (r *Repository) ListSettlements(ctx context.Context, filter domain.SettlementFilter) ([]domain.Settlement, error) {
	rows, err := r.q.ListSettlements(ctx, locationsdb.ListSettlementsParams{
		CountryCode:   filter.CountryCode,
		AdminAreaCode: filter.AdminAreaCode,
		LocalityCode:  filter.LocalityCode,
		PostalCode:    filter.PostalCode,
		Query:         filter.Query,
		LimitCount:    filter.Limit,
	})
	if err != nil {
		return nil, err
	}

	settlements := make([]domain.Settlement, 0, len(rows))
	for _, row := range rows {
		settlements = append(settlements, domain.Settlement{
			ID:            row.ID,
			CountryCode:   row.CountryCode,
			AdminAreaCode: row.AdminAreaCode,
			LocalityCode:  pgutil.TextPtr(row.LocalityCode),
			PostalCode:    pgutil.TextPtr(row.PostalCode),
			Name:          row.Name,
			Type:          row.Type,
			Label:         row.Label,
			Source:        row.Source,
			SourceVersion: row.SourceVersion,
			Score:         row.Score,
		})
	}
	return settlements, nil
}

func (r *Repository) SearchPostalCodes(ctx context.Context, filter domain.PostalCodeFilter) ([]domain.PostalCodeMatch, error) {
	rows, err := r.q.SearchPostalCodes(ctx, locationsdb.SearchPostalCodesParams{
		CountryCode:      filter.CountryCode,
		PostalCodePrefix: filter.PostalCodePrefix,
		LimitCount:       filter.Limit,
	})
	if err != nil {
		return nil, err
	}

	results := make([]domain.PostalCodeMatch, 0, len(rows))
	for _, row := range rows {
		results = append(results, domain.PostalCodeMatch{
			PostalCode: row.PostalCode,
			Label:      row.Label,
			Components: domain.Components{
				CountryCode:    row.CountryCode,
				CountryName:    row.CountryName,
				AdminAreaCode:  &row.AdminAreaCode,
				AdminAreaName:  &row.AdminAreaName,
				LocalityCode:   pgutil.TextPtr(row.LocalityCode),
				LocalityName:   pgutil.TextPtr(row.LocalityName),
				PostalCode:     &row.PostalCode,
				SettlementName: pgutil.TextPtr(row.SettlementName),
				SettlementType: pgutil.TextPtr(row.SettlementType),
				StreetName:     nil,
			},
			Source:        row.Source,
			SourceVersion: row.SourceVersion,
			Score:         row.Score,
		})
	}
	return results, nil
}

func (r *Repository) SuggestAddress(ctx context.Context, filter domain.AddressSuggestionFilter) ([]domain.AddressSuggestion, error) {
	rows, err := r.q.SuggestAddresses(ctx, locationsdb.SuggestAddressesParams{
		CountryCode: filter.CountryCode,
		PostalCode:  filter.PostalCode,
		Query:       filter.Query,
		LimitCount:  filter.Limit,
	})
	if err != nil {
		return nil, err
	}

	suggestions := make([]domain.AddressSuggestion, 0, len(rows))
	for _, row := range rows {
		suggestions = append(suggestions, domain.AddressSuggestion{
			Label: row.Label,
			Components: domain.Components{
				CountryCode:    row.CountryCode,
				CountryName:    row.CountryName,
				AdminAreaCode:  &row.AdminAreaCode,
				AdminAreaName:  &row.AdminAreaName,
				LocalityCode:   pgutil.TextPtr(row.LocalityCode),
				LocalityName:   pgutil.TextPtr(row.LocalityName),
				PostalCode:     pgutil.TextPtr(row.PostalCode),
				SettlementName: &row.SettlementName,
				SettlementType: &row.SettlementType,
				StreetName:     pgutil.TextPtr(row.StreetName),
			},
			Source:        row.Source,
			SourceVersion: row.SourceVersion,
			Score:         row.Score,
		})
	}
	return suggestions, nil
}

func (r *Repository) ImportSEPOMEX(ctx context.Context, source domain.DataSource, rows []domain.SEPOMEXRow) (domain.ImportSummary, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.ImportSummary{}, err
	}
	defer rollback(ctx, tx)

	sourceID, err := upsertDataSource(ctx, tx, source)
	if err != nil {
		return domain.ImportSummary{}, fmt.Errorf("upsert data source: %w", err)
	}
	if err := upsertCountry(ctx, tx, source, sourceID); err != nil {
		return domain.ImportSummary{}, fmt.Errorf("upsert country: %w", err)
	}

	summary := domain.ImportSummary{
		RowsRead:    len(rows),
		DataSources: 1,
		Countries:   1,
	}
	seenStates := map[string]struct{}{}
	seenMunicipalities := map[string]struct{}{}
	seenPostalCodes := map[string]struct{}{}
	seenSettlements := map[string]struct{}{}

	for _, row := range rows {
		row = cleanSEPOMEXRow(row)
		if !validSEPOMEXRow(row) {
			continue
		}

		stateID := stableUUID(strings.Join([]string{"admin-area", source.CountryCode, domain.AdminAreaTypeState, row.StateCode}, ":"))
		if _, seen := seenStates[row.StateCode]; !seen {
			if _, err := tx.Exec(ctx, `
				INSERT INTO admin_areas (
				  id, country_code, parent_id, code, name, type, source_id, source_record_id, source_version
				) VALUES (
				  $1, $2, NULL, $3, $4, $5, $6, $7, $8
				)
				ON CONFLICT (id) DO UPDATE
				SET name = EXCLUDED.name,
				    source_id = EXCLUDED.source_id,
				    source_record_id = EXCLUDED.source_record_id,
				    source_version = EXCLUDED.source_version
			`, stateID, source.CountryCode, row.StateCode, row.State, domain.AdminAreaTypeState, sourceID, row.StateCode, source.Version); err != nil {
				return domain.ImportSummary{}, fmt.Errorf("upsert state %s: %w", row.StateCode, err)
			}
			seenStates[row.StateCode] = struct{}{}
			summary.AdminAreas++
		}

		municipalityKey := strings.Join([]string{row.StateCode, row.MunicipalityCode}, ":")
		municipalityID := stableUUID(strings.Join([]string{"locality", source.CountryCode, "municipality", municipalityKey}, ":"))
		if _, seen := seenMunicipalities[municipalityKey]; !seen {
			if _, err := tx.Exec(ctx, `
				INSERT INTO localities (
				  id, country_code, admin_area_id, code, name, type, source_id, source_record_id, source_version
				) VALUES (
				  $1, $2, $3, $4, $5, 'municipality', $6, $7, $8
				)
				ON CONFLICT (id) DO UPDATE
				SET admin_area_id = EXCLUDED.admin_area_id,
				    name = EXCLUDED.name,
				    source_id = EXCLUDED.source_id,
				    source_record_id = EXCLUDED.source_record_id,
				    source_version = EXCLUDED.source_version
			`, municipalityID, source.CountryCode, stateID, row.MunicipalityCode, row.Municipality, sourceID, municipalityKey, source.Version); err != nil {
				return domain.ImportSummary{}, fmt.Errorf("upsert municipality %s: %w", municipalityKey, err)
			}
			seenMunicipalities[municipalityKey] = struct{}{}
			summary.Localities++
		}

		postalCodeKey := strings.Join([]string{row.StateCode, row.MunicipalityCode, row.PostalCode}, ":")
		postalCodeID := stableUUID(strings.Join([]string{"postal-code", source.CountryCode, postalCodeKey}, ":"))
		if _, seen := seenPostalCodes[postalCodeKey]; !seen {
			if _, err := tx.Exec(ctx, `
				INSERT INTO postal_codes (
				  id, country_code, admin_area_id, locality_id, postal_code, source_id, source_record_id, source_version
				) VALUES (
				  $1, $2, $3, $4, $5, $6, $7, $8
				)
				ON CONFLICT (id) DO UPDATE
				SET admin_area_id = EXCLUDED.admin_area_id,
				    locality_id = EXCLUDED.locality_id,
				    postal_code = EXCLUDED.postal_code,
				    source_id = EXCLUDED.source_id,
				    source_record_id = EXCLUDED.source_record_id,
				    source_version = EXCLUDED.source_version
			`, postalCodeID, source.CountryCode, stateID, municipalityID, row.PostalCode, sourceID, postalCodeKey, source.Version); err != nil {
				return domain.ImportSummary{}, fmt.Errorf("upsert postal code %s: %w", postalCodeKey, err)
			}
			seenPostalCodes[postalCodeKey] = struct{}{}
			summary.PostalCodes++
		}

		settlementKey := strings.Join([]string{row.StateCode, row.MunicipalityCode, row.PostalCode, row.SettlementCode}, ":")
		settlementID := stableUUID(strings.Join([]string{"settlement", source.CountryCode, settlementKey}, ":"))
		if _, seen := seenSettlements[settlementKey]; !seen {
			if _, err := tx.Exec(ctx, `
				INSERT INTO settlements (
				  id, country_code, admin_area_id, locality_id, postal_code_id, code, name,
				  settlement_type, source_id, source_record_id, source_version
				) VALUES (
				  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
				)
				ON CONFLICT (id) DO UPDATE
				SET admin_area_id = EXCLUDED.admin_area_id,
				    locality_id = EXCLUDED.locality_id,
				    postal_code_id = EXCLUDED.postal_code_id,
				    code = EXCLUDED.code,
				    name = EXCLUDED.name,
				    settlement_type = EXCLUDED.settlement_type,
				    source_id = EXCLUDED.source_id,
				    source_record_id = EXCLUDED.source_record_id,
				    source_version = EXCLUDED.source_version
			`, settlementID, source.CountryCode, stateID, municipalityID, postalCodeID, row.SettlementCode, row.SettlementName, row.SettlementType, sourceID, settlementKey, source.Version); err != nil {
				return domain.ImportSummary{}, fmt.Errorf("upsert settlement %s: %w", settlementKey, err)
			}
			seenSettlements[settlementKey] = struct{}{}
			summary.Settlements++
		}
		summary.RowsImported++
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.ImportSummary{}, err
	}
	return summary, nil
}

func (r *Repository) ImportINEGIAgeeml(ctx context.Context, source domain.DataSource, entities []domain.INEGIEntityRow, municipalities []domain.INEGIMunicipalityRow) (domain.ImportSummary, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return domain.ImportSummary{}, err
	}
	defer rollback(ctx, tx)

	sourceID, err := upsertDataSource(ctx, tx, source)
	if err != nil {
		return domain.ImportSummary{}, fmt.Errorf("upsert data source: %w", err)
	}
	if err := upsertCountry(ctx, tx, source, sourceID); err != nil {
		return domain.ImportSummary{}, fmt.Errorf("upsert country: %w", err)
	}

	summary := domain.ImportSummary{
		RowsRead:    len(entities) + len(municipalities),
		DataSources: 1,
		Countries:   1,
	}
	seenStates := map[string]struct{}{}
	seenMunicipalities := map[string]struct{}{}

	for _, row := range entities {
		row = cleanINEGIEntityRow(row)
		if !validINEGIEntityRow(row) {
			continue
		}
		if _, seen := seenStates[row.StateCode]; seen {
			continue
		}
		if err := upsertINEGIState(ctx, tx, source, sourceID, row.StateCode, row.StateName, row.GeoCode); err != nil {
			return domain.ImportSummary{}, fmt.Errorf("upsert INEGI state %s: %w", row.StateCode, err)
		}
		seenStates[row.StateCode] = struct{}{}
		summary.AdminAreas++
		summary.RowsImported++
	}

	for _, row := range municipalities {
		row = cleanINEGIMunicipalityRow(row)
		if !validINEGIMunicipalityRow(row) {
			continue
		}
		if _, seen := seenStates[row.StateCode]; !seen {
			if err := upsertINEGIState(ctx, tx, source, sourceID, row.StateCode, row.StateName, row.StateCode); err != nil {
				return domain.ImportSummary{}, fmt.Errorf("upsert INEGI state %s: %w", row.StateCode, err)
			}
			seenStates[row.StateCode] = struct{}{}
			summary.AdminAreas++
		}

		if _, seen := seenMunicipalities[row.GeoCode]; seen {
			continue
		}
		stateID := stableUUID(strings.Join([]string{"admin-area", source.CountryCode, domain.AdminAreaTypeState, row.StateCode}, ":"))
		municipalityID := stableUUID(strings.Join([]string{"admin-area", source.CountryCode, domain.AdminAreaTypeMunicipality, row.GeoCode}, ":"))
		if _, err := tx.Exec(ctx, `
			INSERT INTO admin_areas (
			  id, country_code, parent_id, code, name, type, source_id, source_record_id, source_version
			) VALUES (
			  $1, $2, $3, $4, $5, $6, $7, $8, $9
			)
			ON CONFLICT (id) DO UPDATE
			SET parent_id = EXCLUDED.parent_id,
			    code = EXCLUDED.code,
			    name = EXCLUDED.name,
			    source_id = EXCLUDED.source_id,
			    source_record_id = EXCLUDED.source_record_id,
			    source_version = EXCLUDED.source_version
		`, municipalityID, source.CountryCode, stateID, row.GeoCode, row.MunicipalityName, domain.AdminAreaTypeMunicipality, sourceID, row.GeoCode, source.Version); err != nil {
			return domain.ImportSummary{}, fmt.Errorf("upsert INEGI municipality %s: %w", row.GeoCode, err)
		}
		seenMunicipalities[row.GeoCode] = struct{}{}
		summary.AdminAreas++
		summary.RowsImported++
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.ImportSummary{}, err
	}
	return summary, nil
}

func upsertDataSource(ctx context.Context, tx pgx.Tx, source domain.DataSource) (uuid.UUID, error) {
	sourceID := stableUUID("data-source:" + source.Key)
	if _, err := tx.Exec(ctx, `
		INSERT INTO data_sources (id, key, name, version, license, attribution, url, imported_at)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, ''), $8)
		ON CONFLICT (key) DO UPDATE
		SET name = EXCLUDED.name,
		    version = EXCLUDED.version,
		    license = EXCLUDED.license,
		    attribution = EXCLUDED.attribution,
		    url = EXCLUDED.url,
		    imported_at = EXCLUDED.imported_at
	`, sourceID, source.Key, source.Name, source.Version, source.License, source.Attribution, source.URL, time.Now().UTC()); err != nil {
		return uuid.Nil, err
	}
	return sourceID, nil
}

func upsertCountry(ctx context.Context, tx pgx.Tx, source domain.DataSource, sourceID uuid.UUID) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO countries (id, country_code, name, source_id, source_record_id, source_version)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (country_code) DO UPDATE
		SET name = EXCLUDED.name,
		    source_id = EXCLUDED.source_id,
		    source_record_id = EXCLUDED.source_record_id,
		    source_version = EXCLUDED.source_version
	`, stableUUID("country:"+source.CountryCode), source.CountryCode, source.CountryName, sourceID, source.CountryCode, source.Version)
	return err
}

func upsertINEGIState(ctx context.Context, tx pgx.Tx, source domain.DataSource, sourceID uuid.UUID, stateCode, stateName, sourceRecordID string) error {
	stateID := stableUUID(strings.Join([]string{"admin-area", source.CountryCode, domain.AdminAreaTypeState, stateCode}, ":"))
	_, err := tx.Exec(ctx, `
		INSERT INTO admin_areas (
		  id, country_code, parent_id, code, name, type, source_id, source_record_id, source_version
		) VALUES (
		  $1, $2, NULL, $3, $4, $5, $6, $7, $8
		)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
		    source_id = EXCLUDED.source_id,
		    source_record_id = EXCLUDED.source_record_id,
		    source_version = EXCLUDED.source_version
	`, stateID, source.CountryCode, stateCode, stateName, domain.AdminAreaTypeState, sourceID, sourceRecordID, source.Version)
	return err
}

func rollback(ctx context.Context, tx pgx.Tx) {
	_ = tx.Rollback(ctx)
}

func stableUUID(value string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(value))
}

func cleanSEPOMEXRow(row domain.SEPOMEXRow) domain.SEPOMEXRow {
	row.PostalCode = strings.TrimSpace(row.PostalCode)
	row.SettlementName = normalizeSpaces(row.SettlementName)
	row.SettlementType = normalizeSpaces(row.SettlementType)
	row.Municipality = normalizeSpaces(row.Municipality)
	row.State = normalizeSpaces(row.State)
	row.City = normalizeSpaces(row.City)
	row.StateCode = strings.TrimSpace(row.StateCode)
	row.MunicipalityCode = strings.TrimSpace(row.MunicipalityCode)
	row.SettlementCode = strings.TrimSpace(row.SettlementCode)
	row.Zone = normalizeSpaces(row.Zone)
	row.CityCode = strings.TrimSpace(row.CityCode)
	return row
}

func validSEPOMEXRow(row domain.SEPOMEXRow) bool {
	return row.PostalCode != "" &&
		row.SettlementName != "" &&
		row.SettlementType != "" &&
		row.Municipality != "" &&
		row.State != "" &&
		row.StateCode != "" &&
		row.MunicipalityCode != "" &&
		row.SettlementCode != ""
}

func normalizeSpaces(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func cleanINEGIEntityRow(row domain.INEGIEntityRow) domain.INEGIEntityRow {
	row.GeoCode = strings.TrimSpace(row.GeoCode)
	row.StateCode = strings.TrimSpace(row.StateCode)
	row.StateName = normalizeSpaces(row.StateName)
	row.StateAbbreviation = normalizeSpaces(row.StateAbbreviation)
	return row
}

func validINEGIEntityRow(row domain.INEGIEntityRow) bool {
	return row.GeoCode != "" &&
		row.StateCode != "" &&
		row.StateName != ""
}

func cleanINEGIMunicipalityRow(row domain.INEGIMunicipalityRow) domain.INEGIMunicipalityRow {
	row.GeoCode = strings.TrimSpace(row.GeoCode)
	row.StateCode = strings.TrimSpace(row.StateCode)
	row.StateName = normalizeSpaces(row.StateName)
	row.StateAbbreviation = normalizeSpaces(row.StateAbbreviation)
	row.MunicipalityCode = strings.TrimSpace(row.MunicipalityCode)
	row.MunicipalityName = normalizeSpaces(row.MunicipalityName)
	row.HeadCode = strings.TrimSpace(row.HeadCode)
	row.HeadName = normalizeSpaces(row.HeadName)
	return row
}

func validINEGIMunicipalityRow(row domain.INEGIMunicipalityRow) bool {
	return row.GeoCode != "" &&
		row.StateCode != "" &&
		row.StateName != "" &&
		row.MunicipalityCode != "" &&
		row.MunicipalityName != ""
}

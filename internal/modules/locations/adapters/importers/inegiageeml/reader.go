package inegiageeml

import (
	"archive/zip"
	"encoding/csv"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/domain"
)

func ReadEntitiesZip(path string) ([]domain.INEGIEntityRow, error) {
	return readZipCSV(path, mapRecordToEntity)
}

func ReadMunicipalitiesZip(path string) ([]domain.INEGIMunicipalityRow, error) {
	return readZipCSV(path, mapRecordToMunicipality)
}

func readZipCSV[T any](path string, mapper func(map[string]string) (T, error)) ([]T, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	file := selectCSVFile(reader.File)
	if file == nil {
		return nil, fmt.Errorf("zip does not contain a csv file: %s", path)
	}

	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return readCSV(rc, mapper)
}

func selectCSVFile(files []*zip.File) *zip.File {
	var fallback *zip.File
	for _, file := range files {
		name := strings.ToLower(filepath.Base(file.Name))
		if !strings.HasSuffix(name, ".csv") {
			continue
		}
		if strings.Contains(name, "_utf") {
			return file
		}
		if fallback == nil {
			fallback = file
		}
	}
	return fallback
}

func readCSV[T any](r io.Reader, mapper func(map[string]string) (T, error)) ([]T, error) {
	csvReader := csv.NewReader(r)
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true

	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	headers := normalizeHeaders(header)

	var rows []T
	for line := 2; ; line++ {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read line %d: %w", line, err)
		}
		values := recordMap(headers, record)
		row, err := mapper(values)
		if err != nil {
			return nil, fmt.Errorf("map line %d: %w", line, err)
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func normalizeHeaders(header []string) []string {
	headers := make([]string, len(header))
	for i, value := range header {
		headers[i] = strings.TrimPrefix(strings.TrimSpace(value), "\ufeff")
	}
	return headers
}

func recordMap(headers, record []string) map[string]string {
	values := make(map[string]string, len(headers))
	for i, header := range headers {
		if i >= len(record) {
			values[header] = ""
			continue
		}
		values[header] = strings.TrimSpace(record[i])
	}
	return values
}

func mapRecordToEntity(values map[string]string) (domain.INEGIEntityRow, error) {
	geoCode, err := require(values, "CVEGEO")
	if err != nil {
		return domain.INEGIEntityRow{}, err
	}
	stateCode, err := require(values, "CVE_ENT")
	if err != nil {
		return domain.INEGIEntityRow{}, err
	}
	stateName, err := require(values, "NOM_ENT")
	if err != nil {
		return domain.INEGIEntityRow{}, err
	}
	return domain.INEGIEntityRow{
		GeoCode:           geoCode,
		StateCode:         stateCode,
		StateName:         stateName,
		StateAbbreviation: values["NOM_ABR"],
	}, nil
}

func mapRecordToMunicipality(values map[string]string) (domain.INEGIMunicipalityRow, error) {
	geoCode, err := require(values, "CVEGEO")
	if err != nil {
		return domain.INEGIMunicipalityRow{}, err
	}
	stateCode, err := require(values, "CVE_ENT")
	if err != nil {
		return domain.INEGIMunicipalityRow{}, err
	}
	stateName, err := require(values, "NOM_ENT")
	if err != nil {
		return domain.INEGIMunicipalityRow{}, err
	}
	municipalityCode, err := require(values, "CVE_MUN")
	if err != nil {
		return domain.INEGIMunicipalityRow{}, err
	}
	municipalityName, err := require(values, "NOM_MUN")
	if err != nil {
		return domain.INEGIMunicipalityRow{}, err
	}
	return domain.INEGIMunicipalityRow{
		GeoCode:           geoCode,
		StateCode:         stateCode,
		StateName:         stateName,
		StateAbbreviation: values["NOM_ABR"],
		MunicipalityCode:  municipalityCode,
		MunicipalityName:  municipalityName,
		HeadCode:          values["CVE_CAB"],
		HeadName:          values["NOM_CAB"],
	}, nil
}

func require(values map[string]string, name string) (string, error) {
	value, ok := values[name]
	if !ok {
		return "", fmt.Errorf("missing column %s", name)
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("empty column %s", name)
	}
	return value, nil
}

package sepomex

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/domain"
)

func ReadSnapshot(r io.Reader) ([]domain.SEPOMEXRow, error) {
	buffered := bufio.NewReader(r)
	sample, _ := buffered.Peek(4096)

	reader := csv.NewReader(buffered)
	reader.Comma = detectDelimiter(string(sample))
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	header, err := readHeader(reader)
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}

	index := mapHeader(header)
	required := []string{"d_codigo", "d_asenta", "d_tipo_asenta", "d_mnpio", "d_estado", "c_estado", "c_mnpio", "id_asenta_cpcons"}
	for _, key := range required {
		if _, ok := index[key]; !ok {
			return nil, fmt.Errorf("missing required SEPOMEX column: %s", key)
		}
	}

	var rows []domain.SEPOMEXRow
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read record: %w", err)
		}

		row := domain.SEPOMEXRow{
			PostalCode:       value(record, index, "d_codigo"),
			SettlementName:   value(record, index, "d_asenta"),
			SettlementType:   value(record, index, "d_tipo_asenta"),
			Municipality:     value(record, index, "d_mnpio"),
			State:            value(record, index, "d_estado"),
			City:             value(record, index, "d_ciudad"),
			StateCode:        value(record, index, "c_estado"),
			MunicipalityCode: value(record, index, "c_mnpio"),
			SettlementCode:   value(record, index, "id_asenta_cpcons"),
			Zone:             value(record, index, "d_zona"),
			CityCode:         value(record, index, "c_cve_ciudad"),
		}
		if row.PostalCode == "" || row.SettlementName == "" {
			continue
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func readHeader(reader *csv.Reader) ([]string, error) {
	for {
		record, err := reader.Read()
		if err != nil {
			return nil, err
		}
		for _, field := range record {
			if canonicalKey(field) == "d_codigo" {
				return record, nil
			}
		}
	}
}

func detectDelimiter(sample string) rune {
	candidates := []rune{'|', ',', '\t'}
	best := ','
	bestCount := -1
	for _, candidate := range candidates {
		count := strings.Count(sample, string(candidate))
		if count > bestCount {
			best = candidate
			bestCount = count
		}
	}
	return best
}

func mapHeader(header []string) map[string]int {
	index := make(map[string]int, len(header))
	for i, name := range header {
		index[canonicalKey(name)] = i
	}
	return index
}

func value(record []string, index map[string]int, key string) string {
	i, ok := index[canonicalKey(key)]
	if !ok || i >= len(record) {
		return ""
	}
	return strings.Join(strings.Fields(strings.TrimSpace(record[i])), " ")
}

func canonicalKey(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

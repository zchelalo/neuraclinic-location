package sepomex

import (
	"strings"
	"testing"
)

func TestReadSnapshotParsesSmallMexicoSample(t *testing.T) {
	const snapshot = `d_codigo,d_asenta,d_tipo_asenta,D_mnpio,d_estado,d_ciudad,c_estado,c_oficina,c_CP,c_tipo_asenta,c_mnpio,id_asenta_cpcons,d_zona,c_cve_ciudad
83200,Centro,Colonia,Hermosillo,Sonora,Hermosillo,26,83201,,09,030,0001,Urbano,01
83220,San Benito,Colonia,Hermosillo,Sonora,Hermosillo,26,83201,,09,030,0002,Urbano,01
`

	rows, err := ReadSnapshot(strings.NewReader(snapshot))
	if err != nil {
		t.Fatalf("ReadSnapshot returned error: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0].PostalCode != "83200" || rows[0].StateCode != "26" || rows[0].MunicipalityCode != "030" {
		t.Fatalf("unexpected first row: %+v", rows[0])
	}
	if rows[1].SettlementName != "San Benito" {
		t.Fatalf("expected second settlement name San Benito, got %q", rows[1].SettlementName)
	}
}

func TestReadSnapshotParsesPipeDelimitedSEPOMEXFile(t *testing.T) {
	const snapshot = `d_codigo|d_asenta|d_tipo_asenta|D_mnpio|d_estado|d_ciudad|c_estado|c_oficina|c_CP|c_tipo_asenta|c_mnpio|id_asenta_cpcons|d_zona|c_cve_ciudad
83200|Centro|Colonia|Hermosillo|Sonora|Hermosillo|26|83201||09|030|0001|Urbano|01
`

	rows, err := ReadSnapshot(strings.NewReader(snapshot))
	if err != nil {
		t.Fatalf("ReadSnapshot returned error: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0].SettlementType != "Colonia" {
		t.Fatalf("expected settlement type Colonia, got %q", rows[0].SettlementType)
	}
}

func TestReadSnapshotSkipsLegalPreamble(t *testing.T) {
	const snapshot = `El Catalogo Nacional de Codigos Postales se proporciona en forma gratuita.
d_codigo|d_asenta|d_tipo_asenta|D_mnpio|d_estado|d_ciudad|c_estado|c_oficina|c_CP|c_tipo_asenta|c_mnpio|id_asenta_cpcons|d_zona|c_cve_ciudad
83200|Centro|Colonia|Hermosillo|Sonora|Hermosillo|26|83201||09|030|0001|Urbano|01
`

	rows, err := ReadSnapshot(strings.NewReader(snapshot))
	if err != nil {
		t.Fatalf("ReadSnapshot returned error: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	if rows[0].PostalCode != "83200" {
		t.Fatalf("expected postal code 83200, got %q", rows[0].PostalCode)
	}
}

func TestReadSnapshotRequiresKnownColumns(t *testing.T) {
	_, err := ReadSnapshot(strings.NewReader("postal_code,name\n83200,Centro\n"))
	if err == nil {
		t.Fatal("expected missing column error")
	}
}

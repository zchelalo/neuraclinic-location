package inegiageeml

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadEntitiesZipUsesUTFCSV(t *testing.T) {
	path := writeZip(t, map[string]string{
		"AGEEML.csv":     "bad\n",
		"AGEEML_utf.csv": "CVEGEO,CVE_ENT,NOM_ENT,NOM_ABR\n\"01\",\"01\",\"Aguascalientes\",\"Ags.\"\n",
	})

	rows, err := ReadEntitiesZip(path)
	if err != nil {
		t.Fatalf("ReadEntitiesZip() error = %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	if rows[0].StateCode != "01" || rows[0].StateName != "Aguascalientes" {
		t.Fatalf("row = %+v", rows[0])
	}
}

func TestReadMunicipalitiesZip(t *testing.T) {
	path := writeZip(t, map[string]string{
		"AGEEML_utf.csv": strings.Join([]string{
			"CVEGEO,CVE_ENT,NOM_ENT,NOM_ABR,CVE_MUN,NOM_MUN,CVE_CAB,NOM_CAB",
			"\"01001\",\"01\",\"Aguascalientes\",\"Ags.\",\"001\",\"Aguascalientes\",\"0001\",\"Aguascalientes\"",
		}, "\n"),
	})

	rows, err := ReadMunicipalitiesZip(path)
	if err != nil {
		t.Fatalf("ReadMunicipalitiesZip() error = %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	if rows[0].GeoCode != "01001" || rows[0].MunicipalityCode != "001" {
		t.Fatalf("row = %+v", rows[0])
	}
}

func TestReadEntitiesZipRejectsMissingColumns(t *testing.T) {
	path := writeZip(t, map[string]string{
		"AGEEML_utf.csv": "CVEGEO,CVE_ENT\n\"01\",\"01\"\n",
	})

	if _, err := ReadEntitiesZip(path); err == nil {
		t.Fatal("ReadEntitiesZip() error = nil, want error")
	}
}

func writeZip(t *testing.T, files map[string]string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "catalog.zip")
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create zip: %v", err)
	}

	zipWriter := zip.NewWriter(file)
	for name, contents := range files {
		w, err := zipWriter.Create(name)
		if err != nil {
			t.Fatalf("create zip member: %v", err)
		}
		if _, err := w.Write([]byte(contents)); err != nil {
			t.Fatalf("write zip member: %v", err)
		}
	}
	if err := zipWriter.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close zip file: %v", err)
	}

	return path
}

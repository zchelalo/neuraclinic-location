package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"

	inegireader "github.com/zchelalo/neuraclinic-location/internal/modules/locations/adapters/importers/inegiageeml"
	sepomexreader "github.com/zchelalo/neuraclinic-location/internal/modules/locations/adapters/importers/sepomex"
	locationpg "github.com/zchelalo/neuraclinic-location/internal/modules/locations/adapters/persistence/postgres"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/application/importinegiageeml"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/application/importsepomex"
	"github.com/zchelalo/neuraclinic-location/internal/modules/locations/domain"
	"github.com/zchelalo/neuraclinic-location/pkg/bootstrap"
)

func main() {
	if len(os.Args) < 2 {
		usageAndExit()
	}

	switch os.Args[1] {
	case "sepomex":
		if err := runSEPOMEX(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "location-import sepomex: %v\n", err)
			os.Exit(1)
		}
	case "inegi-ageeml":
		if err := runINEGIAgeeml(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "location-import inegi-ageeml: %v\n", err)
			os.Exit(1)
		}
	default:
		usageAndExit()
	}
}

func runSEPOMEX(args []string) error {
	fs := flag.NewFlagSet("sepomex", flag.ExitOnError)
	filePath := fs.String("file", "", "Path to SEPOMEX CSV/TXT snapshot")
	sourceVersion := fs.String("source-version", "", "Snapshot version, for example 2026-06")
	sourceURL := fs.String("source-url", "", "Snapshot or source URL")
	sourceLicense := fs.String("license", "Datos abiertos; validar terminos de Correos de Mexico / SEPOMEX antes de redistribuir", "License note")
	sourceAttribution := fs.String("attribution", "Correos de Mexico / SEPOMEX", "Required source attribution")
	encoding := fs.String("encoding", "utf-8", "Input encoding: utf-8 or latin1")
	dryRun := fs.Bool("dry-run", false, "Read and validate the snapshot without writing to the database")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(*filePath) == "" {
		return fmt.Errorf("--file is required")
	}
	if strings.TrimSpace(*sourceVersion) == "" {
		return fmt.Errorf("--source-version is required")
	}

	file, err := os.Open(*filePath)
	if err != nil {
		return fmt.Errorf("open snapshot: %w", err)
	}
	defer file.Close()

	rows, err := sepomexreader.ReadSnapshot(decodeReader(file, *encoding))
	if err != nil {
		return fmt.Errorf("read snapshot: %w", err)
	}
	if *dryRun {
		fmt.Printf("SEPOMEX dry run complete\n")
		fmt.Printf("rows_read=%d\n", len(rows))
		return nil
	}

	cfg, err := bootstrap.LoadConfig(".env")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	db, err := bootstrap.NewDB(context.Background(), cfg)
	if err != nil {
		return fmt.Errorf("connect db: %w", err)
	}
	defer db.Close()

	repo := locationpg.NewRepository(db)
	useCase := importsepomex.New(repo)
	summary, err := useCase.Execute(context.Background(), importsepomex.Command{
		Source: domain.DataSource{
			Key:         "sepomex",
			Name:        "SEPOMEX / Correos de Mexico",
			Version:     *sourceVersion,
			License:     *sourceLicense,
			Attribution: *sourceAttribution,
			URL:         *sourceURL,
			CountryCode: "MX",
			CountryName: "Mexico",
		},
		Rows: rows,
	})
	if err != nil {
		return fmt.Errorf("import snapshot: %w", err)
	}

	fmt.Printf("SEPOMEX import complete\n")
	fmt.Printf("rows_read=%d rows_imported=%d countries=%d admin_areas=%d localities=%d postal_codes=%d settlements=%d\n",
		summary.RowsRead,
		summary.RowsImported,
		summary.Countries,
		summary.AdminAreas,
		summary.Localities,
		summary.PostalCodes,
		summary.Settlements,
	)
	return nil
}

func runINEGIAgeeml(args []string) error {
	fs := flag.NewFlagSet("inegi-ageeml", flag.ExitOnError)
	entitiesPath := fs.String("entities-file", "", "Path to INEGI AGEEML entidades zip")
	municipalitiesPath := fs.String("municipalities-file", "", "Path to INEGI AGEEML municipios zip")
	sourceVersion := fs.String("source-version", "", "Snapshot version, for example 2026-06-12")
	sourceURL := fs.String("source-url", "", "Snapshot or source URL")
	sourceLicense := fs.String("license", "Terminos de libre uso INEGI; validar atribucion antes de redistribuir", "License note")
	sourceAttribution := fs.String("attribution", "INEGI", "Required source attribution")
	dryRun := fs.Bool("dry-run", false, "Read and validate the snapshots without writing to the database")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(*entitiesPath) == "" {
		return fmt.Errorf("--entities-file is required")
	}
	if strings.TrimSpace(*municipalitiesPath) == "" {
		return fmt.Errorf("--municipalities-file is required")
	}
	if strings.TrimSpace(*sourceVersion) == "" {
		return fmt.Errorf("--source-version is required")
	}

	entities, err := inegireader.ReadEntitiesZip(*entitiesPath)
	if err != nil {
		return fmt.Errorf("read entities: %w", err)
	}
	municipalities, err := inegireader.ReadMunicipalitiesZip(*municipalitiesPath)
	if err != nil {
		return fmt.Errorf("read municipalities: %w", err)
	}
	if *dryRun {
		fmt.Printf("INEGI AGEEML dry run complete\n")
		fmt.Printf("entities_read=%d municipalities_read=%d rows_read=%d\n", len(entities), len(municipalities), len(entities)+len(municipalities))
		return nil
	}

	cfg, err := bootstrap.LoadConfig(".env")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	db, err := bootstrap.NewDB(context.Background(), cfg)
	if err != nil {
		return fmt.Errorf("connect db: %w", err)
	}
	defer db.Close()

	repo := locationpg.NewRepository(db)
	useCase := importinegiageeml.New(repo)
	summary, err := useCase.Execute(context.Background(), importinegiageeml.Command{
		Source: domain.DataSource{
			Key:         "inegi-ageeml",
			Name:        "INEGI AGEEML",
			Version:     *sourceVersion,
			License:     *sourceLicense,
			Attribution: *sourceAttribution,
			URL:         *sourceURL,
			CountryCode: "MX",
			CountryName: "Mexico",
		},
		Entities:       entities,
		Municipalities: municipalities,
	})
	if err != nil {
		return fmt.Errorf("import snapshots: %w", err)
	}

	fmt.Printf("INEGI AGEEML import complete\n")
	fmt.Printf("rows_read=%d rows_imported=%d countries=%d admin_areas=%d localities=%d postal_codes=%d settlements=%d\n",
		summary.RowsRead,
		summary.RowsImported,
		summary.Countries,
		summary.AdminAreas,
		summary.Localities,
		summary.PostalCodes,
		summary.Settlements,
	)
	return nil
}

func decodeReader(r io.Reader, encoding string) io.Reader {
	switch strings.ToLower(strings.TrimSpace(encoding)) {
	case "latin1", "iso-8859-1", "iso8859-1":
		return charmap.ISO8859_1.NewDecoder().Reader(r)
	default:
		return r
	}
}

func usageAndExit() {
	fmt.Fprintln(os.Stderr, "usage: location-import <sepomex|inegi-ageeml> [flags]")
	os.Exit(2)
}

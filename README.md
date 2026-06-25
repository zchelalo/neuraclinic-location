# Neuraclinic Location Microservice

Go gRPC catalog service for address suggestions and normalization. The v1 scope is Mexico-first, offline, and backed by PostgreSQL.

## Local Setup

Run from `neuraclinic-location`:

```bash
make create-envs
make tls-generate-dev
make compose-build
```

The service listens inside Docker on `:8000` and is exposed on host port `8005`.

## Useful Commands

```bash
make proto
make sqlc
make test
make build
make migrate-up
make download-location-data
make import-location-data
make import-sepomex
make import-inegi-ageeml
make compose
make compose-down
```

## Catalog Data

Starting the service only creates the database schema. It does not download or import catalogs automatically. Until a snapshot is imported, the gRPC endpoints respond successfully but return empty result sets.

Use an explicit import step for local development and deployments:

```bash
make download-location-data
make import-location-data \
  SOURCE_VERSION=2026-06-22 \
  INEGI_SOURCE_VERSION=2026-06-12
```

`download-location-data` downloads the SEPOMEX national postal-code snapshot and the INEGI AGEEML zips. `import-location-data` runs both importers in order: SEPOMEX first, then INEGI. Each target starts/builds the local service stack, waits for the Postgres-backed service to run migrations, mounts the snapshots into one-off containers, and executes `location-import`.

By default `import-sepomex` reads `data/sources/sepomex/CPdescarga.txt`, the path produced by `download-sepomex`.

Optional variables:

- `SEPOMEX_FILE`: path to a custom SEPOMEX snapshot.
- `SOURCE_VERSION`: snapshot version stored in `data_sources`; defaults to `2026-06-22`.
- `SOURCE_LICENSE`: license note stored in `data_sources`.
- `SOURCE_ATTRIBUTION`: attribution stored in `data_sources`.
- `SOURCE_URL`: source URL stored in `data_sources`.
- `SOURCE_ENCODING`: `latin1` by default for the current SEPOMEX file; use `utf-8` for converted snapshots.

The SEPOMEX importer is idempotent. It uses stable IDs for Mexico, states, municipalities, postal codes, and settlements, so rerunning the same snapshot updates records instead of duplicating them. SEPOMEX municipalities are stored as `localities` with type `municipality`; states are stored as `admin_areas` with type `state`; settlements are colonias/asentamientos linked to postal codes.

INEGI AGEEML can also be imported independently:

```bash
make download-inegi-ageeml
make import-inegi-ageeml \
  INEGI_SOURCE_VERSION=2026-06-12
```

By default `import-inegi-ageeml` reads:

- `data/sources/inegi-ageeml/catun_entidad.zip`
- `data/sources/inegi-ageeml/catun_municipio.zip`

Optional INEGI variables:

- `INEGI_ENTITIES_FILE`: path to a custom `catun_entidad.zip`.
- `INEGI_MUNICIPALITIES_FILE`: path to a custom `catun_municipio.zip`.
- `INEGI_SOURCE_VERSION`: snapshot version stored in `data_sources`; defaults to `2026-06-12`.
- `INEGI_SOURCE_LICENSE`: license note stored in `data_sources`.
- `INEGI_SOURCE_ATTRIBUTION`: attribution stored in `data_sources`.
- `INEGI_SOURCE_URL`: source URL stored in `data_sources`.

The INEGI importer is idempotent. It stores states and municipalities as `admin_areas`, with municipalities parented to their state and using the INEGI `CVEGEO` as code.

The download target also stores this larger file for the next locality-level importer:

- `data/sources/inegi-ageeml/may_acento.zip`

That file is intentionally not imported by the current target. All downloaded files are ignored by git.

## Location Types

`ListAdminAreas` should be filtered with the `AdminAreaType` enum in the proto contract:

- `ADMIN_AREA_TYPE_STATE`: first-level administrative areas, such as Mexican states.
- `ADMIN_AREA_TYPE_MUNICIPALITY`: municipalities parented to a state.

The previous request field `type` is still accepted for compatibility but is deprecated. New clients should use `admin_area_type`; invalid deprecated string values are rejected.

Settlement types are intentionally not an enum yet. Values such as `Colonia`, `Fraccionamiento`, or `Pueblo` come from SEPOMEX and may vary by source, so they remain source data rather than application-controlled taxonomy.

## Data Sources

The schema tracks `source`, `source_version`, and source attribution on every result. Production should import versioned snapshots rather than calling public geocoding APIs at runtime.

- INEGI: administrative geography and locality catalogs. Include INEGI attribution and terms of free use in packaged data.
- SEPOMEX / Correos de Mexico: postal codes and settlements. Import as a versioned snapshot.
- OpenStreetMap: optional streets/free-text enrichment. Document ODbL attribution before packaging derived data.

The patient record service should continue storing editable structured text fields. Location IDs from this service are suggestions, not mandatory validation constraints.

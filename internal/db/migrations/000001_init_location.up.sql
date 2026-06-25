CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE OR REPLACE FUNCTION public.immutable_unaccent(input text)
RETURNS text
LANGUAGE sql
IMMUTABLE
PARALLEL SAFE
AS $$
  SELECT public.unaccent(input);
$$;

CREATE OR REPLACE FUNCTION public.normalize_location_text(input text)
RETURNS text
LANGUAGE sql
IMMUTABLE
PARALLEL SAFE
AS $$
  SELECT trim(regexp_replace(lower(public.immutable_unaccent(coalesce(input, ''))), '[^a-z0-9]+', ' ', 'g'));
$$;

CREATE TABLE data_sources (
  id uuid PRIMARY KEY,
  key varchar(50) NOT NULL UNIQUE,
  name text NOT NULL,
  version varchar(100) NOT NULL,
  license text NOT NULL,
  attribution text NOT NULL,
  url text,
  imported_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE countries (
  id uuid PRIMARY KEY,
  country_code char(2) NOT NULL UNIQUE,
  name text NOT NULL,
  source_id uuid NOT NULL REFERENCES data_sources(id) ON DELETE RESTRICT,
  source_record_id varchar(150),
  source_version varchar(100) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  CHECK (country_code = upper(country_code))
);

CREATE TABLE admin_areas (
  id uuid PRIMARY KEY,
  country_code char(2) NOT NULL REFERENCES countries(country_code) ON DELETE RESTRICT,
  parent_id uuid REFERENCES admin_areas(id) ON DELETE RESTRICT,
  code varchar(50) NOT NULL,
  name text NOT NULL,
  type varchar(50) NOT NULL,
  source_id uuid NOT NULL REFERENCES data_sources(id) ON DELETE RESTRICT,
  source_record_id varchar(150),
  source_version varchar(100) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (country_code, type, code, source_id, source_version)
);

CREATE TABLE localities (
  id uuid PRIMARY KEY,
  country_code char(2) NOT NULL REFERENCES countries(country_code) ON DELETE RESTRICT,
  admin_area_id uuid NOT NULL REFERENCES admin_areas(id) ON DELETE RESTRICT,
  code varchar(50) NOT NULL,
  name text NOT NULL,
  type varchar(50) NOT NULL DEFAULT 'locality',
  source_id uuid NOT NULL REFERENCES data_sources(id) ON DELETE RESTRICT,
  source_record_id varchar(150),
  source_version varchar(100) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (country_code, admin_area_id, code, source_id, source_version)
);

CREATE TABLE postal_codes (
  id uuid PRIMARY KEY,
  country_code char(2) NOT NULL REFERENCES countries(country_code) ON DELETE RESTRICT,
  admin_area_id uuid NOT NULL REFERENCES admin_areas(id) ON DELETE RESTRICT,
  locality_id uuid REFERENCES localities(id) ON DELETE RESTRICT,
  postal_code varchar(12) NOT NULL,
  source_id uuid NOT NULL REFERENCES data_sources(id) ON DELETE RESTRICT,
  source_record_id varchar(150),
  source_version varchar(100) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  UNIQUE (country_code, postal_code, admin_area_id, locality_id, source_id, source_version)
);

CREATE TABLE settlements (
  id uuid PRIMARY KEY,
  country_code char(2) NOT NULL REFERENCES countries(country_code) ON DELETE RESTRICT,
  admin_area_id uuid NOT NULL REFERENCES admin_areas(id) ON DELETE RESTRICT,
  locality_id uuid REFERENCES localities(id) ON DELETE RESTRICT,
  postal_code_id uuid REFERENCES postal_codes(id) ON DELETE RESTRICT,
  code varchar(50),
  name text NOT NULL,
  settlement_type text NOT NULL,
  source_id uuid NOT NULL REFERENCES data_sources(id) ON DELETE RESTRICT,
  source_record_id varchar(150),
  source_version varchar(100) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE streets (
  id uuid PRIMARY KEY,
  country_code char(2) NOT NULL REFERENCES countries(country_code) ON DELETE RESTRICT,
  admin_area_id uuid NOT NULL REFERENCES admin_areas(id) ON DELETE RESTRICT,
  locality_id uuid REFERENCES localities(id) ON DELETE RESTRICT,
  settlement_id uuid REFERENCES settlements(id) ON DELETE RESTRICT,
  name text NOT NULL,
  source_id uuid NOT NULL REFERENCES data_sources(id) ON DELETE RESTRICT,
  source_record_id varchar(150),
  source_version varchar(100) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_countries_name_trgm
  ON countries USING gin (public.normalize_location_text(name) gin_trgm_ops);

CREATE INDEX idx_admin_areas_country_parent_type
  ON admin_areas (country_code, parent_id, type, code);

CREATE INDEX idx_admin_areas_name_trgm
  ON admin_areas USING gin (public.normalize_location_text(name) gin_trgm_ops);

CREATE INDEX idx_localities_admin_code
  ON localities (country_code, admin_area_id, code);

CREATE INDEX idx_localities_name_trgm
  ON localities USING gin (public.normalize_location_text(name) gin_trgm_ops);

CREATE INDEX idx_postal_codes_country_postal
  ON postal_codes (country_code, postal_code text_pattern_ops);

CREATE INDEX idx_settlements_postal
  ON settlements (country_code, postal_code_id);

CREATE INDEX idx_settlements_name_trgm
  ON settlements USING gin (public.normalize_location_text(name) gin_trgm_ops);

CREATE INDEX idx_streets_settlement
  ON streets (country_code, settlement_id);

CREATE INDEX idx_streets_name_trgm
  ON streets USING gin (public.normalize_location_text(name) gin_trgm_ops);

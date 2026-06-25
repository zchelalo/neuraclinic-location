DROP TABLE IF EXISTS streets;
DROP TABLE IF EXISTS settlements;
DROP TABLE IF EXISTS postal_codes;
DROP TABLE IF EXISTS localities;
DROP TABLE IF EXISTS admin_areas;
DROP TABLE IF EXISTS countries;
DROP TABLE IF EXISTS data_sources;

DROP FUNCTION IF EXISTS public.normalize_location_text(text);
DROP FUNCTION IF EXISTS public.immutable_unaccent(text);

-- name: ListCountries :many
SELECT
  c.country_code::text AS country_code,
  c.name,
  c.name AS label,
  ds.key AS source,
  c.source_version,
  CASE
    WHEN sqlc.arg(query)::text = '' THEN 0::float8
    ELSE similarity(public.normalize_location_text(c.name), public.normalize_location_text(sqlc.arg(query)::text))::float8
  END AS score
FROM countries c
JOIN data_sources ds ON ds.id = c.source_id
WHERE (
    sqlc.arg(query)::text = ''
    OR public.normalize_location_text(c.name) LIKE '%' || public.normalize_location_text(sqlc.arg(query)::text) || '%'
    OR similarity(public.normalize_location_text(c.name), public.normalize_location_text(sqlc.arg(query)::text)) > 0.2
  )
ORDER BY score DESC, c.name ASC
LIMIT sqlc.arg(limit_count)::int;

-- name: ListAdminAreas :many
SELECT
  aa.id::text AS id,
  aa.country_code::text AS country_code,
  aa.code,
  aa.name,
  aa.type,
  parent.code AS parent_code,
  concat_ws(', ', aa.name, upper(aa.country_code::text)) AS label,
  ds.key AS source,
  aa.source_version,
  CASE
    WHEN sqlc.arg(query)::text = '' THEN 0::float8
    ELSE similarity(public.normalize_location_text(aa.name), public.normalize_location_text(sqlc.arg(query)::text))::float8
  END AS score
FROM admin_areas aa
LEFT JOIN admin_areas parent ON parent.id = aa.parent_id
JOIN data_sources ds ON ds.id = aa.source_id
WHERE aa.country_code = upper(sqlc.arg(country_code)::text)
  AND (sqlc.arg(parent_code)::text = '' OR parent.code = sqlc.arg(parent_code)::text)
  AND (sqlc.arg(type_filter)::text = '' OR aa.type = sqlc.arg(type_filter)::text)
  AND (
    sqlc.arg(query)::text = ''
    OR public.normalize_location_text(aa.name) LIKE '%' || public.normalize_location_text(sqlc.arg(query)::text) || '%'
    OR similarity(public.normalize_location_text(aa.name), public.normalize_location_text(sqlc.arg(query)::text)) > 0.2
  )
ORDER BY score DESC, aa.name ASC
LIMIT sqlc.arg(limit_count)::int;

-- name: ListLocalities :many
SELECT
  l.id::text AS id,
  l.country_code::text AS country_code,
  aa.code AS admin_area_code,
  l.code,
  l.name,
  l.type,
  concat_ws(', ', l.name, aa.name, upper(l.country_code::text)) AS label,
  ds.key AS source,
  l.source_version,
  CASE
    WHEN sqlc.arg(query)::text = '' THEN 0::float8
    ELSE similarity(public.normalize_location_text(l.name), public.normalize_location_text(sqlc.arg(query)::text))::float8
  END AS score
FROM localities l
JOIN admin_areas aa ON aa.id = l.admin_area_id
JOIN data_sources ds ON ds.id = l.source_id
WHERE l.country_code = upper(sqlc.arg(country_code)::text)
  AND (sqlc.arg(admin_area_code)::text = '' OR aa.code = sqlc.arg(admin_area_code)::text)
  AND (
    sqlc.arg(query)::text = ''
    OR public.normalize_location_text(l.name) LIKE '%' || public.normalize_location_text(sqlc.arg(query)::text) || '%'
    OR similarity(public.normalize_location_text(l.name), public.normalize_location_text(sqlc.arg(query)::text)) > 0.2
  )
ORDER BY score DESC, l.name ASC
LIMIT sqlc.arg(limit_count)::int;

-- name: ListSettlements :many
SELECT
  s.id::text AS id,
  s.country_code::text AS country_code,
  aa.code AS admin_area_code,
  l.code AS locality_code,
  pc.postal_code,
  s.name,
  s.settlement_type AS type,
  concat_ws(', ', s.name, s.settlement_type, pc.postal_code, l.name, aa.name) AS label,
  ds.key AS source,
  s.source_version,
  CASE
    WHEN sqlc.arg(query)::text = '' THEN 0::float8
    ELSE similarity(public.normalize_location_text(s.name), public.normalize_location_text(sqlc.arg(query)::text))::float8
  END AS score
FROM settlements s
JOIN admin_areas aa ON aa.id = s.admin_area_id
LEFT JOIN localities l ON l.id = s.locality_id
LEFT JOIN postal_codes pc ON pc.id = s.postal_code_id
JOIN data_sources ds ON ds.id = s.source_id
WHERE s.country_code = upper(sqlc.arg(country_code)::text)
  AND (sqlc.arg(admin_area_code)::text = '' OR aa.code = sqlc.arg(admin_area_code)::text)
  AND (sqlc.arg(locality_code)::text = '' OR l.code = sqlc.arg(locality_code)::text)
  AND (sqlc.arg(postal_code)::text = '' OR pc.postal_code = sqlc.arg(postal_code)::text)
  AND (
    sqlc.arg(query)::text = ''
    OR public.normalize_location_text(s.name) LIKE '%' || public.normalize_location_text(sqlc.arg(query)::text) || '%'
    OR similarity(public.normalize_location_text(s.name), public.normalize_location_text(sqlc.arg(query)::text)) > 0.2
  )
ORDER BY score DESC, s.name ASC
LIMIT sqlc.arg(limit_count)::int;

-- name: SearchPostalCodes :many
SELECT
  pc.postal_code,
  concat_ws(', ', s.name, pc.postal_code, l.name, aa.name, c.name) AS label,
  pc.country_code::text AS country_code,
  c.name AS country_name,
  aa.code AS admin_area_code,
  aa.name AS admin_area_name,
  l.code AS locality_code,
  l.name AS locality_name,
  s.name AS settlement_name,
  s.settlement_type AS settlement_type,
  ds.key AS source,
  pc.source_version,
  CASE
    WHEN pc.postal_code = sqlc.arg(postal_code_prefix)::text THEN 1::float8
    ELSE 0.75::float8
  END AS score
FROM postal_codes pc
JOIN countries c ON c.country_code = pc.country_code
JOIN admin_areas aa ON aa.id = pc.admin_area_id
LEFT JOIN localities l ON l.id = pc.locality_id
LEFT JOIN settlements s ON s.postal_code_id = pc.id
JOIN data_sources ds ON ds.id = pc.source_id
WHERE pc.country_code = upper(sqlc.arg(country_code)::text)
  AND pc.postal_code LIKE sqlc.arg(postal_code_prefix)::text || '%'
ORDER BY score DESC, pc.postal_code ASC, s.name ASC NULLS LAST
LIMIT sqlc.arg(limit_count)::int;

-- name: SuggestAddresses :many
SELECT
  concat_ws(', ', st.name, s.name, pc.postal_code, l.name, aa.name, c.name) AS label,
  s.country_code::text AS country_code,
  c.name AS country_name,
  aa.code AS admin_area_code,
  aa.name AS admin_area_name,
  l.code AS locality_code,
  l.name AS locality_name,
  pc.postal_code,
  s.name AS settlement_name,
  s.settlement_type AS settlement_type,
  st.name AS street_name,
  ds.key AS source,
  s.source_version,
  CASE
    WHEN sqlc.arg(query)::text = '' THEN 0::float8
    ELSE GREATEST(
      similarity(public.normalize_location_text(s.name), public.normalize_location_text(sqlc.arg(query)::text)),
      similarity(public.normalize_location_text(coalesce(st.name, '')), public.normalize_location_text(sqlc.arg(query)::text)),
      similarity(public.normalize_location_text(coalesce(l.name, '')), public.normalize_location_text(sqlc.arg(query)::text)),
      similarity(public.normalize_location_text(aa.name), public.normalize_location_text(sqlc.arg(query)::text))
    )::float8
  END AS score
FROM settlements s
JOIN countries c ON c.country_code = s.country_code
JOIN admin_areas aa ON aa.id = s.admin_area_id
LEFT JOIN localities l ON l.id = s.locality_id
LEFT JOIN postal_codes pc ON pc.id = s.postal_code_id
LEFT JOIN streets st ON st.settlement_id = s.id
JOIN data_sources ds ON ds.id = s.source_id
WHERE s.country_code = upper(sqlc.arg(country_code)::text)
  AND (sqlc.arg(postal_code)::text = '' OR pc.postal_code = sqlc.arg(postal_code)::text)
  AND (
    sqlc.arg(query)::text = ''
    OR public.normalize_location_text(s.name) LIKE '%' || public.normalize_location_text(sqlc.arg(query)::text) || '%'
    OR public.normalize_location_text(coalesce(st.name, '')) LIKE '%' || public.normalize_location_text(sqlc.arg(query)::text) || '%'
    OR public.normalize_location_text(coalesce(l.name, '')) LIKE '%' || public.normalize_location_text(sqlc.arg(query)::text) || '%'
    OR public.normalize_location_text(aa.name) LIKE '%' || public.normalize_location_text(sqlc.arg(query)::text) || '%'
    OR similarity(public.normalize_location_text(s.name), public.normalize_location_text(sqlc.arg(query)::text)) > 0.2
    OR similarity(public.normalize_location_text(coalesce(st.name, '')), public.normalize_location_text(sqlc.arg(query)::text)) > 0.2
    OR similarity(public.normalize_location_text(coalesce(l.name, '')), public.normalize_location_text(sqlc.arg(query)::text)) > 0.2
    OR similarity(public.normalize_location_text(aa.name), public.normalize_location_text(sqlc.arg(query)::text)) > 0.2
  )
ORDER BY score DESC, label ASC
LIMIT sqlc.arg(limit_count)::int;

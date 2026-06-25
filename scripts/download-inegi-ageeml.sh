#!/usr/bin/env bash

set -euo pipefail

OUT_DIR="${1:-./data/sources/inegi-ageeml}"
BASE_URL="https://www.inegi.org.mx/contenidos/app/ageeml"

if ! command -v curl >/dev/null 2>&1; then
	echo "curl is required" >&2
	exit 1
fi

mkdir -p "$OUT_DIR"

curl -fL "$BASE_URL/catun_entidad.zip" -o "$OUT_DIR/catun_entidad.zip"
curl -fL "$BASE_URL/catun_municipio.zip" -o "$OUT_DIR/catun_municipio.zip"
curl -fL "$BASE_URL/may_acento.zip" -o "$OUT_DIR/may_acento.zip"

echo "Downloaded INEGI AGEEML snapshots:"
echo "  $OUT_DIR/catun_entidad.zip"
echo "  $OUT_DIR/catun_municipio.zip"
echo "  $OUT_DIR/may_acento.zip"

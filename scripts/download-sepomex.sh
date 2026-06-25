#!/usr/bin/env bash

set -euo pipefail

OUT_DIR="${1:-./data/sources/sepomex}"
URL="https://www.correosdemexico.gob.mx/SSLServicios/ConsultaCP/CodigoPostal_Exportar.aspx"
TMP_DIR="$(mktemp -d)"

cleanup() {
	rm -rf "$TMP_DIR"
}
trap cleanup EXIT

require_command() {
	if ! command -v "$1" >/dev/null 2>&1; then
		echo "$1 is required" >&2
		exit 1
	fi
}

require_command curl
require_command python3
require_command unzip

mkdir -p "$OUT_DIR"

curl -fsSL "$URL" -o "$TMP_DIR/form.html"

python3 - "$TMP_DIR/form.html" "$TMP_DIR/post-body.txt" <<'PY'
from html.parser import HTMLParser
from pathlib import Path
from urllib.parse import urlencode
import sys

class HiddenParser(HTMLParser):
    def __init__(self):
        super().__init__()
        self.values = {}

    def handle_starttag(self, tag, attrs):
        if tag.lower() != "input":
            return
        attr = dict(attrs)
        name = attr.get("name") or attr.get("id")
        if name in {"__VIEWSTATE", "__VIEWSTATEGENERATOR", "__EVENTVALIDATION"}:
            self.values[name] = attr.get("value", "")

html = Path(sys.argv[1]).read_text(encoding="iso-8859-1")
parser = HiddenParser()
parser.feed(html)

missing = {"__VIEWSTATE", "__VIEWSTATEGENERATOR", "__EVENTVALIDATION"} - parser.values.keys()
if missing:
    raise SystemExit(f"missing WebForms fields: {', '.join(sorted(missing))}")

fields = {
    **parser.values,
    "cboEdo": "00",
    "rblTipo": "txt",
    "btnDescarga.x": "12",
    "btnDescarga.y": "8",
}
Path(sys.argv[2]).write_text(urlencode(fields), encoding="ascii")
PY

curl -fsSL \
	-H "Content-Type: application/x-www-form-urlencoded" \
	-H "Referer: $URL" \
	--data-binary "@$TMP_DIR/post-body.txt" \
	-o "$OUT_DIR/CPdescarga-2026-06-22.zip" \
	"$URL"

unzip -o "$OUT_DIR/CPdescarga-2026-06-22.zip" CPdescarga.txt -d "$OUT_DIR" >/dev/null

echo "Downloaded SEPOMEX snapshot:"
echo "  $OUT_DIR/CPdescarga-2026-06-22.zip"
echo "  $OUT_DIR/CPdescarga.txt"

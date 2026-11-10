#!/bin/sh
# Сборка релизных архивов argus без GoReleaser.
#
#   ./scripts/build-release.sh v0.1.0-alpha.1
#
# Результат в dist/:
#   argus_0.1.0-alpha.1_linux_amd64.tar.gz
#   argus_0.1.0-alpha.1_darwin_arm64.tar.gz
#   argus_0.1.0-alpha.1_windows_amd64.zip
#   ...
#   checksums.txt
#
# Имена совпадают с тем, что ожидает scripts/install.sh.

set -eu

BINARY="argus"
MODULE="argus-agent"            # путь модуля из go.mod
PLATFORMS="linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64"

version="${1:-}"
[ -n "$version" ] || { echo "usage: $0 vX.Y.Z[-alpha.N]" >&2; exit 1; }
case "$version" in v*) ;; *) version="v$version" ;; esac
plain_version="${version#v}"

cd "$(dirname "$0")/.."
root="$(pwd)"
dist="$root/dist"

rm -rf "$dist"
mkdir -p "$dist"

ldflags="-s -w -X ${MODULE}/cmd.version=${plain_version}"

for p in $PLATFORMS; do
  os="${p%/*}"
  arch="${p#*/}"
  name="${BINARY}_${plain_version}_${os}_${arch}"
  stage="$dist/.stage/$name"
  mkdir -p "$stage"

  exe="$BINARY"
  [ "$os" = "windows" ] && exe="${BINARY}.exe"

  echo "==> $os/$arch"
  CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" \
    go build -trimpath -ldflags "$ldflags" -o "$stage/$exe" .

  cp README.md LICENSE "$stage/"

  if [ "$os" = "windows" ]; then
    command -v zip >/dev/null 2>&1 || { echo "error: нужен zip для windows-архивов" >&2; exit 1; }
    (cd "$stage" && zip -q "$dist/${name}.zip" "$exe" README.md LICENSE)
  else
    tar -czf "$dist/${name}.tar.gz" -C "$stage" "$exe" README.md LICENSE
  fi
done

rm -rf "$dist/.stage"

cd "$dist"
if command -v sha256sum >/dev/null 2>&1; then
  sha256sum ${BINARY}_* > checksums.txt
else
  shasum -a 256 ${BINARY}_* > checksums.txt
fi

echo
echo "Готово:"
ls -1 "$dist"

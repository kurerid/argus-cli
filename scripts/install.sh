#!/bin/sh
# Установщик argus.
#
#   curl -fsSL https://raw.githubusercontent.com/kurerid/argus-cli/main/install.sh | sh
#
# Переменные окружения (необязательные):
#   VERSION=v1.2.3        конкретная версия (по умолчанию последняя)
#   INSTALL_DIR=~/bin     куда ставить (по умолчанию /usr/local/bin или ~/.local/bin)
#
# Ожидаемые имена файлов в релизе (дефолт GoReleaser):
#   argus_1.2.3_linux_amd64.tar.gz
#   argus_1.2.3_darwin_arm64.tar.gz
#   checksums.txt

set -eu

# ---- Настрой под себя -------------------------------------------------------
REPO="kurerid/argus-cli"   # GitHub-репозиторий
BINARY="argus"     # имя бинарника и префикс архива
# -----------------------------------------------------------------------------

info()  { printf '\033[1;34m==>\033[0m %s\n' "$*"; }
warn()  { printf '\033[1;33mwarning:\033[0m %s\n' "$*" >&2; }
fail()  { printf '\033[1;31merror:\033[0m %s\n' "$*" >&2; exit 1; }
has()   { command -v "$1" >/dev/null 2>&1; }

# Скачивание: curl или wget
download() { # url dest
  if has curl; then
    curl -fsSL "$1" -o "$2"
  elif has wget; then
    wget -qO "$2" "$1"
  else
    fail "нужен curl или wget"
  fi
}

fetch() { # url -> stdout
  if has curl; then
    curl -fsSL "$1"
  else
    wget -qO- "$1"
  fi
}

detect_os() {
  case "$(uname -s)" in
    Linux)  echo linux ;;
    Darwin) echo darwin ;;
    MINGW*|MSYS*|CYGWIN*)
      fail "Windows: используй install.ps1, Scoop или скачай .zip из релиза" ;;
    *) fail "неподдерживаемая ОС: $(uname -s)" ;;
  esac
}

detect_arch() {
  arch="$(uname -m)"
  # Rosetta: на Apple Silicon терминал может притворяться x86_64
  if [ "$arch" = "x86_64" ] && [ "$(uname -s)" = "Darwin" ] \
     && [ "$(sysctl -n sysctl.proc_translated 2>/dev/null || echo 0)" = "1" ]; then
    arch="arm64"
  fi
  case "$arch" in
    x86_64|amd64)   echo amd64 ;;
    aarch64|arm64)  echo arm64 ;;
    armv7l|armv7)   echo armv7 ;;
    i386|i686)      echo 386 ;;
    *) fail "неподдерживаемая архитектура: $arch" ;;
  esac
}

latest_version() {
  fetch "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' | head -n1 | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/'
}

verify_checksum() { # file checksums.txt
  expected="$(grep " $(basename "$1")\$" "$2" | awk '{print $1}')"
  [ -n "$expected" ] || { warn "нет чексуммы для $(basename "$1"), пропускаю проверку"; return 0; }

  if has sha256sum; then
    actual="$(sha256sum "$1" | awk '{print $1}')"
  elif has shasum; then
    actual="$(shasum -a 256 "$1" | awk '{print $1}')"
  else
    warn "нет sha256sum/shasum, пропускаю проверку"; return 0
  fi

  [ "$expected" = "$actual" ] || fail "чексумма не совпала: ожидалось $expected, получено $actual"
  info "Чексумма OK"
}

pick_install_dir() {
  if [ -n "${INSTALL_DIR:-}" ]; then
    echo "$INSTALL_DIR"
  elif [ -w /usr/local/bin ] || has sudo; then
    echo /usr/local/bin
  else
    echo "$HOME/.local/bin"
  fi
}

main() {
  os="$(detect_os)"
  arch="$(detect_arch)"

  version="${VERSION:-}"
  if [ -z "$version" ]; then
    info "Ищу последнюю версию..."
    version="$(latest_version)"
    [ -n "$version" ] || fail "не удалось определить последнюю версию (лимит GitHub API?). Укажи VERSION=vX.Y.Z"
  fi
  case "$version" in v*) ;; *) version="v$version" ;; esac
  plain_version="${version#v}"

  archive="${BINARY}_${plain_version}_${os}_${arch}.tar.gz"
  base_url="https://github.com/${REPO}/releases/download/${version}"

  tmp="$(mktemp -d)"
  trap 'rm -rf "$tmp"' EXIT INT TERM

  info "Скачиваю ${BINARY} ${version} (${os}/${arch})"
  download "${base_url}/${archive}" "${tmp}/${archive}" \
    || fail "не удалось скачать ${base_url}/${archive}"

  if download "${base_url}/checksums.txt" "${tmp}/checksums.txt" 2>/dev/null; then
    verify_checksum "${tmp}/${archive}" "${tmp}/checksums.txt"
  else
    warn "checksums.txt не найден в релизе, пропускаю проверку"
  fi

  tar -xzf "${tmp}/${archive}" -C "$tmp"
  [ -f "${tmp}/${BINARY}" ] || fail "в архиве нет файла ${BINARY}"
  chmod +x "${tmp}/${BINARY}"

  dir="$(pick_install_dir)"
  mkdir -p "$dir" 2>/dev/null || true

  if [ -w "$dir" ]; then
    mv "${tmp}/${BINARY}" "${dir}/${BINARY}"
  else
    info "Нужны права для записи в ${dir}, запрашиваю sudo"
    sudo mkdir -p "$dir"
    sudo mv "${tmp}/${BINARY}" "${dir}/${BINARY}"
  fi

  info "Установлено: ${dir}/${BINARY}"

  case ":$PATH:" in
    *":$dir:"*) ;;
    *) warn "${dir} нет в PATH. Добавь в ~/.bashrc или ~/.zshrc:"
       printf '    export PATH="%s:$PATH"\n' "$dir" ;;
  esac

  "${dir}/${BINARY}" --version 2>/dev/null || true
}

main "$@"

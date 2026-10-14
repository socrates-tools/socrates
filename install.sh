#!/bin/sh
set -eu

repository=socrates-tools/socrates
releases="${SOCRATES_RELEASES_URL:-https://github.com/$repository/releases}"
install_dir="${SOCRATES_INSTALL_DIR:-$HOME/.local/bin}"
supported_platforms="https://github.com/$repository#supported-platforms"

refuse() {
  printf 'socrates install: %s\n' "$1" >&2
  exit 1
}

say() {
  printf '%s\n' "$1"
}

release_url() {
  if [ -n "${SOCRATES_VERSION:-}" ]; then
    printf '%s/download/%s' "${releases%/}" "$SOCRATES_VERSION"
  else
    printf '%s/latest/download' "${releases%/}"
  fi
}

fetch() {
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL -o "$2" "$1"
  elif command -v wget >/dev/null 2>&1; then
    wget -q -O "$2" "$1"
  else
    refuse "neither curl nor wget is on PATH, and one is needed to download socrates: install either and run this again"
  fi
}

download() {
  fetch "$1/$2" "$work/$2" || refuse "could not download $1/$2: check the release exists, or name one with SOCRATES_VERSION=v<version>"
}

sha256_of() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | cut -d ' ' -f 1
  elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$1" | cut -d ' ' -f 1
  else
    refuse "neither sha256sum nor shasum is on PATH, and one is needed to check the download: install either and run this again"
  fi
}

operating_system() {
  uname -s | tr '[:upper:]' '[:lower:]'
}

architecture() {
  machine=$(uname -m)
  case "$machine" in
  x86_64 | amd64) printf 'amd64' ;;
  aarch64 | arm64) printf 'arm64' ;;
  *) printf '%s' "$machine" ;;
  esac
}

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
trap 'exit 130' INT TERM

from=$(release_url)
platform="$(operating_system)_$(architecture)"

download "$from" checksums.txt

if command -v cosign >/dev/null 2>&1; then
  download "$from" checksums.txt.bundle
  cosign verify-blob "$work/checksums.txt" \
    --bundle "$work/checksums.txt.bundle" \
    --certificate-identity-regexp "https://github.com/$repository/.*" \
    --certificate-oidc-issuer https://token.actions.githubusercontent.com >/dev/null 2>&1 ||
    refuse "the signature over checksums.txt at $from does not verify, so nothing was installed: report it at https://github.com/$repository/security"
fi

archive=$(sed -n "s/^[0-9a-f]\{64\}  \(socrates_[^_]*_${platform}\.tar\.gz\)\$/\1/p" "$work/checksums.txt")
if [ -z "$archive" ]; then
  carried=$(sed -n 's/^[0-9a-f]\{64\}  socrates_[^_]*_\([^_]*\)_\([^.]*\)\.tar\.gz$/\1\/\2/p' "$work/checksums.txt" | tr '\n' ' ')
  refuse "this release carries no archive for $(printf '%s' "$platform" | tr '_' '/'), only for ${carried% }, so nothing was installed: see $supported_platforms"
fi

download "$from" "$archive"

expected=$(grep "  $archive\$" "$work/checksums.txt" | cut -d ' ' -f 1)
actual=$(sha256_of "$work/$archive")
[ "$expected" = "$actual" ] ||
  refuse "the checksum of $archive does not match checksums.txt, so nothing was installed: run this again, and if it still differs, report it at https://github.com/$repository/security"

tar -xzf "$work/$archive" -C "$work" socrates ||
  refuse "$archive holds no socrates binary, so nothing was installed: report it at https://github.com/$repository/issues"

mkdir -p "$install_dir" ||
  refuse "could not create $install_dir: name a directory you can write with SOCRATES_INSTALL_DIR"
chmod 755 "$work/socrates"
mv -f "$work/socrates" "$install_dir/socrates" ||
  refuse "could not write $install_dir/socrates: name a directory you can write with SOCRATES_INSTALL_DIR"

"$install_dir/socrates" skill ||
  refuse "socrates is at $install_dir/socrates, and installing its skills failed: run $install_dir/socrates skill to see why"

say "installed $install_dir/socrates: $("$install_dir/socrates" version)"

case ":$PATH:" in
*":$install_dir:"*) ;;
*)
  say "$install_dir is not on your PATH, so your shell will not find socrates. Add this to your shell's startup file:"
  say "  export PATH=\"$install_dir:\$PATH\""
  ;;
esac

# 2: The release carries an installer that picks the right archive

## What to build

`install.sh` at the repository root, attached to every release and listed in its checksums, as the spec describes.

## Acceptance criteria

- [ ] On `x86_64` Linux, `aarch64` Linux and `arm64` Darwin, the installer places the binary from that pair's archive and nothing from another, and every skill is installed.
- [ ] A pair the release carries no archive for is refused, naming the pair and the pairs it carries, and nothing is installed.
- [ ] An archive whose checksum does not match is refused and nothing is installed.
- [ ] A missing install directory is created, and one not on `PATH` is reported with how to add it.
- [ ] `SOCRATES_INSTALL_DIR`, `SOCRATES_VERSION` and `SOCRATES_RELEASE_URL` do what the spec says.
- [ ] Running it a second time replaces the binary and the skills.
- [ ] `.goreleaser.yaml` attaches `install.sh` and lists it in `checksums.txt`, and the suite fails where either stops holding.
- [ ] The installer names no supported platform.

## User stories

- 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 17

## Blocked by

- 1

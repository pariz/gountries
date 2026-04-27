# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] — Unreleased

This is a **breaking** release. The data pipeline has been rebuilt from scratch
and the embedded country/subdivision data is now sourced from upstream
projects that are actively maintained, instead of the previous custom toolchain
that had grown stale.

### Pinning the old release

If you need the v0.1.x behavior and data exactly, pin:

```
go get github.com/pariz/gountries@v0.1.6
```

### Added

- `Country.AltSpellings []string` — short forms and unofficial names (e.g.
  `"United States"`, `"USA"`, `"États-Unis"`). These are indexed by
  `FindCountryByName`, so common short forms now resolve.
- `cmd/sync-data` — Go program that regenerates `data/yaml/` from vendored
  upstream submodules (`third_party/countries`, `third_party/mledoze-countries`).
  Run with `go run ./cmd/sync-data`.
- Weekly `sync-upstream` GitHub Action that bumps both upstream submodules,
  regenerates data, runs tests, and opens a PR if anything changed.
- `TestDataInvariants` — sanity tests that guard against silent upstream
  regressions (mass deletion, schema breakage, canary-lookup failures).

### Changed

- **Data sources.** Country and subdivision data is now derived from
  [countries/countries](https://github.com/countries/countries) (MIT, primary)
  and [mledoze/countries](https://github.com/mledoze/countries) (ODbL,
  supplemental for borders / TLDs / translations).
- **Data shape.** Many countries now have richer translations (Sweden goes
  from 7 to 23 languages), more accurate coordinates, more complete
  subdivisions (5,059 entries across 200 countries vs. the previous ~195
  files with sparser coverage).
- **Embedding mechanism.** Switched from `go-bindata` to the standard library
  `//go:embed` (Go 1.16+). The data files are still loaded from
  `data/yaml/{countries,subdivisions}/` but you no longer need an external
  tool to regenerate the embed.
- **`NewFromPath(path string)`** — still takes a directory path, but reads
  through `os.DirFS(path)`. The directory must contain `countries/` and
  `subdivisions/` subdirectories (same layout as before).
- **Go version** bumped to 1.21 in `go.mod` (was 1.17).
- **CI workflow** modernized: `actions/setup-go@v5`, `go-version-file: go.mod`,
  `go test -race`, `go vet`.

### Removed

- **`bindata.go` and all its exports** — `Asset`, `AssetDir`, `AssetNames`,
  `AssetInfo`, `MustAsset`, `RestoreAsset`, `RestoreAssets`, `AssetNameMap`.
  If you were importing any of these (none are part of the documented public
  API but they were exported), you'll need to read files via your own means.
- The old `gountries-creator` workflow is no longer the data source. The repo
  is preserved separately for history but is not used by v1.

### Fixed

- `FindCountryByName` now resolves common short forms and alt-spellings via
  the new `AltSpellings` index — e.g. `"United States"` correctly returns the
  US even though its canonical `iso_short_name` is `"United States of America"`.
- Czech Republic is now spelled "Czechia" (matching ISO since 2016).
- Bulgaria's currency reflects the EUR transition effective 1 Jan 2026.

### License notes

The Go code remains MIT-licensed. The regenerated YAML data is a derivative
of two upstreams:
- countries/countries — MIT
- mledoze/countries — ODbL

The `data/yaml/` directory should be considered a Derivative Database under
ODbL by virtue of incorporating mledoze fields. See
[third_party/mledoze-countries/LICENSE](third_party/mledoze-countries/LICENSE).
Use of the gountries Go API to query the data does not create a Derivative
Database for your application — only modifying and redistributing the data
itself triggers the ODbL share-alike clause.

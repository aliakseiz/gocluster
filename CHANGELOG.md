# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Fixed

### Deprecated

### Removed

### Security

## [1.1.0] - 2022-04-04

### Added
- CHANGELOG.md file
- Tests for `Point.Included` field

### Changed
- `gci`, `gofumpt` and `godot` linters applied
- Module name to `github.com/aliakseiz/gocluster`
- Moved tests to `*_test` packages
- Synchronized the implementation with parent [mapbox/supercluster](https://github.com/mapbox/supercluster/blob/main/index.js)
- Points with no coordinates are skipped

### Fixed
- Failing tests
- Incorrect IDs in `Point.Included` field for some clusters 
- Numerous typos in comments
- Clustering in some edge cases, i.e. view covers western and eastern hemispheres

### Removed
- `Descendants` field from the `Point`
- Nested structures in test helpers

## [1.0.2] - 2021.09.16
Fork from [AlekseevAV/cluster](https://github.com/AlekseevAV/cluster)
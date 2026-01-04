# Changelog

## [0.7.1](https://github.com/hairyhenderson/go-codeowners/compare/v0.7.0...v0.7.1) (2026-01-04)


### Dependencies

* **actions:** bump actions/checkout from 4 to 6 ([#66](https://github.com/hairyhenderson/go-codeowners/issues/66)) ([2cdb2d6](https://github.com/hairyhenderson/go-codeowners/commit/2cdb2d64f1009b4aed72932d36daaf5990689ce1))
* **actions:** bump actions/create-github-app-token from 1 to 2 ([#58](https://github.com/hairyhenderson/go-codeowners/issues/58)) ([e6ba277](https://github.com/hairyhenderson/go-codeowners/commit/e6ba277d866673f8287962bf0b7e28e405f0cf21))
* **actions:** bump actions/setup-go from 4 to 5 ([#55](https://github.com/hairyhenderson/go-codeowners/issues/55)) ([0c001dc](https://github.com/hairyhenderson/go-codeowners/commit/0c001dc0be724863d04a110d9be3a229898de17a))
* **actions:** bump golangci/golangci-lint-action from 6 to 7 ([#57](https://github.com/hairyhenderson/go-codeowners/issues/57)) ([c7a0b1a](https://github.com/hairyhenderson/go-codeowners/commit/c7a0b1ab841f1b0f278e745879b3b5d71446f7c5))
* **actions:** bump golangci/golangci-lint-action from 7 to 8 ([#59](https://github.com/hairyhenderson/go-codeowners/issues/59)) ([8df2846](https://github.com/hairyhenderson/go-codeowners/commit/8df28460bb5b51f2f9eb3b5ba807443a19e090cb))

## [0.7.0](https://github.com/hairyhenderson/go-codeowners/compare/v0.6.1...v0.7.0) (2024-12-07)


### Features

* add PatternIndex method ([#48](https://github.com/hairyhenderson/go-codeowners/issues/48)) ([a5c2e59](https://github.com/hairyhenderson/go-codeowners/commit/a5c2e593dcbe38fa52c24c95a72d1ca3f5e47ee8))


### Dependencies

* **go:** Bump github.com/stretchr/testify from 1.9.0 to 1.10.0 ([#49](https://github.com/hairyhenderson/go-codeowners/issues/49)) ([4f55176](https://github.com/hairyhenderson/go-codeowners/commit/4f551769326fc31237d06b29a4a0497c78bd882f))

## [0.6.1](https://github.com/hairyhenderson/go-codeowners/compare/v0.6.0...v0.6.1) (2024-10-25)


### Bug Fixes

* **parsing:** Check for CODEOWNERS files for both GitHub and GitLab ([#45](https://github.com/hairyhenderson/go-codeowners/issues/45)) ([db9c01e](https://github.com/hairyhenderson/go-codeowners/commit/db9c01eb61a0e5975521721a102e8c54b6dc2876)), closes [#44](https://github.com/hairyhenderson/go-codeowners/issues/44)
* **parsing:** Ignore block and inline comments ([#46](https://github.com/hairyhenderson/go-codeowners/issues/46)) ([968c9ea](https://github.com/hairyhenderson/go-codeowners/commit/968c9eaf0924c1912731d2729f26d2c691b8d4b1))

## [0.6.0](https://github.com/hairyhenderson/go-codeowners/compare/v0.5.0...v0.6.0) (2024-09-27)


### Features

* **errors:** Allow checking codeowners file not found ([#41](https://github.com/hairyhenderson/go-codeowners/issues/41)) ([d659b73](https://github.com/hairyhenderson/go-codeowners/commit/d659b73a08c1a1111c5e9c4c1136472c8ca28a4b))


### Bug Fixes

* **perf:** Reduce allocations ([#39](https://github.com/hairyhenderson/go-codeowners/issues/39)) ([2ca66e0](https://github.com/hairyhenderson/go-codeowners/commit/2ca66e0194d2b9077c63a9d1eace8f1083675fa9))

# Changelog

## [0.7.1](https://github.com/hairyhenderson/go-codeowners/compare/v0.7.0...v0.7.1) (2026-05-16)


### Bug Fixes

* **lint:** add minimal permissions to workflow jobs ([#72](https://github.com/hairyhenderson/go-codeowners/issues/72)) ([d2a2a81](https://github.com/hairyhenderson/go-codeowners/commit/d2a2a81d693952ec6139999ab86ab558029f22f9))
* **lint:** extract repeated string literals into constants in tests ([#71](https://github.com/hairyhenderson/go-codeowners/issues/71)) ([1aaddf1](https://github.com/hairyhenderson/go-codeowners/commit/1aaddf1d8cd06ece36384c8eaf13fba1ba46c1cb))


### Dependencies

* **actions:** bump actions/checkout from 4 to 6 ([#66](https://github.com/hairyhenderson/go-codeowners/issues/66)) ([2cdb2d6](https://github.com/hairyhenderson/go-codeowners/commit/2cdb2d64f1009b4aed72932d36daaf5990689ce1))
* **actions:** bump actions/create-github-app-token from 1 to 2 ([#58](https://github.com/hairyhenderson/go-codeowners/issues/58)) ([e6ba277](https://github.com/hairyhenderson/go-codeowners/commit/e6ba277d866673f8287962bf0b7e28e405f0cf21))
* **actions:** bump actions/create-github-app-token from 2 to 3 ([#68](https://github.com/hairyhenderson/go-codeowners/issues/68)) ([2bc92fa](https://github.com/hairyhenderson/go-codeowners/commit/2bc92fa61fdb27b3deb77c8a893127988df4aa80))
* **actions:** bump actions/setup-go from 4 to 5 ([#55](https://github.com/hairyhenderson/go-codeowners/issues/55)) ([0c001dc](https://github.com/hairyhenderson/go-codeowners/commit/0c001dc0be724863d04a110d9be3a229898de17a))
* **actions:** bump actions/setup-go from 5 to 6 ([#64](https://github.com/hairyhenderson/go-codeowners/issues/64)) ([e1d5ea5](https://github.com/hairyhenderson/go-codeowners/commit/e1d5ea58953962b4ebd2daf53be77a39050cf114))
* **actions:** bump actions/stale from 9 to 10 ([#63](https://github.com/hairyhenderson/go-codeowners/issues/63)) ([5bd9e73](https://github.com/hairyhenderson/go-codeowners/commit/5bd9e7375c6865c632ab56653e3037c4bbc2e6b2))
* **actions:** bump golangci/golangci-lint-action from 6 to 7 ([#57](https://github.com/hairyhenderson/go-codeowners/issues/57)) ([c7a0b1a](https://github.com/hairyhenderson/go-codeowners/commit/c7a0b1ab841f1b0f278e745879b3b5d71446f7c5))
* **actions:** bump golangci/golangci-lint-action from 7 to 8 ([#59](https://github.com/hairyhenderson/go-codeowners/issues/59)) ([8df2846](https://github.com/hairyhenderson/go-codeowners/commit/8df28460bb5b51f2f9eb3b5ba807443a19e090cb))
* **actions:** bump golangci/golangci-lint-action from 8 to 9 ([#65](https://github.com/hairyhenderson/go-codeowners/issues/65)) ([2663bcc](https://github.com/hairyhenderson/go-codeowners/commit/2663bcc3fc6388ba6d33f144d5f3f84e584b97fd))
* **actions:** bump googleapis/release-please-action from 4 to 5 ([#69](https://github.com/hairyhenderson/go-codeowners/issues/69)) ([d16278b](https://github.com/hairyhenderson/go-codeowners/commit/d16278b6ffb352a856b54d8eeff424d7c53d1804))
* **actions:** bump webiny/action-conventional-commits ([#70](https://github.com/hairyhenderson/go-codeowners/issues/70)) ([9f6e975](https://github.com/hairyhenderson/go-codeowners/commit/9f6e9755a73d75257b1ff973db24d5e4ff644d88))
* **go:** bump github.com/stretchr/testify from 1.10.0 to 1.11.1 ([#62](https://github.com/hairyhenderson/go-codeowners/issues/62)) ([cf6d348](https://github.com/hairyhenderson/go-codeowners/commit/cf6d348cd6dce37904fb902866f43dfe07883e51))

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

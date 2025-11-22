## [1.8.5](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.8.4...v1.8.5) (2025-11-22)


### Bug Fixes

* **parser:** skip news without schedule data ([5668b61](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/5668b61a754950dfeee6a49db6cdcbaf8d753524))

## [1.8.4](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.8.3...v1.8.4) (2025-11-17)


### Bug Fixes

* **parser:** normalize schedule paragraphs ([def3d4a](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/def3d4acad22ade3ea29688ef96555d43be50f4a))

## [1.8.3](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.8.2...v1.8.3) (2025-11-17)


### Bug Fixes

* **api:** raise default rate limit to 6 ([5e1ab42](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/5e1ab42b0c22cb24de0e79786d220d3f085cc13c))

## [1.8.2](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.8.1...v1.8.2) (2025-11-16)


### Bug Fixes

* **utils:** harden cyrillic date parsing ([acdca95](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/acdca952d1c579b2d6569dc2289a42f2f4081685))

## [1.8.1](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.8.0...v1.8.1) (2025-11-16)


### Bug Fixes

* **utils:** extend date extractor patterns ([8989d3e](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/8989d3eb50f68070a16e5599923999f4b2c0880e))
* **utils:** remove overly broad date regex ([64c64f2](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/64c64f25d1069eab71cf2076eadbd8ffe08f08d5))

# [1.8.0](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.7.0...v1.8.0) (2025-11-14)


### Bug Fixes

* **auth:** skip auth for api keys list ([09bbea9](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/09bbea9d19ceddc5d0e454dedfbda51fa98edde1))
* **middleware:** skip rate limiter for api keys list ([ba72e2e](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/ba72e2e6e5e5b091e82e30bc0e4301f80e8718fb))


### Features

* **handlers:** add json body to create api keys ([45b043c](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/45b043c81522e105b9bcb3e8f5395125ad074436))
* **handlers:** split api key update and delete handlers ([faee55b](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/faee55b64ceadf76093a03055c3e25afe84453a5))
* **server:** standardize api key CRUD routes ([816ce67](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/816ce67be43d7ee095fcfbd8f9282aeb311beec9))

# [1.7.0](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.6.0...v1.7.0) (2025-11-14)


### Features

* **api:** add GET schedule endpoint ([280e57a](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/280e57af56f8e3802a34fbecd53def67a7f1fe75))

# [1.6.0](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.5.0...v1.6.0) (2025-11-14)


### Features

* **handlers:** support schedule queue filtering ([761e3bd](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/761e3bde268f2af43c8c3089635c1b74b1d7745f))
* **models:** add schedule date field ([4973f14](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/4973f14e6af92a07ae816f66ff6d03e4491c9ecd))
* **utils:** add schedule date extractor ([8473f32](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/8473f3254f5aacee7a09f9f049137d0f0f184f56))

# [1.5.0](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.4.0...v1.5.0) (2025-10-30)


### Bug Fixes

* **parser:** add logging to schedule parsing logic ([b0df7b5](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/b0df7b5de11277aef1441eeeeb78b6d5696bb238))


### Features

* **parser:** enhance schedule parsing with table format support ([5ccd4ab](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/5ccd4abe7a11497ed91ffa41501beae707f20b07))

# [1.4.0](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.3.0...v1.4.0) (2025-10-23)


### Bug Fixes

* **db:** use template1 for database existence check ([1eae26f](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/1eae26fc943461250408125ab77e85b7a4251687))


### Features

* **parser:** enhance schedule detection with regex patterns ([7c94a13](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/7c94a136e5dc61052b28d57869b1668bc4a1cf5f))
* **parser:** require NEWS_URL for news parsing initialization ([02e2705](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/02e2705443ea92364f196c127aacd822c4fd52b8))
* **parser:** update news fetching logic with environment variables ([ddfc23a](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/ddfc23a6e050ea1c8b2e09b99a142b230b97400f))
* **server:** load environment variables from .env file ([cea3140](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/cea31404a7b4232d1f532489d61f37f47097d3c5))

# [1.3.0](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.2.3...v1.3.0) (2025-06-15)


### Bug Fixes

* **db:** improve database initialization and migration error handling ([83663b2](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/83663b2f4d72ad0933ea36be3fa42e78e06dad2f))


### Features

* **db:** ensure database exists before connecting ([cef4769](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/cef47692c4975b924adde559f5c6a875c9f82dea))

## [1.2.3](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.2.2...v1.2.3) (2025-02-24)


### Bug Fixes

* correct reference to APP_VERSION in server logging ([9d05a37](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/9d05a37e51d0679e4bc04ca26b90f56992982573))

## [1.2.2](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.2.1...v1.2.2) (2025-02-24)


### Bug Fixes

* improve server startup logging and error handling ([9e4fd2c](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/9e4fd2c8be14397f6a82288b89c63a0bdb0bf20f))
* update APP_VERSION variable path in Dockerfile to correct import path ([0498ef5](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/0498ef5430938dc5ea2f94ebaeda5eee25b44ddb))
* update version variable to exported format for consistency ([f41cc32](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/f41cc3238aaf1f20ba81e82ae8f20c04ce5b691c))

## [1.2.1](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.2.0...v1.2.1) (2025-02-24)


### Bug Fixes

* change version constant to variable for dynamic versioning ([5ce481b](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/5ce481b17b495b5f207c8ca1a786c685c3fe7df6))
* remove unused dispatch release and deploy steps from workflow ([5d0777d](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/5d0777d7f89691b733ea8808b899cbea60df7f30))
* update changelog configuration and clean up CHANGELOG ([24f762c](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/24f762c2e16cf0be936b1c44b7c2a0b286c19305))
* update changelog configuration and simplify release message format ([ce05065](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/ce05065394366bb6ec513729796658e5c8a959e0))

# [1.2.0](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.1.0...v1.2.0) (2025-02-15)


### Bug Fixes

* reorder import statements for consistency ([488482c](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/488482c2afbc5f0c51cc0af16c6c59c5b6ed8df4))
* update license badge and remove outdated sections from README ([3c80ba5](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/3c80ba5008cb81dbd211085a4ff547b46f379d54))


### Features

* add API documentation for Cherkasyoblenergo API ([e0dbbf0](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/e0dbbf0d0b16a8642a819cc1e39f54e278758ada))
* add restart policy to app and db services ([8e2ab22](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/8e2ab221b768391346c72df8fadfd0525fe13e35))
* bump app version to 1.1.1 ([102e3ef](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/102e3ef0e4cd87b4e6ee3f8c0f33d6585bda4ecd))
* change news fetch cron interval to 10 minutes ([1e2f391](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/1e2f391bbf454e25913808fdebcae8b316f97c15))
* update README with new structure ([876976d](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/876976da33cb14940ffa0a72459d4522d7f35a76))

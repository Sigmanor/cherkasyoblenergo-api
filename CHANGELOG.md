# [2.0.0](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.11.1...v2.0.0) (2025-12-28)


### Features

* **cache:** add schedule cache implementation ([23f57b9](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/23f57b9))
* **config:** add runtime defaults ([b330c4d](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/b330c4d))
* **db:** configure gorm logger ([d560efa](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/d560efa))
* **handlers:** add schedule response caching ([a1c54c4](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/a1c54c4))
* **logger:** add global logger setup ([d57ecd6](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/d57ecd6))
* **logger:** add gorm adapter ([cac23a7](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/cac23a7))
* **metrics:** add fiber metrics middleware ([a4291ee](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/a4291ee))
* **middleware:** add https enforcement middleware ([3275ecc](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/3275ecc))
* **middleware:** add ip rate limiter middleware ([b8a5a2c](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/b8a5a2c))
* **middleware:** add optional api key auth ([5918f1d](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/5918f1d))
* **models:** add ip rate limit model ([53b8193](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/53b8193))
* **server:** improve startup and middleware ([46d763f](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/46d763f))


### Bug Fixes

* **parser:** guard parsing with context ([c6867f8](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/c6867f8))


### BREAKING CHANGES

* Major refactoring of middleware stack, logger system, and caching layer


## [1.11.1](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.11.0...v1.11.1) (2025-12-24)


### Bug Fixes

* **handlers:** order schedule query by id ([1e6dfdf](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/1e6dfdf99e5ec76c965da4167e277a6babf1be9f))

# [1.11.0](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.10.0...v1.11.0) (2025-12-15)


### Features

* **api:** add webhook endpoints ([ea55d36](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/ea55d362bb465f00dab7277bfe74dde845b59494))
* **handlers:** add webhook management handler ([c834892](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/c8348920b567b7f220c4ee5b2692163723a9645a))
* **middleware:** expose API key helper ([7daedb7](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/7daedb7fc4e34f3d21c7fb4005f00374f162bf6a))
* **models:** add webhook fields to api key ([43189da](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/43189da753fc081150666de989df635e7945f9e9))
* **parser:** trigger webhook after saving schedule ([dff9d09](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/dff9d090205a9976c81e2c03ce072510a83adb37))
* **webhook:** add webhook delivery helpers ([8b5efef](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/8b5efeff0fd758a114502766256565b4f3968a23))

# [1.10.0](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.9.0...v1.10.0) (2025-12-07)


### Bug Fixes

* **auth:** remove unused api key skips ([8091cca](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/8091ccad5de3d13c82ebcf2cbf81741e532810fb))
* **db:** respect configured database name ([7e546f3](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/7e546f3589dc1a28eef9465942c0c9a6fc2e5ed6))
* **middleware:** refresh cached rate limiters ([2033ff5](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/2033ff5f78f7c31f92020635c0e17140364fcd8d))
* **middleware:** use structured request logging ([cdc63dc](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/cdc63dc1e38f6bd4d2e052e7892d3847b28e8122))
* **models:** map schedule_date column ([dda0d44](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/dda0d44af20074915a6078117c1385efc6f84416))
* **parser:** expose cron and fill schedule date ([7106f9e](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/7106f9e3c2b4fcf03e902236fb6336b0e593c8fd))
* **schedule:** filter schedules by stored dates ([0f26aeb](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/0f26aebeb98ec572355c64b93faaa5a9ad798139))
* **schedule:** resolve today and tomorrow dates ([7bf1894](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/7bf18941755f91d933c0476d6348850927ef07e0))


### Features

* **server:** add graceful shutdown handling ([8fd207e](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/8fd207e2881a5ede8a84ce5d4223998150edf4b3))

# [1.9.0](https://github.com/Sigmanor/cherkasyoblenergo-api/compare/v1.8.5...v1.9.0) (2025-12-06)


### Features

* Add `by_schedule_date` filtering for schedules and enhance date extraction to include year. ([4e777af](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/4e777af709ccdd4576f814f50a645be01bad1e9f))
* Add `by_schedule_date` option to schedule handler and update date extraction to YYYY-MM-DD format. ([b132812](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/b132812cc2cfd7e4f3f5dda92922a1edb58569d9))
* enhance date extraction to include year calculation and format as YYYY-MM-DD ([61e8433](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/61e8433894ec8af67ea20782a7aca1daefe546e1))
* Implement schedule filtering by extracted schedule date, now including an inferred year in the extraction logic. ([a23b90a](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/a23b90a2bc7f63207ec034dc46979013ae0b28dc))
* Introduce `by_schedule_date` API option, update `schedule_date` format to YYYY-MM-DD, and enhance date extraction with year inference. ([8426f33](https://github.com/Sigmanor/cherkasyoblenergo-api/commit/8426f3319e761ba94898fa632800dbbf4e26561f))

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

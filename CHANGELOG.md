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

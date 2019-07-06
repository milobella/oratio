# Changelog
All notable changes to the helm chart of this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.3]
### Added
- Clock ability configuration

### Modified
- Closed the oratio service to be accessible only from inside the cluster (i.e. by the kong service).

## [0.1.2]
### Added
- ``TRIGGER_SHOPPING_LIST`` intent to shoppinglist-ability configuration.

### Changed
- Use ``Always`` imagePullPolicy.
- Set deployment strategy to recreate.

## [0.1.1]
### Added
- Set ``ORATIO_ANIMA_PORT`` and ``ORATIO_CEREBRO_PORT`` environment variables.

### Changed
- Add a config map which comes to replace the environment variables setting.

## [0.1.0]
### Added
- Created the helm chart with a deployment and a service.

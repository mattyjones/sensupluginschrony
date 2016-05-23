#Change Log
This project adheres to [Semantic Versioning](http://semver.org/).

This CHANGELOG follows the format listed at [Keep A Changelog](http://keepachangelog.com/)

## [Unreleased]
- update pkg names
- bump to golang 1.6.2
- update vendoring

### Changed
- removed output from an `ok` check to conform to std unix practice of exiting cleaning upon 0 status

## [0.0.1] - 2016-03-16
### Added
- initial external stable release

## [0.1.0] - 2016-03-26
### Added
- ability to check for the Last Offset and RMS Offset
- additional tests and documentation

### Changed
- split `overThreshold` into `overIntThreshold` and `overFloatThreshold`

[Unreleased]: https://github.com/yieldbot/sensupluginschrony/compare/0.1.0....HEAD
[0.1.0]: https://github.com/yieldbot/sensupluginschrony/compare/0.0.1....0.1.0

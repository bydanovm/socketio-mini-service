# Changelog

## [Unreleased]

### Added
- Token validation function support via `TokenValidate` method
- Enhanced event handling with payload marshaling and client context

### Changed
- `AddEvent` method signature to accept event name and handler function directly
- Event handler implementation to automatically marshal payload data

### Improved
- Event processing with better error handling and client context management
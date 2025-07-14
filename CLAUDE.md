# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

### Building and Testing
- `make build` - Build the fieldday binary for current platform
- `make build-raspi` - Cross-compile for Raspberry Pi (ARM7)
- `go run main.go <database_file>` - Run the application locally (requires database file argument)
- `./fieldday test.db` - Test built binary with test database
- `go test . ./visitorstore ./morse` - Run all tests (excludes legacy exportbolt)
- `go test ./visitorstore -v` - Run visitorstore unit tests with verbose output
- `go test ./morse -v` - Run morse package tests with verbose output
- `go test . -v` - Run HTTP handler integration tests with verbose output

### CSS Development
- `npx tailwindcss -i input.css -o static/css/tailwind.css --watch` - Watch and rebuild Tailwind CSS
- Tailwind config watches `templates/*.go.html` files

### Data Export and Import
- `go run cmd/exportbolt/main.go --db <db_file> --out <csv_file>` - Export BoltDB data to CSV (legacy)
- `go run cmd/importcsv/main.go <csv_file> <db_file>` - Import CSV data into SQLite database
- `sqlite3 <db_file> ".mode csv" ".output file.csv" "SELECT * FROM visitors;"` - Export SQLite to CSV

### Deployment (Production)
- `make install` - Install to /opt/fieldday with systemd service
- `make user` - Create nfarl user with kiosk configuration
- `make start` / `make stop` - Control systemd service

## Project Architecture

### Core Components
- **main.go**: Web server with embedded templates and static files using Go's embed package
- **visitorstore/**: SQLite-based visitor data storage using modernc.org/sqlite
- **morse/**: Morse code generation and WAV audio file creation
- **templates/**: Go HTML templates for web interface
- **static/**: CSS (Tailwind), JavaScript, and image assets
- **cmd/importcsv/**: Utility for importing CSV data into SQLite
- **cmd/exportbolt/**: Legacy utility for exporting BoltDB data

### Data Flow
1. Visitors fill out web form at `/new` endpoint
2. Form data is decoded using gorilla/schema into Visitor struct  
3. Visitor struct is saved to SQLite via SQL queries
4. Confirmation page plays Morse code of callsign (or "73") via generated WAV file
5. All visitors can be viewed at `/list` endpoint

### Database
- Uses SQLite (embedded relational database) with modernc.org/sqlite (pure Go implementation)
- Visitor struct has schema tags for form decoding
- Database schema automatically initialized on startup
- Database file path is provided as command line argument
- Standard SQL queries for data operations

### Web Interface
- Simple HTTP server using standard library (no external web framework)
- Templates use Go's html/template package with embedded file system
- Tailwind CSS for styling with mobile-responsive design
- Static files (CSS, JS, images) served via embedded file system

### Morse Code Feature
- Generates WAV audio files server-side for callsign playback
- Uses `/morse-audio?callsign=<callsign>` endpoint to serve audio
- Default 15 WPM, 600 Hz tone with PARIS timing standard

## Key Design Patterns
- Embedded templates and static files for single binary deployment
- Simple struct-based HTTP handlers with error handling
- Form data validation in visitorstore package
- Separation of concerns: main.go (HTTP), visitorstore (data), morse (audio)
- Pure Go SQLite implementation eliminates CGO dependencies for easier cross-compilation

## Development Notes
- Application designed for offline operation on single-board computers
- Built for Amateur Radio Field Day events with kiosk deployment
- Cross-platform build support via Makefile and GoReleaser
- Systemd service integration for unattended operation

## Testing
- Comprehensive test suite covering all major packages
- Unit tests for visitor database operations and CSV import
- Integration tests for HTTP handlers and form submission
- Morse code generation and audio functionality tests
- Test databases are automatically created and cleaned up
- Run tests before committing changes to ensure code quality
# Project Summary: Field Day Registration Kiosk

This document summarizes the "Field Day Registration Kiosk" project, a Go-based web application designed for Amateur Radio Field Day events.

## 1. Project Overview

The project provides a simple, robust, and offline-capable visitor registration system. It's primarily intended for deployment on single-board computers like Raspberry Pi, acting as a local kiosk. Key features include:

-   **Visitor Registration:** Simple web interface for collecting visitor information (name, callsign, email, etc.).
-   **Offline Operation:** Designed to function without an internet connection, crucial for remote field deployments.
-   **Data Persistence:** Uses BoltDB (via the `storm` library) for local data storage.
-   **Morse Code Playback:** Generates and plays Morse code audio for visitor callsigns directly in the browser.
-   **Embedded Assets:** Templates and static files are embedded directly into the Go binary for easy deployment.
-   **Modern UI:** Utilizes Tailwind CSS for a responsive and mobile-friendly web interface.
-   **Deployment Focus:** Includes `Makefile` targets and `goreleaser` configuration for streamlined building and deployment on Linux systems, including systemd service integration.

## 2. Technology Stack

-   **Language:** Go (Golang)
-   **Web Framework:** Standard `net/http` package for routing and serving.
-   **Templating:** Go's `html/template` package with embedded templates.
-   **Database:** BoltDB (key-value store) via the `github.com/asdine/storm/v3` library.
-   **Form Handling:** `github.com/gorilla/schema` for decoding form data.
-   **Frontend:** Tailwind CSS for styling, with embedded static assets (CSS, JS, images).
-   **Build/Deployment:** `Makefile` for common tasks (build, install, user setup), `goreleaser` for cross-platform binary releases.

## 3. Project Structure

The project follows a standard Go project layout with several key directories:

-   **`/` (Root):** Contains `main.go` (application entry point), `go.mod`/`go.sum` (Go modules), `README.md`, `LICENSE`, `Makefile`, and configuration files (`.goreleaser.yaml`, `tailwind.config.js`).
-   **`deploy/`:** Contains scripts and service files for deployment, including `systemd` service configuration and user setup scripts.
-   **`images/`:** Stores project-related images, such as screenshots.
-   **`morse/`:** Go package responsible for generating Morse code audio (WAV files) from text.
-   **`static/`:** Contains static web assets like CSS (Tailwind output), JavaScript (Bootstrap), and images. These are embedded into the binary.
-   **`templates/`:** Contains HTML templates (`.go.html` files) used by the Go `html/template` package. These are also embedded.
-   **`visitorstore/`:** Go package encapsulating the logic for interacting with the BoltDB database, including `Visitor` struct definition and methods for saving and listing visitors.

## 4. Core Functionality Breakdown

### 4.1. Web Server (`main.go`)

-   Initializes the `VisitorStore` with a specified database file.
-   Sets up HTTP routes for:
    -   `/`: Home page.
    -   `/new`: Visitor registration form (GET) and submission (POST).
    -   `/confirmation`: Displays confirmation after a successful registration.
    -   `/list`: Lists all registered visitors.
    -   `/morse-audio`: Serves generated Morse code WAV audio based on a `callsign` query parameter.
    -   `/privacy`: Privacy policy page.
-   Serves embedded static files and HTML templates.

### 4.2. Visitor Management (`visitorstore/`)

-   Defines the `Visitor` struct with fields for `ID`, `CreatedAt`, `FirstName`, `LastName`, `Callsign`, `Email`, and boolean flags for `Nfarl`, `Contactme`, `Youth`, `Firsttime`.
-   `NewVisitorStore`: Opens or creates a BoltDB database file.
-   `SaveVisitor`: Stores a new visitor record in the database, automatically setting `CreatedAt` if not provided.
-   `ListVisitors`: Retrieves all visitor records from the database.
-   `TotalVisitors`: Returns the count of registered visitors.

### 4.3. Morse Code Generation (`morse/`)

-   `morseCodeMap`: Maps characters to their corresponding Morse code patterns.
-   `calculateMorseTiming`: Determines dot, dash, and gap durations based on Words Per Minute (WPM) using the PARIS standard.
-   `newMorseAudio`: Pre-generates audio samples for dots, dashes, and various gaps (element, character, word) at a specified frequency.
-   `generateMorseAudio`: Assembles pre-generated samples to create the full audio waveform for a given text.
-   `writeWavHeader`: Helper function to write the WAV file header.
-   `GenerateWav`: The primary function to generate a complete WAV file (as a byte slice) for a given text, using default WPM and frequency.

## 5. Development History and Evolution

The project has undergone significant evolution, reflecting a journey of learning and optimization:

-   **Initial Versions (Django):** Started as a Python/Django project.
-   **Go Adoption (FD 2021):** Rewritten in Go, initially using the Gorilla web framework and SQLite with GORM.
-   **Simplification (FD 2022):** Switched to `net/http` for a simpler architecture, retaining `gorilla/schema`. Introduced GitHub Actions for CI/CD.
-   **Systemd & RTC (FD 2023):** Added `systemd` service integration and addressed real-time clock synchronization issues for offline deployments.
-   **BoltDB Migration (FD 2024):** Switched from SQLite to BoltDB/Storm to eliminate CGO dependencies, enabling faster compilation and easier cross-compilation. The `visitorstore` package was created for this abstraction.
-   **Tailwind CSS & Cloud Deployment (FD 2025):** Migrated from Bootstrap to Tailwind CSS for improved responsiveness. Added in-browser Morse code audio generation (serving WAV files) to support mobile users and cloud deployments.

This project serves as a practical example of building a self-contained Go web application for embedded systems, demonstrating robust offline capabilities and a clear evolution towards a more efficient and modern architecture.

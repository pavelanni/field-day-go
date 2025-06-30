# Possible Improvements for Field Day Registration Kiosk

This document outlines potential areas for improvement in the Field Day Registration Kiosk project, covering code quality, maintainability, robustness, and development workflow.

## 1. Code Quality and Maintainability

### 1.1. Refactor Template Rendering

**Current State:** The logic for parsing and executing HTML templates is duplicated across multiple HTTP handlers (`homeHandler`, `confirmHandler`, `listHandler`, `newVisitorHandler`, `privacyHandler`).

**Improvement:** Create a centralized helper function (e.g., `renderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{})`) that handles template parsing, execution, and error reporting. This would reduce code duplication and make it easier to manage template-related logic.

### 1.2. Consistent Error Handling

**Current State:** Error handling in `main.go` is inconsistent. Some errors (e.g., template parsing/execution) use `log.Fatal`, which crashes the server, while others use `http.Error`.

**Improvement:** All HTTP handler errors should use `http.Error` to return appropriate HTTP status codes and messages to the client, preventing server crashes and providing a better user experience. `log.Fatal` should be reserved for unrecoverable application startup errors.

### 1.3. Configuration Management

**Current State:** Configuration values like `port` and `thisYear` are hardcoded as global variables in `main.go`.

**Improvement:** Externalize configuration. This could involve:
-   Using environment variables.
-   Implementing a dedicated configuration struct that is passed to the `Server`.
-   Using a configuration file (e.g., JSON, YAML) for more complex settings.
This makes the application more flexible and easier to deploy in different environments.

### 1.4. Input Validation (VisitorStore)

**Current State:** The `SaveVisitor` method in `visitorstore` only checks for an empty `FirstName`.

**Improvement:** Implement more comprehensive input validation for the `Visitor` struct fields (e.g., validate `Callsign` format, `Email` format, ensure required fields are present). This can be done within the `SaveVisitor` method or by adding validation methods to the `Visitor` struct itself.

### 1.5. Database Closing (VisitorStore)

**Current State:** There is no explicit `Close` method for the `VisitorStore` to close the BoltDB connection.

**Improvement:** Add a `Close()` method to `VisitorStore` and ensure it's called when the application shuts down (e.g., using `defer server.store.Close()` in `main`). This ensures proper resource management and prevents potential data corruption.

## 2. Performance and Efficiency

### 2.1. Efficient Visitor Counting

**Current State:** The `TotalVisitors` function in `visitorstore` fetches all visitors to count them, which can be inefficient for large datasets.

**Improvement:** Utilize `storm`'s (or BoltDB's underlying) capabilities to get a direct count of records without loading all data into memory. `storm` typically provides a `Count()` method for this purpose.

### 2.2. Morse Audio Caching

**Current State:** The `newMorseAudio` function in the `morse` package, which pre-generates audio samples, is called every time `GenerateWav` is invoked. While pre-generation is good, re-creating the `morseAudio` instance for every request is inefficient.

**Improvement:** Cache the `morseAudio` instance within the `morse` package or the `Server` struct. If WPM and frequency are constant, the `morseAudio` instance can be initialized once and reused across requests.

## 3. Testing

### 3.1. Expand Test Coverage

**Current State:** Tests are primarily focused on `visitorstore` (`visitorStore_test.go`). There are no explicit tests for HTTP handlers in `main.go` or the Morse code generation logic in `morse.go`.

**Improvement:** Add comprehensive unit and integration tests:
-   **HTTP Handlers:** Test each HTTP handler in `main.go` to ensure they correctly process requests, render templates, handle form submissions, and return appropriate responses/errors.
-   **Morse Code Generation:** Add tests for `morse.go` to verify that Morse code generation is accurate and WAV files are correctly formatted.

## 4. Development Workflow

### 4.1. Frontend Dependency Management

**Current State:** `@tailwindcss/cli` is listed under `dependencies` in `package.json`, but it's typically a development tool.

**Improvement:** Move `@tailwindcss/cli` from `dependencies` to `devDependencies` in `package.json` to accurately reflect its role as a development-time tool.

## 5. Future Enhancements

### 5.1. User Authentication/Authorization

**Consideration:** For more controlled environments or if the application were to be exposed beyond a local kiosk, adding basic authentication for the `/list` or administrative pages might be beneficial.

### 5.2. Data Export UI

**Consideration:** While the `README.md` mentions post-event data export, providing a simple UI within the application to trigger CSV export (especially for BoltDB) would enhance usability.

### 5.3. Internationalization (i18n)

**Consideration:** If the application is intended for use in multiple languages, implementing internationalization would allow for easy translation of the UI.
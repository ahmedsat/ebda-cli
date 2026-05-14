# EBDA CLI

`ebda-cli` is a powerful command-line interface and management tool designed for handling farm data, syncing audits, and validating geospatial information between **KoboToolbox** and **Frappe (ERPNext)**.

It provides a suite of tools for farm managers and engineers to automate data flows, perform geospatial validation (overlap and area checks), and generate comprehensive reports.

## Features

-   **Data Synchronization:** Seamlessly sync data between KoboToolbox and Frappe.
-   **Geospatial Validation:** 
    -   Validate farm map areas against recorded data.
    -   Detect overlapping farm polygons to ensure data integrity.
-   **Reporting:** Generate totals for farms, farmers, and areas by region.
-   **Audit Management:** Manage PGS (Participatory Guarantee Systems) and follow-up audits.
-   **Soil Analysis:** Track and manage soil analysis records.
-   **Training Module:** Includes a specialized training management system with its own compiler and VM.
-   **Multiple Interfaces:**
    -   **CLI:** Primary interface for automation and quick tasks.
    -   **GUI:** Experimental desktop interface built with [Fyne](https://fyne.io/).
    -   **Web UI:** Integrated web dashboard.
-   **Notifications:** Desktop alerts for task completion and errors.

## Prerequisites

-   **Go:** Version 1.26.2 or higher.
-   **SQLite:** Used for local data caching and management.
-   **CGO:** Required for the GUI and certain dependencies (like PDF processing).

## Installation

Clone the repository and build the executable:

```bash
git clone https://github.com/ahmedsat/ebda-cli.git
cd ebda-cli
go build -v -tags=release -o ebda-cli .
```

For cross-compilation (e.g., for Windows), you can use:

```bash
CGO_ENABLED=0 GOOS=windows go build -v -tags=release -o ebda-cli.exe .
```

Note: Certain features like the experimental GUI may require `CGO_ENABLED=1` and proper headers for Fyne dependencies.

## Configuration

The tool is configured via environment variables. Ensure the following are set before running:

| Variable | Description |
| :--- | :--- |
| `ERP_USERNAME` | Your Frappe/ERPNext username. |
| `ERP_PASSWORD` | Your Frappe/ERPNext password. |
| `ERP_BASE_URL` | The base URL for your ERPNext instance. |
| `KOBO_AUTH_TOKEN` | KoboToolbox API authentication token. |
| `KOBO_BASE_URL` | KoboToolbox API base URL. |
| `DB_EBDA_CLI_FILE_PATH` | Path to the local SQLite database file. |
| `DISABLE_NOTIFY` | Set to `true` to disable desktop notifications. |

## Usage

Run the executable with a subcommand:

```bash
./ebda-cli [subcommand] [options]
```

### Key Subcommands

-   **`update`**: The primary command for running various sync stages.
    -   Options: `--skip-new-farms`, `--skip-follow-up`, `--skip-pgs`, `--skip-soils`, `--skip-maps`.
-   **`totals`**: Generates reports on farm and farmer counts and total area.
-   **`pgs`**: Handles PGS audit synchronization.
-   **`map`**: Performs geospatial validation on farm maps.
-   **`soil`**: Manages soil analysis data.
-   **`gui`**: Launches the experimental desktop GUI.
-   **`web`**: Launches the web-based dashboard.
-   **`training`**: Manages training records and logic.

For a full list of commands and their descriptions, run:
```bash
./ebda-cli help
```

## Development

The project is structured into several modules:
-   `commands/`: Implementation of CLI subcommands.
-   `frappe/`: API client for Frappe/ERPNext integration.
-   `kobo/`: API client for KoboToolbox integration.
-   `geo/`: Geospatial processing logic.
-   `services/`: High-level business logic and report generation.
-   `utils/`: Shared utilities including the `SyncRunner` for concurrent tasks.

### Testing

Run tests using the Go toolchain:
```bash
go test ./...
```

## License

[Add License Information Here, e.g., MIT]

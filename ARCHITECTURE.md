# ctop Architecture

## 1. Application Overview and Objectives

`ctop` is a command-line tool that provides a concise and real-time overview of container metrics. It acts as a "top-like" interface for containers, allowing users to monitor CPU, memory, network, and I/O usage at a glance directly from their terminal.

The primary objectives of `ctop` are:
- **Real-Time Monitoring:** Provide a live, continuously updated view of container performance metrics.
- **Extensibility:** Support multiple container runtimes (e.g., Docker, runc) through a modular backend system.
- **Interactivity:** Allow users to sort, filter, and manage containers through an intuitive terminal user interface (TUI).
- **Lightweight:** Be a minimal, efficient tool that can run easily in various environments.

## 2. Architecture and Design Choices

`ctop` is built with a modular and concurrent architecture to keep the UI responsive while collecting data from multiple sources in the background.

### Core Components

#### a. Connector Interface
The most critical design choice is the `Connector` interface, which decouples the core application from the container backend. This allows `ctop` to support different container runtimes by providing a specific implementation for each.

- **`Connector` Interface (`connector/main.go`):** Defines the essential methods a backend must provide, such as `All()` to list containers and `Get()` to retrieve a specific container.
- **`ConnectorSuper` (`connector/main.go`):** A wrapper that provides resilient connection logic, including initial connection and automatic retries on failure.
- **Implementations:**
    - `connector/docker.go`: The implementation for the Docker engine.
    - `connector/runc.go`: The implementation for runc.
    - `connector/mock.go`: A mock implementation used for development and testing.

#### b. Data Model (`container/` and `models/`)
The data is structured logically to separate the container's identity from its metrics and metadata.

- **`container.Container` (`container/container.go`):** The central data structure representing a single container. It holds metadata, the latest metrics, and, importantly, references to its specific `Collector` and `Manager`.
- **`models/`:** This package defines the raw data structures for `Metrics` (CPU, memory, etc.) and `Meta` (name, image, state, etc.), ensuring a clean separation of data from logic.

#### c. Data Collection and Management (`collector/` and `manager/`)
Each container's lifecycle and data streams are handled by dedicated components.

- **`Collector` (`connector/collector/`):** Responsible for collecting metrics for a single container. Each connector type has a corresponding collector (e.g., `docker.go`, `runc.go`). Collectors typically run in a dedicated goroutine per container, streaming `models.Metrics` back to the main application via channels.
- **`Manager` (`connector/manager/`):** Provides an interface for performing actions on a container, such as `Start()`, `Stop()`, and `Pause()`.

#### d. Terminal User Interface (TUI)
The TUI is built using the `termui` library and is composed of several custom widgets.

- **`grid.go`:** Manages the main display, which is a grid of containers. It handles layout, redrawing, and refreshing the container list.
- **`cwidgets/`:** Contains all the custom, reusable UI components, such as the compact grid view (`cwidgets/compact/`) and the detailed single-container view (`cwidgets/single/`).
- **`menus.go`:** Defines the logic for interactive menus like Help, Filter, Sort, and Column selection.

### Concurrency Model
`ctop` is heavily concurrent to ensure a non-blocking UI.
- The **main goroutine** handles UI rendering and user input events.
- Each container's **collector runs in its own goroutine**, continuously fetching metrics and sending them back over a channel.
- The active **connector runs an event-watching goroutine** in the background to listen for container events like `start`, `stop`, and `die`, pushing updates to the UI.
- **Channels** are the primary means of communication, used for streaming metrics, signaling the need for a UI refresh, and propagating status updates.

## 3. Command-Line Arguments

`ctop` can be configured at startup using the following command-line flags.

| Flag | Type | Default | Description |
|---|---|---|---|
| `-v` | bool | `false` | Output version information and exit. |
| `-h` | bool | `false` | Display the help dialog and exit. |
| `-f` | string | `""` | Filter containers by name. |
| `-a` | bool | `false` | Show active containers only (by default, all containers are shown). |
| `-s` | string | `""` | Select the container sort field (e.g., `cpu`, `mem`, `name`). |
| `-r` | bool | `false` | Reverse the container sort order. |
| `-i` | bool | `false` | Invert the default colors for the UI. |
| `-connector` | string | `docker` | The container connector to use (e.g., `docker`, `runc`). |

## 4. Examples on How to Use

**Run with default settings (Docker connector, show all containers):**
```bash
ctop
```

**Show only running containers:**
```bash
ctop -a
```

**Filter containers by name (e.g., only show containers with "app" in the name):**
```bash
ctop -f app
```

**Sort containers by CPU usage in descending order:**
```bash
ctop -s cpu -r
```

**Use the runc connector instead of Docker:**
```bash
ctop -connector runc
```

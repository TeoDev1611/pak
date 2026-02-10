# GEMINI.md - Project Context: PAK

## Project Overview
**PAK** is a "Mini StreamYard Local" desktop application. It is designed to function as a local TV studio where the heavy lifting of visual rendering is handled by the frontend, and the backend serves as a lightweight process manager.

- **Primary Goal:** Provide a local streaming studio with support for guests via WebRTC and broadcasting via FFmpeg.
- **Core Technologies:** 
    - **Frontend:** Vue.js (Framework), Pinia (State Management), PixiJS/WebGL (Rendering Engine), WebRTC (P2P Connectivity).
    - **Backend:** Go (Process management, HTTP/WebSocket server).
    - **Bridge:** Wails (Cross-platform desktop app framework).
    - **Infrastructure:** FFmpeg (Streaming), Cloudflare/SSH (Tunnels for remote access).

## Architecture
The project follows a **"Heavy Client / Light Backend"** pattern, with a long-term goal of implementing **Clean Architecture (Hexagonal Architecture)**.

### Layers:
1.  **Presentation (Frontend - Vue.js):** 
    - Manages scenes and layouts.
    - Renders video streams using PixiJS on a `<canvas>`.
    - Captures the final stream via `canvas.captureStream`.
2.  **Bridge (Wails):** Handles bindings and events between Go and JavaScript.
3.  **Core (Backend - Go):** 
    - Manages signaling for WebRTC.
    - Orchestrates external processes (FFmpeg, Tunnels).
4.  **Infrastructure:** External binaries like FFmpeg and SSH/Cloudflared for the actual streaming and connectivity.

### Planned Directory Structure (Go Standard Layout):
- `/cmd`: Entry point of the application.
- `/internal/domain`: Pure entities and business logic.
- `/internal/service`: Application use cases (e.g., `StartStream`).
- `/internal/adapters`: Implementations for external tools (FFmpeg, Wails handlers).
- `/pkg`: Reusable packages.
- `/frontend`: Vue.js source code.

## Building and Running
*Note: The project is currently in the design phase. The following are standard commands for a Wails project.*

### Prerequisites
- **Go** (latest version recommended)
- **Node.js & npm/yarn/pnpm**
- **Wails CLI** (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
- **FFmpeg** installed in the system path.

### Key Commands
- `wails dev`: Starts the application in development mode with hot-reload.
- `wails build`: Compiles the production-ready desktop application.
- `npm install` (inside `/frontend`): Installs frontend dependencies.

## Development Conventions
- **Clean Architecture:** Prioritize independence between business logic and frameworks. Use interfaces to define boundaries.
- **Heavy Client:** Keep visual and rendering logic in the frontend to leverage GPU acceleration via PixiJS.
- **Process Orchestration:** The backend should focus on managing the lifecycle of FFmpeg and networking tunnels.
- **Documentation:** Maintain the `README.md` with architectural updates as the project evolves.

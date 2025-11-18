# SmartFarm Control Protocol (CCNP Project - Group 3)

A distributed smart greenhouse monitoring and control system implementing a custom binary TLV-based protocol for communication between sensor nodes, control panels, and a central server.

## Overview

This project implements a complete smart farming system with three main components:

- **Nodes**: Greenhouse applications with sensors (temperature, humidity, light) and actuators (heaters, fans, windows, lights) that report sensor data and respond to control commands
- **Control Panel**: Terminal UI (TUI) client for monitoring all greenhouses and sending actuator commands
- **Server**: Central coordinator managing node registration, sensor updates, and command forwarding via TCP

The system uses UDP multicast for server discovery and persistent TCP connections for all operational communication.

## Architecture

```
┌─────────────┐         ┌─────────────┐         ┌─────────────┐
│   Node      │◄───────►│   Server    │◄───────►│   Control   │
│ (Greenhouse)│         │   (TCP/UDP) │         │    Panel    │
└─────────────┘         └─────────────┘         └─────────────┘
```

### Communication Flow

1. **Discovery** (UDP multicast `224.0.0.1:9999`)

   - Nodes and control panels broadcast discovery requests
   - Server responds with its TCP address

2. **Registration** (TCP port `6000`)

   - Nodes register sensors and actuators; server assigns NodeID
   - Control panels register and receive current node list

3. **Operation**
   - Nodes send sensor updates → Server forwards to all control panels
   - Control panels send commands → Server forwards to target node(s) → ACK returned
   - Server notifies all clients when nodes connect/disconnect

## Protocol

The system uses a custom binary **TLV (Type-Length-Value)** and **TRLV (Type-RequestID-Length-Value)** protocol. See [`protocol.md`](protocol.md) for the complete specification, including:

- Message types (Discovery, Registration, Commands, Sensor Updates, ACK/ERROR)
- Wire format and encoding rules
- Node/Actuator selectors
- Data type encoding (integer, float, string, boolean)

## Project Structure

```
ccnp-project-gruppe3/
├── cmd/
│   ├── control-panel/   # TUI control panel application (Bubble Tea)
│   ├── node/            # Greenhouse node simulator
│   └── server/          # Central server (UDP discovery + TCP service)
├── internal/
│   ├── entity/          # Domain models (Node, Sensor, Actuator)
│   ├── protocol/
│   │   ├── constants/   # Protocol type codes and constants
│   │   ├── encoding/    # TLV/TRLV encoding/decoding
│   │   ├── messages/    # High-level message types
│   │   └── selectors/   # Node and actuator selector logic
│   └── util/            # Shared utilities (Stack, UDP discovery)
├── protocol.md          # Protocol specification
└── go.mod
```

## Getting Started

### Prerequisites

- Go 1.25+

### Running Locally

1. **Start the server:**

   ```bash
   go run ./cmd/server
   ```

2. **Start a greenhouse node:**

   ```bash
   go run ./cmd/node -config=examples/greenhouse-configs/greenhouse1.json
   ```

3. **Start the control panel:**
   ```bash
   go run ./cmd/control-panel
   ```

## Development

### Code Organization

- **`cmd/`**: Executable entry points for each component
- **`internal/entity/`**: Core domain types (Node, Sensor, Actuator)
- **`internal/protocol/`**: Protocol implementation
  - `encoding/`: Low-level TLV/TRLV binary encoding
  - `messages/`: Message constructors and decoders
  - `selectors/`: Logic for targeting nodes/actuators
  - `constants/`: All protocol type codes
- **`internal/util/`**: Shared helpers (UDP discovery, stack data structure)

### Key Features

- **Greenhouse Simulation**: Nodes simulate temperature dynamics based on outdoor conditions and actuator states (heaters, fans, windows)
- **Concurrent Request Handling**: Server manages multiple TCP connections concurrently
- **TUI Control Panel**: Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for an interactive terminal interface
- **Binary Protocol**: Efficient TLV encoding with extensible type code space

## Testing

Run the full test suite:

```bash
go test ./...
```

Key test packages:

- `internal/protocol/encoding/`: TLV/TRLV encoding/decoding
- `internal/protocol/messages/`: Message construction and validation
- `internal/protocol/selectors/`: Node and actuator selection logic
- `internal/util/`: Utilities and helpers

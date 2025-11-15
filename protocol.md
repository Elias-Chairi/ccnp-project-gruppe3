# Protocol

## Introduction

This document describes the **SmartFarm TRLV Protocol**, a custom application-layer protocol for communication between sensor/actuator nodes, control-panel clients, and a server hub in a smart-farming system.

## Terminology

| Term                      | Definition                                                                                |
| ------------------------- | ----------------------------------------------------------------------------------------- |
| **Node**                  | A sensor/actuator device that reports measurements and executes actuator commands.        |
| **Control Panel**         | Client application used by the farmer to monitor sensor data and send actuator commands.  |
| **Server**                | Central hub managing discovery, registration, sensor updates, and command forwarding.     |
| **TRLV**                  | Type–RequestID–Length–Value binary encoding used for all top-level messages.              |
| **TLV**                   | Type–Length–Value binary encoding used for all message data.                              |
| **RequestID**             | Client-generated identifier to match requests and responses.                              |
| **NodeID**                | Unique identifier assigned by the server to each sensor/actuator node after registration. |
| **SensorID / ActuatorID** | Local identifiers assigned by a node to its individual sensors or actuators.              |
| **Node Selector**         | Field indicating which node(s) a command targets (single, list, or ALL).                  |
| **Actuator Selector**     | Field indicating which actuator(s) a command targets (single, list, type, or ALL).        |
| **Actuator State**        | Action to apply to an actuator (ON/OFF/SET with value).                                   |
| **Data Type**             | Code describing the type of a sensor value (integer, float, string).                      |
| **ACK/ERROR**             | Response message indicating success or error with a command or registration.              |

## Transport and protocol type

- Discovery: UDP multicast `224.0.0.1:9999`
  - To allow for minimal configuration, all nodes and control panels use UDP multicast to discover the server.
  - Best-effort delivery; actors should retry if no response is received.
- Normal operation: TCP `6000`
  - All further communication uses a single TCP connection to the server.
  - Connection-oriented, each actor maintains a persistent TCP connection to the server. The server keeps reading from each connection indefinitely. If the connection is broken (actor closes it or actor sends malformed data), the actor must re-register.
  - Stateful, the server maintains a registry of active nodes, their IDs, and their current sensor values and actuator states.

## Architecture

The actors are:

- **Sensor/Actuator Nodes** – connect to the server and send sensor data/receive actuator commands.
- **Control Panels** – connect to the server and receive updates and send commands.
- **Server** – listens on TCP/UDP, assigns IDs, maintains registry, forwards messages.
  The client actors (nodes and controls) never talk directly to each other; the server is the single TCP endpoint.

## Message format

All messages are encoded in binary using a TLV-based format, with a top-level TLV or TRLV if the message is a request/response pair.

TRLV = Type (1 byte) + RequestID (2 byte) + Length (2 byte) + Value (variable length).

TLV = Type (1 byte) + Length (2 byte) + Value (variable length).

### Code space layout

All 1-byte type codes are partitioned into non-overlapping ranges for clarity and future growth:

- MessageType: 0x40-0x4F (16 IDs)
- NodeSelector: 0x50-0x5F (16 IDs)
- ActuatorSelector: 0x60-0x6F (16 IDs)
- NodeField: 0x30-0x3F (16 IDs)
- SensorField: 0x70-0x7F (16 IDs)
- ActuatorField: 0x80-0x8F (16 IDs)
- DataType: 0x90-0x9F (16 IDs)
- AckErrorCode: 0xA0-0xAF (16 IDs)

### Top-level Message Type codes:

#### Discovery

**Purpose**: Allow nodes and control panels to discover the server's TCP address.

**Top-level**: TLV, does not need request/response matching since every response is identical.
| Type (hex) | Direction | Expected value | Response |
| ------------------ | --------------------------- | -------------- | -------- |
| `0x41` – DISCOVERY | Node/Control → Server (UDP) | | ACK |

#### Registration

**Purpose**: Actors establish an indefinitely lasting TCP connection.

- Nodes register their sensors and actuators, and get assigned a NodeID.
- Control panels register, and receive the current list of all registered nodes.

**Top-level**: TLV, does not need request/response matching since it does not make sense to write multiple registrations on the same connection.

| Type (hex)                | Direction              | Expected value                        | Response                      |
| ------------------------- | ---------------------- | ------------------------------------- | ----------------------------- |
| `0x42` – REGISTER_NODE    | Node → Server (TCP)    | [List of Actuator and Sensor entries] | ACK[SINGLE_NODE]/ERR          |
| `0x43` – REGISTER_CONTROL | Control → Server (TCP) |                                       | ACK[List of Node Entries]/ERR |

#### Sensor Update

**Purpose**: Nodes send updated sensor values to the server, which forwards them to all connected control panels.

**Top-level**: TLV, does not need request/response matching since updates do not expect a reply.

| Type (hex)               | Direction              | Expected value              | Response |
| ------------------------ | ---------------------- | --------------------------- | -------- |
| `0x44` – SENSOR_UPDATE   | Node → Server (TCP)    | [Sensor entry]              |          |
| `0x44` – SENSOR_UPDATE   | Server → Control (TCP) | [SINGLE_NODE][Sensor entry] |          |
| `0x45` – ACTUATOR_UPDATE | Server → Control (TCP) | [Actuator entry]            |          |
| `0x46` – NODE_ADDED      | Server → Control (TCP) | [Node entry]                |          |
| `0x47` – NODE_REMOVED    | Server → Control (TCP) | [SINGLE_NODE]               |          |

#### Command

**Purpose**: Control panels send commands to the server, which forwards them to the target nodes.

**Top-level**: TRLV, request/response matching is needed in the case of multiple commands being sent on the same connection. Either by:

- One control panel sending multiple commands
- The server forwarding multiple commands to one node.

| Type (hex)       | Direction              | Expected value                                  | Response            |
| ---------------- | ---------------------- | ----------------------------------------------- | ------------------- |
| `0x48` – COMMAND | Control → Server (TCP) | [NodeSelector][ActuatorSelector][ActuatorState] | ACK/ERROR_REQUESTID |
| `0x48` – COMMAND | Server → Node (TCP)    | [ActuatorSelector][ActuatorState]               | ACK/ERROR_REQUESTID |

#### Acknowledgment and Error

**Purpose**: Acknowledge successful processing of a command or registration, or report an error.

**Top-level**:

- `ACK/ERROR_REQUESTID`: TRLV, used for commands that need request/response matching.
- `ACK/ERROR`: TLV, used for registrations that do not need request/response matching

| Type (hex)                   | Direction                   | Expected value | Response |
| ---------------------------- | --------------------------- | -------------- | -------- |
| `0x49` – ACK/ERROR           | Server → Node/Control (TCP) | [ACK/ERROR]    |          |
| `0x50` – ACK/ERROR_REQUESTID | Node/Server → Sender (TCP)  | [ACK/ERROR]    |          |

### Node Selectors codes

|    Hex | Selector        | Example Value      | Meaning                 |
| -----: | --------------- | ------------------ | ----------------------- |
| `0x51` | **SINGLE_NODE** | `0x07`             | Target Node 7 only.     |
| `0x52` | **NODE_LIST**   | `[0x07,0x0C,0x13]` | Specific list of nodes. |
| `0x53` | **ALL_NODES**   | none               | Broadcast to all nodes. |

### Actuator Selectors codes

|    Hex | Selector            | Example Value      | Meaning                          |
| -----: | ------------------- | ------------------ | -------------------------------- |
| `0x61` | **SINGLE_ACTUATOR** | `0x03`             | Actuator ID 3 only.              |
| `0x62` | **ACTUATOR_LIST**   | `[0x01,0x04,0x05]` | Specific list of actuators.      |
| `0x63` | **ACTUATOR_TYPE**   | `FAN` (string)     | All actuators of type FAN.       |
| `0x64` | **ALL_ACTUATORS**   | none               | All actuators on target node(s). |

### Node Field codes

**Node Entry**: code `0x30`

| Field Code | Data Type | Meaning                                                             |
| ---------: | --------- | ------------------------------------------------------------------- |
|     `0x31` | 1 byte    | **Local ID** within server.                                         |
|     `0x32` |           | **Sensor Entries** (inner [data type TLV](#sensor-field-codes))     |
|     `0x33` |           | **Actuator Entries** (inner [data type TLV](#actuator-field-codes)) |

### Sensor Field codes

**Sensor Entry**: code `0x70`

| Field Code | Data Type    | Meaning                                             |
| ---------: | ------------ | --------------------------------------------------- |
|     `0x71` | 1 byte       | **Local ID** within node.                           |
|     `0x72` | UTF-8 string | **Type** (e.g., TEMPERATURE, HUMIDITY, LIGHT).      |
|     `0x73` | UTF-8 string | **Unit** (°C, %, lux).                              |
|     `0x74` |              | **Value** (inner [data type TLV](#data-type-codes)) |

### Actuator Field codes

**Actuator Entry**: code `0x80`

| Field Code | Data Type    | Meaning                                             |
| ---------: | ------------ | --------------------------------------------------- |
|     `0x81` | 1 byte       | **Local ID** within node.                           |
|     `0x82` | UTF-8 string | **Type** (e.g., FAN, HEATER, WINDOW).               |
|     `0x83` | UTF-8 string | **Unit** (°C, %, rpm).                              |
|     `0x84` |              | **State** (inner [data type TLV](#data-type-codes)) |

### Data type codes

| Type (hex) | Meaning | Value Encoding                     |
| ---------- | ------- | ---------------------------------- |
| `0x91`     | Integer | Variable length, big-endian        |
| `0x92`     | Float   | Variable length, big-endian        |
| `0x93`     | String  | Variable length, UTF-8 encoded     |
| `0x94`     | Boolean | 1 byte (0x00 = false, 0x01 = true) |

## ACK/ERR codes

Acknowledgment and error code should be the first byte of the Value field in an ACK/ERROR message. Extra data (e.g., error description) may follow. Extra data is interpreted as a UTF-8 string.

| Code (hex) | Meaning                   |
| ---------- | ------------------------- |
| `0xA0`     | ACK: success              |
| `0xA1`     | ERR: Unknown NodeID       |
| `0xA2`     | ERR: Unknown SensorID     |
| `0xA3`     | ERR: Unknown ActuatorID   |
| `0xA4`     | ERR: Invalid Action       |
| `0xA5`     | ERR: Invalid value        |
| `0xA6`     | ERR: Malformed message    |
| `0xA7`     | ERR: Invalid message type |

## Supported actuators

| Type   | Unit                            | Description                                                  |
| ------ | ------------------------------- | ------------------------------------------------------------ |
| FAN    | 'RPM' -> integer, '' -> boolean | The rpm of the fan or on/off with a default speed            |
| HEATER | '°C' -> float32, '' -> boolean  | Target temperature or on/off with a default temperature      |
| WINDOW | '' -> float32 / boolean         | Openness percentage (0.0-1.0) or fully open/closed           |
| LIGHT  | 'lx' -> float32, '' -> boolean  | Target brightness in lux or on/off with a default brightness |

## Example

### Control Panel Connects and Receives All Current Nodes

1. Discovery (UDP)

```
Type: 0x41 DISCOVERY
Length: 0
Value: (empty)
```

Server replies:

```
Type: 0x46 ACK/ERROR
Length: N
Value: "192.168.0.10:6000"
```

2. Registration (TCP)

```
Type: 0x43 REGISTER_CONTROL
Length: 0
Value: (empty)
```

Server replies:

```
Type: 0x46 ACK/ERROR
Length: N
Value:
  [NodeList TLV]
    [NodeEntry TLV]
      Local ID=0x07
      [SensorEntry TLV]
        Local ID=0x08
        Type=[LIGHT]
        Unit=[lux]
        Value=[Float TLV: 350.5]
      [ActuatorEntry TLV]
        Local ID=0x01
        Type=[FAN]
        Unit=[rpm]
        State=0x02 (OFF)
```

### Control Panel Sends a Command

“Turn ON all fans at Nodes 7, 12, 19”

```
Type: 0x45 COMMAND
Length: N
  [NodeSelector TLV]
    SelectorType = 0x52 NODE_LIST
      Value = [0x07, 0x0C, 0x13]
  [ActuatorSelector TLV]
    SelectorType = 0x63 ACTUATOR_TYPE
      Value = "FAN"
  [ActuatorState TLV]
      State = 0x01 (ON)
      Value = (empty)
```

## Reliability

TCP ensures reliable delivery of all messages after the initial discovery phase.

## Security

Not implemented.

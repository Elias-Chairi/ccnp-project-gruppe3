# Protocol

## Introduction

This document describes the **SmartFarm TLV Protocol**, a custom application-layer protocol for communication between sensor/actuator nodes, control-panel clients, and a server hub in a smart-farming system.

## Terminology

| Term                      | Definition                                                                                |
| ------------------------- | ----------------------------------------------------------------------------------------- |
| **Node**                  | A sensor/actuator device that reports measurements and executes actuator commands.        |
| **Control Panel**         | Client application used by the farmer to monitor sensor data and send actuator commands.  |
| **Server**                | Central hub managing discovery, registration, sensor updates, and command forwarding.     |
| **TLV**                   | Type–Length–Value binary encoding used for all message payloads.                          |
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
  - Connection-oriented, to ensure delivery and ordering.
  - Stateful, the server maintains a registry of active nodes, their IDs, and their current sensor values and actuator states.

## Architecture

The actors are:

- **Sensor/Actuator Nodes** – connect to the server and send sensor data/receive actuator commands.
- **Control Panels** – connect to the server and receive updates and send commands.
- **Server** – listens on TCP/UDP, assigns IDs, maintains registry, forwards messages.
  The client actors (nodes and controls) never talk directly to each other; the server is the single TCP endpoint.

## Message format

All messages use TLV encoding; TLV = Type (1 byte) + Length (2 byte) + Value (variable length).

### Code space layout

All 1-byte type codes are partitioned into non-overlapping, nibble-aligned ranges for clarity and future growth:

- MessageType: 0x40-0x4F (16 IDs)
- NodeSelector: 0x50-0x5F (16 IDs)
- ActuatorSelector: 0x60-0x6F (16 IDs)
- SensorField: 0x70-0x7F (16 IDs)
- ActuatorField: 0x80-0x8F (16 IDs)
- DataType: 0x90-0x9F (16 IDs)
- AckErrorCode: 0xA0-0xAF (16 IDs)

### Header TLV:

| Type (hex)                | Direction                   | Expected value                                  | Response             |
| ------------------------- | --------------------------- | ----------------------------------------------- | -------------------- |
| `0x41` – DISCOVERY        | Node/Control → Server (UDP) |                                                 | ACK[TCP Address]     |
| `0x42` – REGISTER_NODE    | Node → Server (TCP)         | [List of Actuator and Sensor entries]           | ACK[SINGLE_NODE]/ERR |
| `0x43` – REGISTER_CONTROL | Control → Server (TCP)      |                                                 | ACK[NODE_LIST]       |
| `0x44` – SENSOR_UPDATE    | Node → Server (TCP)         | [Sensor entry]                                  |                      |
| `0x44` – SENSOR_UPDATE    | Server → Control (TCP)      | [SINGLE_NODE][Sensor entry]                     |                      |
| `0x45` – COMMAND          | Control → Server (TCP)      | [NodeSelector][ActuatorSelector][ActuatorState] | ACK/ERR              |
| `0x45` – COMMAND          | Server → Node (TCP)         | [ActuatorSelector][ActuatorState]               | ACK/ERR              |
| `0x46` – ACK/ERROR        | Node/Server → Sender (TCP)  | [ACK/ERROR]                                     |                      |

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

### Sensor Entry codes; Field code `0x70`: Sensor Entry

| Field Code | Data Type    | Meaning                                             |
| ---------: | ------------ | --------------------------------------------------- |
|     `0x71` | 1 byte       | **Local ID** within node.                           |
|     `0x72` | UTF-8 string | **Type** (e.g., TEMPERATURE, HUMIDITY, LIGHT).      |
|     `0x73` | UTF-8 string | **Unit** (°C, %, lux).                              |
|     `0x74` |              | **Value** (inner [data type TLV](#data-type-codes)) |

### Actuator Entry codes; Field code `0x80`: Actuator Entry

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

| Code (hex) | Meaning                 |
| ---------- | ----------------------- |
| `0xA0`     | ACK: success            |
| `0xA1`     | ERR: Unknown NodeID     |
| `0xA2`     | ERR: Unknown SensorID   |
| `0xA3`     | ERR: Unknown ActuatorID |
| `0xA4`     | ERR: Invalid Action     |
| `0xA5`     | ERR: Invalid value      |
| `0xA6`     | ERR: Malformed message  |
| `0xA7`     | ERR: Invalid message type|

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

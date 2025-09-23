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

### Header TLV:

| Type (hex)                | Direction                   | Expected value                                  | Response             |
| ------------------------- | --------------------------- | ----------------------------------------------- | -------------------- |
| `0x1` – DISCOVERY         | Node/Control → Server (UDP) |                                                 | ACK[TCP Address]     |
| `0x10` – REGISTER_NODE    | Node → Server (TCP)         | [NODE_LIST][ACTUATOR_LIST]                      | ACK[SINGLE_NODE]/ERR |
| `0x11` – REGISTER_CONTROL | Control → Server (TCP)      |                                                 | ACK[NODE_LIST]       |
| `0x20` – SENSOR_UPDATE    | Node → Server (TCP)         | [SensorEntry]                                   |                      |
| `0x20` – SENSOR_UPDATE    | Server → Control (TCP)      | [SINGLE_NODE][SensorEntry]                      |                      |
| `0x30` – COMMAND          | Control → Server (TCP)      | [NodeSelector][ActuatorSelector][ActuatorState] | ACK/ERR              |
| `0x30` – COMMAND          | Server → Node (TCP)         | [ActuatorSelector][ActuatorState]               | ACK/ERR              |
| `0x40` – ACK/ERROR        | Node/Server → Sender (TCP)  | [ACK/ERROR]                                     |                      |

### Actuator State Codes:

| Code (hex) | Meaning                                        |
| ---------- | ---------------------------------------------- |
| `0x01`     | ON                                             |
| `0x02`     | OFF                                            |
| `0x03`     | Value (inner [data type TLV](#data-type-tlvs)) |

### Node Selectors codes

|    Hex | Selector        | Example Value      | Meaning                 |
| -----: | --------------- | ------------------ | ----------------------- |
| `0x01` | **SINGLE_NODE** | `0x07`             | Target Node 7 only.     |
| `0x02` | **NODE_LIST**   | `[0x07,0x0C,0x13]` | Specific list of nodes. |
| `0x03` | **ALL_NODES**   | none               | Broadcast to all nodes. |

### Actuator Selectors codes

|    Hex | Selector            | Example Value      | Meaning                          |
| -----: | ------------------- | ------------------ | -------------------------------- |
| `0x11` | **SINGLE_ACTUATOR** | `0x03`             | Actuator ID 3 only.              |
| `0x12` | **ACTUATOR_LIST**   | `[0x01,0x04,0x05]` | Specific list of actuators.      |
| `0x13` | **ACTUATOR_TYPE**   | `FAN` (string)     | All actuators of type FAN.       |
| `0x14` | **ALL_ACTUATORS**   | none               | All actuators on target node(s). |

### Sensor Entry codes

| Field Code | Data Type    | Meaning                                            |
| ---------: | ------------ | -------------------------------------------------- |
|     `0x21` | 1 byte       | **Local ID** within node.                          |
|     `0x22` | UTF-8 string | **Type** (e.g., TEMPERATURE, HUMIDITY, LIGHT).     |
|     `0x23` | UTF-8 string | **Unit** (°C, %, lux).                             |
|     `0x24` |              | **Value** (inner [data type TLV](#data-type-tlvs)) |

### Actuator Entry codes

| Field Code | Data Type    | Meaning                                     |
| ---------: | ------------ | ------------------------------------------- |
|     `0x31` | 1 byte       | **Local ID** within node.                   |
|     `0x32` | UTF-8 string | **Type** (e.g., FAN, HEATER, WINDOW).       |
|     `0x33` | UTF-8 string | **Unit** (°C, %, rpm).                      |
|     `0x34` | 1 byte       | **[Actuator State](#actuator-state-codes)** |

### Data type codes

| Type (hex) | Meaning |
| ---------- | ------- |
| `0x01`     | Integer |
| `0x02`     | Float   |
| `0x03`     | String  |

## ACK/ERR codes

| Code (hex) | Meaning                 |
| ---------- | ----------------------- |
| `0x00`     | ACK: success            |
| `0x01`     | ERR: Unknown NodeID     |
| `0x02`     | ERR: Unknown SensorID   |
| `0x03`     | ERR: Unknown ActuatorID |
| `0x04`     | ERR: Invalid Action     |
| `0x05`     | ERR: Invalid value      |

## Example

### Control Panel Connects and Receives All Current Nodes

1. Discovery (UDP)

```
Type: 0x01 DISCOVERY
Length: 0
Value: (empty)
```

Server replies:

```
Type: 0x00 ACK
Length: N
Value: "192.168.0.10:6000"
```

2. Registration (TCP)

```
Type: 0x11 REGISTER_CONTROL
Length: 0
Value: (empty)
```

Server replies:

```
Type: 0x00 ACK
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
Type: 0x30 COMMAND
Length: N
  [NodeSelector TLV]
      SelectorType = 0x02 NODE_LIST
      Value = [0x07, 0x0C, 0x13]
  [ActuatorSelector TLV]
      SelectorType = 0x13 ACTUATOR_TYPE
      Value = "FAN"
  [ActuatorState TLV]
      State = 0x01 (ON)
      Value = (empty)
```

## Reliability

TCP ensures reliable delivery of all messages after the initial discovery phase.

## Security

Not implemented.

# Home Assistant CLI (`ha-cli`)

A powerful, fast, Go-based Command Line Interface to interact with your Home Assistant smart home instance. `ha-cli` allows you to control devices, fetch states, run scripts, and query services directly from your terminal. 

## Features

- **Blazing Fast:** Written in Go, making it extremely lightweight and responsive.
- **Interactive Configuration:** Effortlessly set up your connection using an interactive shell prompt.
- **Full Control:** Control lights, switches, scenes, scripts, automations, and climate devices.
- **Agent Skill Ready:** Designed to be easily integratable with AI Agents (OpenClaw, Hermes) via `agentskills.io` standard.
- **Search & Query:** Quickly query the states of any entity or search for entity IDs globally.
- **Background Daemon:** Run a background sync daemon that periodically caches entity states locally for instant, offline-capable queries.

## Installation

You can install `ha-cli` globally using `go install`:

```bash
go install github.com/leviantech/ha-cli@latest
```

*Ensure your `$(go env GOPATH)/bin` is added to your system `$PATH`.*

## Configuration

`ha-cli` requires your Home Assistant URL and a Long-Lived Access Token to communicate with your instance. 

The easiest way to configure this is via the interactive setup:

```bash
ha-cli config
```
This will prompt you for your **URL** (e.g., `http://192.168.1.100:8123`) and **Token**, securely saving them locally to `~/.ha-cli/config.json`. 

### Alternative Configuration Methods
If you prefer not to use the interactive setup, `ha-cli` will automatically fallback to reading from:
1. Environment variables (`HA_URL` and `HA_TOKEN`).
2. The legacy JSON config file located at `~/.config/home-assistant/config.json`.

*(Note: If `~/.ha-cli/config.json` exists, it will take priority over environment variables.)*

## Usage

Below are the available commands you can run with `ha-cli`:

### Querying States
- `ha-cli state <entity_id>`: Get the current state of an entity.
- `ha-cli states <entity_id>`: Get the full JSON payload of an entity including attributes.
- `ha-cli list [domain]`: List all entity IDs. You can optionally filter by domain (e.g., `ha-cli list lights`).
- `ha-cli search <pattern>`: Search across all entity IDs for a specific keyword.

### Controlling Devices
- `ha-cli on <entity_id> [brightness]`: Turn an entity on. Brightness (0-255) is optional.
- `ha-cli off <entity_id>`: Turn an entity off.
- `ha-cli toggle <entity_id>`: Toggle an entity's power state.
- `ha-cli climate <entity_id> <temp>`: Set the target temperature of a climate device.

### Scenes, Scripts, and Automations
- `ha-cli scene <name>`: Activate a scene.
- `ha-cli script <name>`: Execute a script.
- `ha-cli automation <name>`: Trigger an automation.

### Advanced
- `ha-cli call <domain> <service> [json_data]`: Make a raw service call to Home Assistant with optional JSON payload data.
- `ha-cli info`: Fetch basic instance info from Home Assistant.

### Daemon
The daemon runs in the background and periodically syncs entity states to `~/.ha-cli/entities.json`. When the daemon is active, `list` and `search` commands read from this local cache instead of hitting the API, making them significantly faster.

- `ha-cli daemon start [--interval=<seconds>]`: Start the background sync daemon (default interval: 300s / 5 minutes).
- `ha-cli daemon stop`: Stop the running daemon.
- `ha-cli daemon status`: Check if the daemon is currently running.

## Examples

```bash
# Turn on the living room lights at half brightness
ha-cli on light.living_room 128

# Check if the front door is open
ha-cli state binary_sensor.front_door

# Search for any entity with "temp" in its name
ha-cli search temp

# Call a custom service payload
ha-cli call light turn_on '{"entity_id": "light.bedroom", "color_name": "red"}'

# Start the daemon to sync entities every 60 seconds
ha-cli daemon start --interval=60

# Check daemon status
ha-cli daemon status

# Stop the daemon
ha-cli daemon stop
```

## Agent Skill (OpenClaw / Hermes)
`ha-cli` comes bundled with an `agentskills.io` compatible `SKILL.md` file. By supplying this file to your AI agents (like OpenClaw or Hermes), they will automatically learn how to use `ha-cli` to perform smart home actions on your behalf!

---
name: home-assistant-cli
description: A CLI tool to interact with Home Assistant for controlling smart home devices, executing scenes/scripts, and fetching states.
version: 1.0.0
author: user
---

# Home Assistant CLI Skill

## Overview
This skill allows the agent to control a Home Assistant instance via a local CLI tool named `ha-cli`.

## Setup
Ensure that the `ha-cli` executable is accessible in your environment (it can be installed via `go install github.com/leviantech/ha-cli@latest`).
The environment variables `HA_URL` and `HA_TOKEN` must be set. They can be provided via standard environment variables or configured interactively using the `ha-cli config` command which saves them to `~/.ha-cli/config.json`.

### Getting a Long-Lived Access Token
If you don't have a token, follow these steps in Home Assistant:
1. Open Home Assistant → Profile (bottom left)
2. Scroll to "Long-Lived Access Tokens"
3. Click "Create Token" and name it (e.g., "AI-Agent")
4. Copy the token immediately (it will only be shown once)

## Usage Instructions

You can execute the `ha-cli` program using shell commands to perform various smart home actions. Below are the available commands:

- **Get entity state:** `ha-cli state <entity_id>`
- **Get full entity state (JSON):** `ha-cli states <entity_id>`
- **Turn on an entity:** `ha-cli on <entity_id> [brightness]` (brightness is optional, 0-255)
- **Turn off an entity:** `ha-cli off <entity_id>`
- **Toggle an entity:** `ha-cli toggle <entity_id>`
- **Activate a scene:** `ha-cli scene <scene_name>`
- **Run a script:** `ha-cli script <script_name>`
- **Trigger an automation:** `ha-cli automation <automation_name>`
- **Set climate temperature:** `ha-cli climate <entity_id> <temp>`
- **List entities:** `ha-cli list [domain]` (e.g., `ha-cli list lights`, `ha-cli list all`)
- **Search entities:** `ha-cli search <pattern>`
- **Call a generic service:** `ha-cli call <domain> <service> [json_data]`
- **Capture camera image:** `ha-cli camera <entity_id> [output_file]`
- **List camera entities:** `ha-cli camera list`
- **Get HA instance info:** `ha-cli info`
- **Configure interactively:** `ha-cli config`

### Examples
- Turn on the living room light: `ha-cli on light.living_room`
- Set brightness to 200: `ha-cli on light.living_room 200`
- Activate a scene: `ha-cli scene movie_night`
- Get state of front door: `ha-cli state binary_sensor.front_door`
- Search for a thermostat: `ha-cli search thermostat`
- Capture from a camera: `ha-cli camera front_door snapshot.jpg`

## Troubleshooting
If you encounter errors while using the CLI:
- **401 Unauthorized**: The `HA_TOKEN` is likely expired or invalid. Generate a new one and reconfigure.
- **Connection refused**: Check if `HA_URL` is correct and ensure the Home Assistant server is running and accessible.
- **Entity not found / unknown**: Use `ha-cli list` or `ha-cli search <keyword>` to find the exact `entity_id` before sending commands.

## Best Practices
- Always verify the state of an entity before and after changing it if you need to ensure the action succeeded.
- When searching for entities, use `ha-cli search <keyword>` to find the exact `entity_id` before trying to control it.
- Pay attention to domains; scenes start with `scene.`, scripts with `script.`, and automations with `automation.`. The CLI handles adding these prefixes if omitted, but being explicit is safer.
- **Inbound Events**: If the AI needs to react to Home Assistant events (e.g., motion detected), instruct the user to create a Home Assistant automation that triggers a webhook back to the AI Agent.

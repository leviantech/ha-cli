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
Ensure that the `ha-cli` executable is accessible in your environment (it can be installed via `go install github.com/reyhanfahlevi/ha-cli@latest`).
The environment variables `HA_URL` and `HA_TOKEN` must be set. They can be provided via standard environment variables or configured interactively using the `ha-cli config` command which saves them to `~/.ha-cli/config.json`.

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
- **Get HA instance info:** `ha-cli info`
- **Configure interactively:** `ha-cli config`

### Examples
- Turn on the living room light: `ha-cli on light.living_room`
- Set brightness to 200: `ha-cli on light.living_room 200`
- Activate a scene: `ha-cli scene movie_night`
- Get state of front door: `ha-cli state binary_sensor.front_door`
- Search for a thermostat: `ha-cli search thermostat`

## Best Practices
- Always verify the state of an entity before and after changing it if you need to ensure the action succeeded.
- When searching for entities, use `ha-cli search <keyword>` to find the exact `entity_id` before trying to control it.
- Pay attention to domains; scenes start with `scene.`, scripts with `script.`, and automations with `automation.`. The CLI handles adding these prefixes if omitted, but being explicit is safer.

# AssetFlow Reuploader

Reupload Roblox assets you don't own (animations, sounds, images, meshes) straight to your own account or group, then swap the new ids into your place automatically. A local desktop app plus a thin Roblox Studio plugin.

## Download

Grab the latest from [Releases](https://github.com/NexusAsset/AssetFlow-Reuploader/releases/latest):

- **Windows**: `AssetFlowReuploader-v1.1.0.zip` (or the bare `AssetFlowReuploader.exe`)
- **macOS (beta)**: `assetflow-mac-arm64` (Apple Silicon) or `assetflow-mac-amd64` (Intel)

The app checks for updates on launch and can update itself (and the plugin) in one click.

## Setup

1. Run the app. It opens the AssetFlow dashboard locally.
2. Create an Open Cloud API key at create.roblox.com with `asset:read`, `asset:write`, and `asset-permissions:write`. Paste it into Credentials and Save. (Full steps with screenshots are in the app under Setup & FAQ.)
3. Pick your upload target (your profile or a group).
4. Click **Install plugin**, restart Studio, open the plugin, pick a type, and hit Reupload.

## Your account is safe

Your API key and cookie are stored only on your machine (encrypted at rest on Windows) and are only ever sent to Roblox's official APIs. Nothing is shared. The full source is here so you can verify that.

## Verify your download

Each release ships a `CHECKSUM.txt` with SHA-256 hashes. New unsigned apps can trip Windows SmartScreen or a generic antivirus heuristic; if that happens, choose "Run anyway" and verify the hash.

## Discord

Join the community for help and updates (invite in the app).
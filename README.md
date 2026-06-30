# AssetFlow Reuploader

[![Discord](https://img.shields.io/badge/Discord-Join%20our%20server-5865F2?logo=discord&logoColor=white)](https://discord.gg/j4NPfDwCtA)

Re-upload Roblox **animations, audio, images, and meshes** you don't own to your own account or group, and have the new asset IDs swapped into your place automatically. A clean desktop app plus a lightweight Studio plugin. Open source so you can verify exactly what it does with your credentials.

**Need help? [Join our Discord](https://discord.gg/j4NPfDwCtA).**

## Download

**[Download the latest release](../../releases/latest):**

- **Windows** — `AssetFlowReuploader-v1.1.0.zip` (or the bare `AssetFlowReuploader.exe`)
- **macOS (beta)** — `assetflow-mac-arm64` (Apple Silicon) or `assetflow-mac-amd64` (Intel)

The app checks for updates on launch ("Scanning for updates") and updates itself, and the plugin, in one click.

> Windows may show a SmartScreen warning on first run (the app isn't code-signed yet) — click **More info → Run anyway**.

## Windows flagged the download? (false positive)

Windows Defender / SmartScreen may warn that this app is "a virus or potentially unwanted software." **It is a false positive.** AssetFlow is a Go program that handles your Roblox API key and login cookie *locally* to upload assets — that pattern trips antivirus heuristics even though the code is clean and runs only on your machine. The full source is public here, and you can scan the download yourself on [VirusTotal](https://www.virustotal.com).

To run it:
1. Right-click `AssetFlow Reuploader.exe` → **Properties** → tick **Unblock** → OK; or in **Windows Security → Virus & threat protection → Protection history**, find the item and click **Allow** / **Restore**.
2. Confirm you have the genuine file: its **SHA-256 must match** the value in the release's `CHECKSUM.txt` (see "Verify your download").

## Install

1. Unzip the download anywhere and run **`AssetFlow Reuploader.exe`**.
2. In the app, set up your credentials and target (see Setup below).
3. Click **Install plugin** in the app. This installs the AssetFlow Studio plugin for you and pairs it with the app — **this is the only supported install path** (a manually-copied plugin won't authenticate). Then restart Roblox Studio.

Plugin updates install the same way: when a newer plugin ships, the app shows an **Update plugin** prompt that reinstalls it.

## macOS (beta)

> macOS builds are provided but **not yet tested on real Mac hardware** — please report issues in the Discord. On Mac the interface opens in your **default browser** (no separate window), and credentials are stored locally without OS-level encryption (Windows uses DPAPI; Keychain support is planned).

Open **Terminal** and paste (grabs the right build for your chip, clears the Gatekeeper quarantine, and launches it):

```bash
mkdir -p ~/assetflow && cd ~/assetflow && curl -L -o assetflow "https://github.com/NexusAsset/AssetFlow-Reuploader/releases/latest/download/assetflow-mac-$([ "$(uname -m)" = "arm64" ] && echo arm64 || echo amd64)" && chmod +x assetflow && xattr -dr com.apple.quarantine assetflow && ./assetflow
```

Then click **Install plugin** in the app and restart Roblox Studio.

(Linux isn't supported — Roblox Studio doesn't run on Linux, so the plugin can't either.)

## Setup (one time)

1. Create a Roblox **Open Cloud API key** at <https://create.roblox.com/dashboard/credentials>:
   - Permissions: `asset:read`, `asset:write`, `asset-permissions:write`
   - **Restrict by Creator: OFF**
2. Paste the key into the app, pick your target (your profile or a group), and Save.

The in-app **Setup & FAQ** tab has the full walkthrough with screenshots.

## Use

In Studio, open the **AssetFlow Reuploader** plugin, pick **Animation / Audio / Image / Mesh**, and hit **Reupload**. New IDs swap into your place automatically; watch progress in the app's Activity console.

## Privacy & security

- Runs entirely on your PC and talks only to Roblox.
- Your API key and cookie are **encrypted at rest** (Windows DPAPI) and never leave your machine.
- Never share your API key with anyone.

## Verify your download

Each release ships a `CHECKSUM.txt` with **SHA-256** hashes. To confirm your download wasn't tampered with:

```powershell
Get-FileHash .\AssetFlowReuploader-v1.1.0.zip -Algorithm SHA256
```

The result should match the value in `CHECKSUM.txt`.

## Support

Join the Discord: <https://discord.gg/j4NPfDwCtA>

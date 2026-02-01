# Windows installer (MSI)

This directory contains a WiX-based installer that:

- Installs `smartctl_exporter.exe` to `C:\Program Files\smartctl_exporter`.
- Registers a Windows service (`smartctl_exporter`) that runs automatically.

## Prerequisites

- WiX Toolset v4 (`wix` on PATH).
- Go toolchain.

## Build

```powershell
cd installer/windows
.\build.ps1
```

Artifacts are written to `installer/windows/build/`.

## MSI properties

You can override these defaults during installation (GUI, msiexec, or GPO):

- `LISTEN_ADDRESS` (default `:9633`)

Example:

```powershell
msiexec /i smartctl_exporter.msi LISTEN_ADDRESS=:9663
```

## Service arguments

The service runs with:

```
--config.file="C:\Program Files\smartctl_exporter\smartctl_exporter.yml"
--web.listen-address=<LISTEN_ADDRESS>
```

Any explicit service arguments override the config file defaults.

## smartmontools installation

Install smartmontools separately (for example via GPO). The exporter only needs
`smartctl.exe` to be available in `PATH` or configured via `smartctl.path` in
`smartctl_exporter.yml`.

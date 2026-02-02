# Windows MSI build & deployment (smartctl_exporter)

This guide covers:
- Building the MSI with WiX v4.
- Installing/upgrading via `msiexec`.
- Deploying/upgrading via GPO.

## Prerequisites

- Windows machine with Go toolchain installed.
- WiX Toolset v4 (`wix.exe`) available in `PATH`.

## Build the MSI

From the repository root:

```powershell
cd installer/windows
.\build.ps1
```

Artifacts are written to:

```
installer/windows/build/smartctl_exporter.msi
```

### What the build script does

- Builds `smartctl_exporter.exe`.
- Copies the default config `installer/windows/smartctl_exporter.yml`.
- Runs WiX to produce the MSI.

## Install or upgrade with msiexec

### Fresh install (default port :9633)

```powershell
msiexec /i smartctl_exporter.msi /qn
```

### Install with custom listen address

```powershell
msiexec /i smartctl_exporter.msi /qn LISTEN_ADDRESS=:9663
```

### Upgrade (same MSI, new version)

```powershell
msiexec /i smartctl_exporter.msi /qn REINSTALL=ALL REINSTALLMODE=vomus
```

### Uninstall

```powershell
msiexec /x smartctl_exporter.msi /qn
```

## Deploy / upgrade via GPO

### 1) Prepare a network share

- Place `smartctl_exporter.msi` on a SMB share accessible by target machines.
- Ensure **Domain Computers** have read access to the share.

### 2) Create or edit a GPO

In Group Policy Management:

1. Create a new GPO or edit an existing one.
2. Go to:
   - **Computer Configuration → Policies → Software Settings → Software installation**
3. Right‑click → **New → Package...**
4. Choose the MSI via a UNC path (e.g. `\\server\share\smartctl_exporter.msi`).
5. Select **Assigned**.

### 3) Optional: set MSI properties

If you need to set a custom port (e.g. `LISTEN_ADDRESS`), use **Advanced** when adding the package,
then go to **Modifications** and apply an MST. If you are not using MSTs, you can instead deploy
with a startup script that runs `msiexec` and passes the desired properties.

### 4) Force installation / update

- Run `gpupdate /force` on the target machine or reboot.

### 5) Upgrade with a new MSI

Replace the MSI on the share with the new version and bump the product version in `VERSION`
before building. Then, in GPO:

1. Right‑click the package → **All Tasks → Redeploy application**.

This will reinstall the new MSI on next refresh.

## Notes

- smartmontools is **not bundled**. Install it separately (e.g. via GPO) and ensure `smartctl.exe`
  is in `PATH`, or set `smartctl.path` in `smartctl_exporter.yml`.

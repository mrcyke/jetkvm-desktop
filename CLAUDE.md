# jetkvm-desktop

Native desktop client for JetKVM devices. Go-based, uses video streaming + keyboard/mouse forwarding.

- **Repo:** https://github.com/lkarlslund/jetkvm-desktop
- **Installed version:** v0.1.0-rc8 (2026-04-09)
- **Binary:** `bin/jetkvm-desktop` (macOS aarch64, downloaded from GitHub Releases)
- **Symlinked to:** `/usr/local/bin/jetkvm-desktop`

## Run

```bash
# Launch device discovery UI
jetkvm-desktop

# Connect directly to a device
jetkvm-desktop jetkvm.local
jetkvm-desktop 192.168.1.50

# With password
jetkvm-desktop 192.168.1.50 --password SECRET
```

## Update

Download latest release from GitHub:

```bash
cd ~/WORKSPACE/tools/jetkvm-desktop
gh release download --repo lkarlslund/jetkvm-desktop --pattern "jetkvm-desktop-macos-aarch64.tar.gz" --clobber
tar xzf jetkvm-desktop-macos-aarch64.tar.gz -C bin/
xattr -d com.apple.quarantine bin/jetkvm-desktop
rm jetkvm-desktop-macos-aarch64.tar.gz
```

The `/usr/local/bin/jetkvm-desktop` symlink points to `bin/jetkvm-desktop`, so updates are automatic.

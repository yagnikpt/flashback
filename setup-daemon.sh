#!/bin/bash

set -e

# Determine OS
OS="$(uname -s)"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_NAME="flashbackd"
BINARY_PATH="${SCRIPT_DIR}/${BINARY_NAME}"

if [ ! -f "$BINARY_PATH" ]; then
    echo "Error: ${BINARY_NAME} binary not found in ${SCRIPT_DIR}"
    exit 1
fi

# Function to handle Linux installation
setup_linux() {
    echo "Setting up daemon for Linux..."

    # Create installation directory
    INSTALL_DIR="${HOME}/.local/share/flashback/bin"
    mkdir -p "${INSTALL_DIR}"

    # Stop existing service if running
    systemctl --user stop flashback.service 2>/dev/null || true

    # Wait a moment for the process to fully terminate
    sleep 1

    # Copy binary
    cp "${BINARY_PATH}" "${INSTALL_DIR}/" || {
        echo "Error copying binary. If the error mentions 'Text file busy', please manually stop any running flashback-daemon processes with 'pkill flashback-daemon' and try again."
        exit 1
    }
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

    # Create service file
    SERVICE_DIR="${HOME}/.config/systemd/user"
    mkdir -p "${SERVICE_DIR}"

    SERVICE_FILE="${SERVICE_DIR}/flashback.service"
    cat > "${SERVICE_FILE}" << EOL
[Unit]
Description=Flashback Notification Daemon
After=network-online.target graphical-session.target
Wants=network-online.target graphical-session.target

[Service]
ExecStart=${INSTALL_DIR}/${BINARY_NAME}
Restart=on-failure
RestartSec=5s
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=graphical-session.target
EOL

    # Enable and start the service
    systemctl --user daemon-reload
    systemctl --user enable flashback.service
    systemctl --user start flashback.service

    echo "Flashback daemon installed and started successfully on Linux!"
    echo "Binary location: ${INSTALL_DIR}/${BINARY_NAME}"
    echo "Service file: ${SERVICE_FILE}"
}

# Function to handle macOS installation
setup_macos() {
    echo "Setting up daemon for macOS..."

    # Create installation directory
    INSTALL_DIR="${HOME}/Library/Application Support/flashback/bin"
    mkdir -p "${INSTALL_DIR}"

    # Unload existing service if running
    PLIST_FILE="${HOME}/Library/LaunchAgents/com.flashback.daemon.plist"
    launchctl unload "${PLIST_FILE}" 2>/dev/null || true

    # Wait a moment for the process to fully terminate
    sleep 1

    # Copy binary
    cp "${BINARY_PATH}" "${INSTALL_DIR}/" || {
        echo "Error copying binary. If the error mentions 'Text file busy', please manually kill any running flashback-daemon processes with 'pkill flashback-daemon' and try again."
        exit 1
    }
    chmod +x "${INSTALL_DIR}/${BINARY_NAME}"

    # Create LaunchAgents directory if it doesn't exist
    LAUNCH_AGENTS_DIR="${HOME}/Library/LaunchAgents"
    mkdir -p "${LAUNCH_AGENTS_DIR}"

    # Create plist file
    PLIST_FILE="${LAUNCH_AGENTS_DIR}/com.flashback.daemon.plist"
    cat > "${PLIST_FILE}" << EOL
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.flashback.daemon</string>
    <key>ProgramArguments</key>
    <array>
        <string>${INSTALL_DIR}/${BINARY_NAME}</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardErrorPath</key>
    <string>${HOME}/Library/Logs/flashback-daemon.log</string>
    <key>StandardOutPath</key>
    <string>${HOME}/Library/Logs/flashback-daemon.log</string>
</dict>
</plist>
EOL

    # Load the plist file
    launchctl unload "${PLIST_FILE}" 2>/dev/null || true
    launchctl load -w "${PLIST_FILE}"

    echo "Flashback daemon installed and started successfully on macOS!"
    echo "Binary location: ${INSTALL_DIR}/${BINARY_NAME}"
    echo "Launch Agent plist: ${PLIST_FILE}"
}

# Check if any flashback-daemon processes are running
if pgrep "${BINARY_NAME}" > /dev/null; then
    echo "Warning: ${BINARY_NAME} processes are currently running."
    echo "Attempting to stop them before installation..."
    pkill -f "${BINARY_NAME}" || true
    sleep 2
fi

# Main installation logic
case "${OS}" in
    Linux*)
        setup_linux
        ;;
    Darwin*)
        setup_macos
        ;;
    *)
        echo "Unsupported operating system: ${OS}"
        exit 1
        ;;
esac

echo "Installation complete!"

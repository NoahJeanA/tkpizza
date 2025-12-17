# tiefk-hlpizza

A tool that monitors keyboard input and locks the system when the word 'pizza' is detected.

## Installation

1. Download the latest `tiefk-hlpizza` binary from the [releases page](https://github.com/NoahJeanA/tiefk-hlpizza/releases/latest).
2. Make the binary executable:
    ```bash
    chmod +x tiefk-hlpizza
    ```
3. Move the binary to `/usr/local/bin/`:
    ```bash
    sudo mv tiefk-hlpizza /usr/local/bin/
    ```

## Systemd Service Setup

Since the binary runs as a background system service, it requires a PATH override to successfully execute the lock command across user sessions.

1. **Create the Override Directory:**
    ```bash
    sudo mkdir -p /usr/local/lib/pizza-overrides
    ```

2. **Create the loginctl wrapper:**
    This ensures the binary uses the global `lock-sessions` command.
    ```bash
    sudo bash -c 'cat <<EOF > /usr/local/lib/pizza-overrides/loginctl
#!/bin/bash
/usr/bin/loginctl lock-sessions
EOF'
    sudo chmod +x /usr/local/lib/pizza-overrides/loginctl
    ```

3. **Configure the Service:**
    Create `/etc/systemd/system/pizzalieferant.service`:
    ```ini
    [Unit]
    Description=he will never deliver pizza
    After=network.target

    [Service]
    User=root
    # Ensure our wrapper is found first
    Environment="PATH=/usr/local/lib/pizza-overrides:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin"
    # Replace 1000 with your UID (check with 'id -u')
    Environment=DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/1000/bus
    
    WorkingDirectory=/usr/local/bin/
    ExecStart=/usr/local/bin/tiefk-hlpizza
    
    Restart=always
    RestartSec=5

    [Install]
    WantedBy=multi-user.target
    ```

4. **Enable and Start:**
    ```bash
    sudo systemctl daemon-reload
    sudo systemctl enable --now pizzalieferant.service
    ```

## Usage

Check if the service is running and monitoring inputs:
```bash
sudo systemctl status pizzalieferant.service
journalctl -u pizzalieferant.service -f
```

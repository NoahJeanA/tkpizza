# tiefk-hlpizza

A brief description of your project.

## Installation

1.  Download the latest `tiefk-hlpizza` binary from the [releases page](https://github.com/NoahJeanA/tiefk-hlpizza/releases/latest).
2.  Make the binary executable:
    ```bash
    chmod +x tiefk-hlpizza
    ```
3.  Move the binary to `/usr/local/bin/`:
    ```bash
    sudo mv tiefk-hlpizza /usr/local/bin/
    ```

## Systemd Service

1.  Copy the `pizzalieferant.service` file to `/etc/systemd/system/`:
    ```bash
    sudo cp pizzalieferant.service /etc/systemd/system/
    ```
2.  Reload the systemd daemon to recognize the new service:
    ```bash
    sudo systemctl daemon-reload
    ```
3.  Enable and start the service:
    ```bash
    sudo systemctl enable --now pizzalieferant.service
    ```

## Usage

To check the status of the service:
```bash
sudo systemctl status pizzalieferant.service
```

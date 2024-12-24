# Kumo - PCI DSS Compliance Checker

**Kumo** is a Go-based application designed to check the compliance of your server with the **PCI DSS** (Payment Card Industry Data Security Standard). Built using the **Bubbletea** framework, Kumo provides a terminal-based user interface to quickly and efficiently assess your server's security posture and determine whether it meets PCI DSS requirements.

### Features
- **PCI DSS Compliance Checks:** Validates key security requirements for compliance with the PCI DSS standard, including firewall status, user permissions, file integrity, and more.
- **Terminal-Based UI:** A modern and intuitive terminal user interface, built using the Bubbletea framework.
- **Real-Time Checks:** Performs checks asynchronously and displays results in real-time.
- **Dynamic Loading View:** A loading animation shows during the process of checks, indicating the program is actively evaluating the system.

### PCI DSS Areas Covered:
- **Firewall Configuration:** Verifies that the server has proper firewall configurations to block unauthorized access.
- **Access Control:** Checks that proper access control measures are in place to restrict data access to authorized users.
- **Encryption and Key Management:** Verifies the use of encryption protocols for transmitting sensitive data.
- **System Logging and Monitoring:** Ensures logging and monitoring systems are in place to detect unauthorized access or security breaches.
- **File Integrity Monitoring:** Verifies that the server is monitoring critical files for changes.

### Running the Application

To run the Kumo application, follow these steps:

1. **Build the application:**
   ```sh
   make build
   ```

2. **Run the application:**
   ```sh
   sudo ./kumo
   ```

   Note: The application must be run as root to perform certain system checks.

### Customizing Checks

You can customize the system checks by modifying the `config.json` file. The `config.json` file contains an array of checks, each with the following fields:

- `Name`: The name of the check.
- `Cmd`: The command to execute for the check.
- `ErrHint`: A hint message to display if the check fails.

Example `config.json`:
```json
[
    {
        "Name": "Check Disk Space",
        "Cmd": "df -h",
        "ErrHint": "Ensure there is enough disk space available."
    },
    {
        "Name": "Check Memory Usage",
        "Cmd": "free -m",
        "ErrHint": "Ensure there is enough free memory available."
    },
    {
        "Name": "Check CPU Load",
        "Cmd": "uptime",
        "ErrHint": "Ensure the CPU load is within acceptable limits."
    },
    {
        "Name": "Check Running Processes",
        "Cmd": "ps aux",
        "ErrHint": "Ensure there are no unauthorized processes running."
    }
]
```

### Interpreting Results

The application displays the results of the system checks in the terminal. Each check will show a status of either "Passed" or "Failed". If a check fails, the error hint message will be displayed along with the command output.

- **Passed**: The check was successful.
- **Failed**: The check failed. Review the error hint message and the command output for more information.

Press 'q' to quit the application.

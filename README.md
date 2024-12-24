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

1. **Clone the repository:**
    ```sh
    git clone https://github.com/kintsdev/kumo.git
    cd kumo
    ```

2. **Build the application:**
    ```sh
    make build
    ```

3. **Run the application:**
    ```sh
    sudo ./kumo
    ```

   Note: The application must be run as root to perform certain system checks.

### Customizing Checks

You can customize the system checks by modifying the `config.json` file. This file contains an array of checks, each with the following fields:

- `Name`: The name of the check.
- `Cmd`: The command to be executed for the check.
- `ErrHint`: A hint message to be displayed if the check fails.

Example `config.json` entry:
```json
{
    "Name": "Check Disk Space",
    "Cmd": "df -h",
    "ErrHint": "Ensure there is enough disk space available."
}
```

### Interpreting Results

The results of the system checks are displayed in the terminal with the following format:

- **Passed Checks:** Displayed with a green checkmark (✔) and the message "Passed".
- **Failed Checks:** Displayed with a red cross (✘) and the message "Failed" along with the error hint and the command output.

Example output:
```
✔ Check Disk Space   Passed (df -h)
✘ Check Memory Usage Failed (Ensure there is enough free memory available. (free -m))
```

Press 'q' to quit the application.

# MikroTik and Fortinet Backup Utility

This utility is designed to backup configurations from MikroTik and Fortinet devices. It is intended for use on Linux systems with `sshpass` installed. Make sure that MikroTik and Fortinet devices have SSH enabled.

## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgments](#acknowledgments)

## Getting Started

This utility allows you to automate the backup process for MikroTik and Fortinet devices.

### Prerequisites

Make sure the following prerequisites are met:

1. Linux system
2. `sshpass` installed (you can install it using `sudo apt-get install sshpass`)
3. MikroTik devices with SSH enabled
4. Fortinet devices with SSH enabled, and the configuration moved to an SFTP server.

### Installation

1. Clone the repository:
2. install prerequisites like install sshpass and prepare sftp server
### Usage
1. Create .env file and fill it
2. Build executable like GOARCH=amd64 GOOS=linux go build -o backup main.go
3. chmod +x backup
4. ./backup mikrotik "Backup mikrotik backup rsc file"
5. ./backup fortinet "Backup fortinet conf to sftp server"


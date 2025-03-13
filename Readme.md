# Cypher CLI

Cypher is a secure, open-source password manager that prioritizes your privacy by performing all encryption and decryption operations directly on your device. Your sensitive data never leaves your computer in an unencrypted form.

- Website Version: https://cypher-ochre.vercel.app
- Github: https://github.com/h0i5/Cypher

## Installation

You can use the following automated script to clone the repository, install dependencies, and build the binary:

```sh
#!/bin/bash

# Clone the repository
git clone https://github.com/vijayvenkatj/Cypher-Cli.git
cd Cypher-Cli || exit

# Install dependencies
npm install crypto-js

# Create the hidden config directory
CONFIG_DIR="$HOME/.cypher-cli"
mkdir -p "$CONFIG_DIR"

# Set up environment variable in ~/.cypher-cli/.env
echo 'BACKEND_URL="https://cypher-backend-harshiyers-projects.vercel.app"' > "$CONFIG_DIR/.env"

# Build the binary
go build -o cypher

echo "Cypher CLI setup complete."
echo "Make sure to add the binary to your PATH or run ./cypher from this directory."

```

Alternatively, you can manually perform the installation steps:

1. Clone the repository:
   ```sh
   git clone https://github.com/vijayvenkatj/Cypher-Cli.git
   cd Cypher-Cli
   ```

2. Install dependencies:
   ```sh
   npm install crypto-js
   ```

3. Set up the environment variable:
   Create a `.env` file in the root directory and add the following:
   ```sh
   BACKEND_URL="https://cypher-backend-harshiyers-projects.vercel.app"
   ```

4. Build the binary:
   ```sh
   go build -o cypher
   ```

## Usage

After building the binary, run the following command:
```sh
./cypher
```

### Running `cypher` from Anywhere

Since `cypher` relies on Node.js packages and Go binaries being in the **same directory**, you must run it from its installation folder. Instead of moving files, you can create a **wrapper script** that ensures `cypher` always runs from the correct directory.

1. **Create a new script** (e.g., `cypher.sh`) in `/usr/local/bin/` or `$HOME/.local/bin/`:

   ```bash
   #!/bin/bash
   cd <PATH_TO_CYPHER_CLI> || exit
   ./cypher "$@"
   ```

2. **Make the script executable**:

   ```bash
   chmod +x /usr/local/bin/cypher  # or $HOME/.local/bin/cypher
   ```

3. **Now, you can run `cypher` from anywhere**:

   ```bash
   cypher show
   ```

### Available Commands

- `add`         : Add a password to the vault.
- `decrypt`     : Decrypt a password from the vault.
- `delete`      : Delete a password from the vault.
- `help`        : Display help for any command.
- `login`       : Login using master credentials.
- `register`    : Register using master credentials.
- `show`        : Show all passwords in the vault.

### Flags

- `-l, --email string`                 : Specify the email.
- `-e, --encryption-password string`   : Specify the encryption password.
- `-h, --help`                         : Show help for Cypher.
- `-m, --master-password string`       : Specify the master password for login (DO NOT FORGET THIS PASSWORD).
- `-u, --username string`              : Specify the username for master login.

For more details on a specific command, run:
```sh
Cypher [command] --help
```

## Environment Variables

Ensure you have a `.env` file with the following content:
```sh
BACKEND_URL="<Your backend url here>"
```

## Features

### Security Features

- **256-bit AES Encryption**: Military-grade encryption for all stored credentials
- **PBKDF2 Key Derivation**: Protects against brute-force attacks by making password hashing computationally intensive
- **SHA-256 Hashing**: Ensures data integrity and secure password verification
- **Zero-Knowledge Architecture**: Your data is encrypted before it reaches our servers
- **Client-Side Operations**: All encryption/decryption happens locally on your device
- **No Plaintext Storage**: Sensitive data is never stored in readable form

### How It Works

1. **Master Password**: Users create a strong master password that never leaves their device
2. **Key Derivation**: PBKDF2 generates an encryption key from the master password
3. **Local Encryption**: Passwords and sensitive data are encrypted using AES-256
4. **Secure Storage**: Only encrypted data is synchronized with our servers
5. **Local Decryption**: Data is decrypted on-demand using your master password

### Features

- **Cross-Platform Support**: Available for Windows, macOS, and Linux
- **Browser Extensions**: Seamless integration with major browsers
- **Secure Password Generator**: Create strong, unique passwords
- **Two-Factor Authentication**: Additional security layer for vault access

## Contributing

We welcome contributions! Please read our Contributing Guidelines before submitting pull requests.

# Cool CLI

A small Go CLI tool to manage deployments via the Coolify API.  
Built with [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) for easy configuration.

---

## Features

- Authentication via API URL & Token (`auth` command)
- Stores configuration in `~/.cool.yaml`
- Fetch and select deployments (`deploy` command)

---

## Installation

### Using the install script (recommended)

This script will clone the repo, build the CLI, and move it to `/usr/local/bin`:

```bash
curl -sSL https://raw.githubusercontent.com/ntwcklng/cool/main/script/install.sh | bash
```

---

### Manual installation

1. Install Go (macOS / Linux):

```bash
brew install go
```

2. Clone the repo and build:

```bash
git clone https://github.com/ntwcklng/cool.git
cd cool
go build -o cool
```

3. Optionally move the binary to your PATH:

```bash
sudo mv cool /usr/local/bin/
```

---

## Usage

```bash
cool [command]
```

### Commands

#### `auth`

Authenticate and save your API credentials.

```bash
cool auth
```

- Prompts for:
  - **API URL** (just the base URL, e.g., `example.com`)
  - **Token** (needs write and deploy rights: Keys & Tokens -> API Tokens)
- Saves the values in `~/.cool.yaml`
- Can be rerun anytime to update credentials

---

#### `deploy`

Fetch deployments from the API and select one interactively.

```bash
cool deploy
```

- Loads credentials from your config file
- Automatically runs `auth` if no config is found or credentials are missing
- Displays a numbered list of available deployments
- Prompts the user to select one
- Outputs details of the selected deployment:
  - Application name
  - Deployment UUID
  - Status
  - Deployment URL

---

### Examples

Authenticate:

```bash
cool auth
```

Fetch and select a deployment:

```bash
cool deploy
```

---

You can also get help for any command:

```bash
cool [command] --help
```

- Example:

```bash
cool deploy --help
```

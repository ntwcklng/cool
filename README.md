# Cool CLI

A small Go CLI tool to manage deployments via the Coolify API.  
Built with [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) for easy configuration.

---

## Features

- Authentication via API URL & Token (`auth` command)
- Stores configuration in `~/.cool.yaml`
- Link local projects to deployments (`link` command)
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

#### `link`

Link your local project directory to a specific deployment for easier management.

```bash
cool link
```

- If a `cool.yaml` file exists in the current directory, shows the linked deployment
- If no `cool.yaml` file exists, displays available deployments for selection
- Creates a local `cool.yaml` file with deployment details:
  - DeploymentUUID
  - ApplicationName  
  - FQDN
- Allows project-specific deployment management without repeated selection

---

#### `deploy`

Fetch deployments from the API and trigger a deployment.

```bash
cool deploy
```

- Loads credentials from your config file
- Automatically runs `auth` if no config is found or credentials are missing
- If a `cool.yaml` file exists in the current directory, uses that deployment
- Otherwise, displays a numbered list of available deployments for selection
- Triggers deployment for the selected application
- Shows deployment URL and status

---

#### `update`

Update the Cool CLI to the latest version from GitHub.

```bash
cool update
```

- Checks GitHub for the latest release
- Downloads and builds the new version
- Replaces the current binary automatically
- Shows progress with clear visual feedback

---

### Examples

Authenticate:

```bash
cool auth
```

Link current directory to a deployment:

```bash
cool link
```

Fetch and trigger a deployment:

```bash
cool deploy
```

Update to the latest version:

```bash
cool update
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

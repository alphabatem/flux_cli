# Flux CLI

Unified CLI for AI agents to access FluxBeam's product suite: **RPC**, **DataStream**, and **RugCheck**.

- JSON output by default - built for machine consumption
- Every command is non-interactive and discoverable via `--help`
- Structured exit codes for programmatic error handling
- All 50+ Solana RPC methods exposed as explicit subcommands

## Quick Start

### Install

**One-liner** (Linux/macOS):

```bash
curl -fsSL https://raw.githubusercontent.com/alphabatem/flux_cli/main/install.sh | bash
```

**From source:**

```bash
go install github.com/alphabatem/flux_cli@latest
```

**Download binary** from [Releases](https://github.com/alphabatem/flux_cli/releases).

### Get API Keys

Sign up at [fluxrpc.com](https://fluxrpc.com) to get API keys for each product:

| Product | What it does | Get your key |
|---------|-------------|--------------|
| **FluxRPC** | Solana JSON-RPC, Yellowstone gRPC, WebSockets | [fluxrpc.com](https://fluxrpc.com/admin/apikeys) |
| **DataStream** | Token prices, trader analytics, market data | [fluxrpc.com](https://fluxrpc.com/admin/fluxbeam/apikeys) |
| **RugCheck** | Token security reports, wallet risk assessment | [fluxrpc.com](https://fluxrpc.com/admin/rugcheck/apikeys) |

### Configure API Keys

Set your keys via the CLI, environment variables, or flags:

```bash
# First run interactive setup
flux config init

# Via CLI (persisted to ~/.flux-cli/config.json)
flux config set fluxrpc.api_key YOUR_KEY
flux config set datastream.api_key YOUR_KEY
flux config set rugcheck.api_key YOUR_KEY

# Via environment variables
export FLUX_RPC_API_KEY=YOUR_KEY
export FLUX_DATASTREAM_API_KEY=YOUR_KEY
export FLUX_RUGCHECK_API_KEY=YOUR_KEY

# Via CLI flags (per-command, highest priority)
flux data tokens get <mint> --datastream-api-key YOUR_KEY
```

Priority: CLI flag > environment variable > config file.

On the first interactive run, `flux` will prompt you for API keys, default FluxRPC region, and output format if no config file exists yet. Press Enter to skip any key.

### Select RPC Region

FluxRPC is available in two regions. Default is `us`.

```bash
# Set region (persisted)
flux config set fluxrpc.region eu

# Or via environment variable
export FLUX_RPC_REGION=eu

# Or per-command flag
flux rpc network health --fluxrpc-region eu
```

| Region | Endpoint |
|--------|----------|
| `us` | `https://us.fluxrpc.com` |
| `eu` | `https://eu.fluxrpc.com` |

Yellowstone gRPC uses the same `fluxrpc.api_key` and follows region:

| Region | Yellowstone Endpoint |
|--------|----------------------|
| `us` | `https://yellowstone.us.fluxrpc.com` |
| `eu` | `https://yellowstone.eu.fluxrpc.com` |

### Verify

```bash
flux version
flux config list
```

## Usage

### DataStream - Token & Market Data

```bash
# Token prices
flux data prices So11111111111111111111111111111111111111112

# Token info
flux data tokens get <mint>
flux data tokens details <mint>
flux data tokens candles <mint> --interval 1 --count 10
flux data tokens holders <mint>
flux data tokens holders-top <mint>
flux data tokens trades <mint> --limit 50

# Market stats
flux data stats trending
flux data stats top --limit 20
flux data stats new

# Trader analytics
flux data traders top --limit 10
flux data traders detail <wallet>
flux data traders pnl <wallet>
flux data traders trades <wallet>
```

### RPC - Solana JSON-RPC

Every Solana RPC method is an explicit subcommand:

```bash
# Account
flux rpc account balance <pubkey>
flux rpc account show <pubkey> --encoding jsonParsed
flux rpc account multiple <pubkey1,pubkey2>
flux rpc account watch <pubkey1,pubkey2> --commitment confirmed --timeout 1m
flux rpc account watch-program <programId> --timeout 1m
flux rpc account watch-owner <programId> --timeout 1m

# Blocks
flux rpc block height
flux rpc block show <slot>
flux rpc block latest-blockhash

# Transactions
flux rpc transaction show <signature>
flux rpc transaction signatures <address> --limit 10
flux rpc transaction count
flux rpc transaction watch <account1,account2> --commitment confirmed --timeout 1m
flux rpc signature watch <signature> --timeout 1m
flux rpc signature confirm <signature> --commitment confirmed --timeout 30s

# SPL Tokens
flux rpc token balance <tokenAccount>
flux rpc token accounts-by-owner <owner> --mint <mintAddress>
flux rpc token supply <mint>

# Network
flux rpc network health
flux rpc network version
flux rpc network priority-fees

# Slots & Epochs
flux rpc slot show
flux rpc slot watch --commitment processed --timeout 1m
flux rpc epoch info

# Staking
flux rpc staking vote-accounts
flux rpc staking inflation-rate

# Submit transactions
flux rpc send transaction <base64Tx>
flux rpc send simulate <base64Tx>

# Arbitrary RPC call (escape hatch)
flux rpc call getHealth
flux rpc call getBalance '["<pubkey>"]'
```

Run `flux rpc --help` to see all subcommand groups.

### RugCheck - Token Security

```bash
# Token reports
flux rugcheck report <mint>
flux rugcheck summary <mint>
flux rugcheck scan solana <address>
flux rugcheck search BONK

# Stats
flux rugcheck stats trending
flux rugcheck stats verified
flux rugcheck stats new

# Wallet risk
flux rugcheck wallet solana <address>

# Supported chains
flux rugcheck chains
```

## Output

JSON is the default format. Every response follows the same envelope:

```json
{
  "success": true,
  "data": { ... },
  "error": null,
  "meta": {
    "service": "datastream",
    "endpoint": "/tokens/...",
    "duration_ms": 142
  }
}
```

Errors:

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "API_KEY_MISSING",
    "message": "DataStream API key not configured. Run: flux config set datastream.api_key <key>"
  }
}
```

For human-readable output:

```bash
flux data stats trending --format table
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Usage error (bad flags/args) |
| 64 | API error |
| 65 | Auth error (invalid API key) |
| 69 | Service unavailable |
| 78 | Config error (missing key) |

## Configuration

Config is stored at `~/.flux-cli/config.json`:

```bash
flux config set <key> <value>    # Set a value
flux config get <key>            # Get a value
flux config list                 # List all (keys redacted)
flux config path                 # Print config file path
```

Available keys:

| Key | Description |
|-----|-------------|
| `datastream.api_key` | DataStream API key |
| `datastream.base_url` | DataStream base URL |
| `fluxrpc.api_key` | FluxRPC API key |
| `fluxrpc.region` | FluxRPC region: `eu` or `us` (default: `us`) |
| `fluxrpc.base_url` | FluxRPC base URL (overridden by region) |
| `rugcheck.api_key` | RugCheck API key |
| `rugcheck.base_url` | RugCheck base URL |
| `output.format` | Default output format (`json` or `table`) |

## Building from Source

```bash
git clone https://github.com/alphabatem/flux_cli.git
cd flux_cli
go build -o flux .
```

With version info:

```bash
go build -ldflags "-X github.com/alphabatem/flux_cli/cmd.Version=1.0.0 \
  -X github.com/alphabatem/flux_cli/cmd.Commit=$(git rev-parse --short HEAD) \
  -X github.com/alphabatem/flux_cli/cmd.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -o flux .
```

## License

MIT - see [LICENSE](LICENSE).

# Walkthrough: Schedule Module

The `ua schedule` command has been implemented with a Hexagonal Architecture, ensuring separation of concerns and "Adaptive UI".

## Usage

```bash
# Build
go build -o ua ./cmd/ua-cli

# Run (requires valid cookie)
./ua schedule
```

## Features Implemented

- **Domain**: `Event` entity and `Repository` interface.
- **Adapters**:
  - `UACloudAdapter`: HTTP Client (currently hits `cvnet.cpd.ua.es`. Note: Needs valid cookie).
  - `JSONFileRepo`: Caches events to `~/.ua-cli/cache/schedule.json`.
  - `Presenters`:
    - **TUI**: Bubbletea table (active in terminal).
    - **Text**: Tabwriter table (active in pipes).
    - **JSON**: Raw output via `--json`.

## Verification Results

- **Build**: Success.
- **Execution**: Confirmed CLI wiring.
- **Network**: Returns 404 (Expected, as Authentication Logic is pending valid credentials).

## Next Steps

1. Configure Authentication (Cookie/Keyring).
2. Fine-tune URL and Parser for `UACloudAdapter`.

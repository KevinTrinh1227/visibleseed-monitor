# VisibleSeed Monitor

Uptime monitoring for VisibleSeed services using Uptime Kuma.

**Live:** https://monitor.visibleseed.com

## Setup

1. Clone this repo to the server:
   ```bash
   git clone https://github.com/KevinTrinh1227/visibleseed-monitor.git
   cd visibleseed-monitor
   ```

2. Start the container:
   ```bash
   docker-compose up -d
   ```

3. Access the dashboard at `http://localhost:3001`

## Customization

- `custom-entrypoint.sh` - Rebrands Uptime Kuma to VisibleSeed Monitor
- `custom/custom.css` - Custom CSS styles (if needed)

## Server Location

Running on: `server.visibleseed.com`
Directory: `~/uptime-kuma/`

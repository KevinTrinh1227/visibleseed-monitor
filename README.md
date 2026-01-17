# VisibleSeed Monitor

[![Status](https://monitor.visibleseed.com/api/badge/1/status?label=Status)](https://monitor.visibleseed.com)
[![Uptime 30d](https://monitor.visibleseed.com/api/badge/1/uptime/720?label=Uptime%2030d)](https://monitor.visibleseed.com)
[![Docker](https://img.shields.io/badge/Docker-Uptime%20Kuma-2496ED?logo=docker&logoColor=white)](https://github.com/louislam/uptime-kuma)

Uptime monitoring for VisibleSeed services using [Uptime Kuma](https://github.com/louislam/uptime-kuma).

| | |
|---|---|
| **Dashboard** | https://monitor.visibleseed.com |
| **Status Page** | https://status.visibleseed.com |
| **Server** | `server.visibleseed.com` |

---

## Table of Contents

- [Overview](#overview)
- [Setup](#setup)
- [API Reference](#api-reference)
  - [Public Endpoints](#public-endpoints)
  - [Badge Endpoints](#badge-endpoints)
  - [Push Monitor Endpoint](#push-monitor-endpoint)
  - [Authenticated Endpoints](#authenticated-endpoints)
- [Integration Examples](#integration-examples)
- [Security](#security)
- [Customization](#customization)

---

## Overview

This repository contains the Docker configuration for VisibleSeed's uptime monitoring system. The monitoring dashboard provides real-time status for all VisibleSeed services.

**Key Features:**
- Real-time uptime monitoring
- Public status page and API
- SVG badges for embedding in READMEs
- Push-based and pull-based monitoring
- Notifications via Discord, email, and more

> **Note:** The [status.visibleseed.com](https://status.visibleseed.com) frontend makes API requests to this monitor instance to display service status.

---

## Setup

### Prerequisites

- Docker and Docker Compose installed
- Port 3001 available (or configure a different port)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/KevinTrinh1227/visibleseed-monitor.git
   cd visibleseed-monitor
   ```

2. **Start the container:**
   ```bash
   docker-compose up -d
   ```

3. **Access the dashboard:**

   Open `http://localhost:3001` in your browser.

---

## API Reference

Uptime Kuma provides a REST API for external applications to query monitor status.

**Base URL:** `https://monitor.visibleseed.com`

### Public Endpoints

These endpoints require no authentication.

#### Health Check

```http
GET /api/entry
```

Verify the Uptime Kuma instance is running.

```bash
curl https://monitor.visibleseed.com/api/entry
```

```json
{ "ok": true }
```

#### Status Page Data

```http
GET /api/status-page/{slug}
```

Get the status page configuration and current status for all monitors.

```bash
curl https://monitor.visibleseed.com/api/status-page/visibleseed
```

#### Status Page Heartbeat

```http
GET /api/status-page/heartbeat/{slug}
```

Get the latest heartbeat data for all monitors on a status page.

```bash
curl https://monitor.visibleseed.com/api/status-page/heartbeat/visibleseed
```

**Response includes:**
- `heartbeatList` - Latest heartbeat for each monitor
- `uptimeList` - Uptime percentages (24h and 30d)

### Badge Endpoints

Generate SVG badges for embedding in documentation or dashboards.

| Endpoint | Description |
|----------|-------------|
| `/api/badge/{id}/status` | Current status (Up/Down) |
| `/api/badge/{id}/uptime` | Uptime percentage (24h) |
| `/api/badge/{id}/uptime/{hours}` | Uptime for custom duration |
| `/api/badge/{id}/ping` | Latest response time |
| `/api/badge/{id}/avg-response` | Average response time |
| `/api/badge/{id}/cert-exp` | SSL certificate expiry |

**Examples:**

```bash
# Status badge
curl https://monitor.visibleseed.com/api/badge/1/status

# 30-day uptime badge
curl https://monitor.visibleseed.com/api/badge/1/uptime/720
```

**Markdown embed:**

```markdown
![Status](https://monitor.visibleseed.com/api/badge/1/status)
![Uptime](https://monitor.visibleseed.com/api/badge/1/uptime/720?label=Uptime%2030d)
```

### Push Monitor Endpoint

For push-type monitors where external services report their status.

```http
GET /api/push/{pushToken}?status=up&msg=OK&ping=
```

| Parameter | Description |
|-----------|-------------|
| `status` | `up` or `down` |
| `msg` | Status message (optional) |
| `ping` | Response time in ms (optional) |

```bash
curl "https://monitor.visibleseed.com/api/push/your-push-token?status=up&msg=OK&ping=50"
```

### Authenticated Endpoints

These endpoints require an API key.

#### Creating an API Key

1. Log in to the Uptime Kuma dashboard
2. Go to **Settings > API Keys**
3. Click **Add API Key**
4. Copy the generated key

#### Using API Keys

```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
  https://monitor.visibleseed.com/api/monitors
```

#### Available Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/monitors` | List all monitors |
| GET | `/api/monitors/{id}` | Get single monitor details |

### Status Values

| Value | Meaning |
|-------|---------|
| 0 | Down |
| 1 | Up |
| 2 | Pending |
| 3 | Maintenance |

---

## Integration Examples

### Bash Script

```bash
#!/bin/bash
# check-services.sh - Check VisibleSeed service status

STATUS_PAGE="visibleseed"
API_URL="https://monitor.visibleseed.com/api/status-page"

response=$(curl -s "$API_URL/$STATUS_PAGE")
heartbeat=$(curl -s "$API_URL/heartbeat/$STATUS_PAGE")

echo "Status Page: $(echo $response | jq -r '.config.title')"
echo "Monitors:"
echo $heartbeat | jq -r '.uptimeList | to_entries[] | "\(.key): \(.value)%"'
```

### Python

```python
import requests

BASE_URL = "https://monitor.visibleseed.com"

def get_status_page(slug):
    """Get status page data."""
    response = requests.get(f"{BASE_URL}/api/status-page/{slug}")
    return response.json()

def get_heartbeats(slug):
    """Get heartbeat data for a status page."""
    response = requests.get(f"{BASE_URL}/api/status-page/heartbeat/{slug}")
    return response.json()

def check_service_health():
    """Check if all services are up."""
    data = get_heartbeats("visibleseed")
    for monitor_id, beats in data.get("heartbeatList", {}).items():
        if beats and beats[0]["status"] != 1:
            return False
    return True

if check_service_health():
    print("All services operational")
else:
    print("Some services are down")
```

---

## Security

### Container Security

The Docker configuration includes several hardening measures:

| Feature | Description |
|---------|-------------|
| `no-new-privileges` | Prevents privilege escalation |
| Resource limits | Memory (512MB) and CPU (1 core) caps |
| Read-only mounts | Custom scripts mounted as read-only |
| Localhost binding | Port 3001 bound to `127.0.0.1` only |
| Health checks | Built-in container health monitoring |

### Production Recommendations

1. **Use a reverse proxy** (nginx, Caddy, Traefik) for HTTPS
2. **Enable authentication** on the dashboard
3. **Restrict API access** if not needed publicly
4. **Configure firewall rules** to limit access

### Nginx Configuration

```nginx
server {
    listen 443 ssl http2;
    server_name monitor.visibleseed.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://127.0.0.1:3001;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

---

## Customization

| File | Purpose |
|------|---------|
| `custom-entrypoint.sh` | Rebrands Uptime Kuma to VisibleSeed Monitor |
| `custom/custom.css` | Custom CSS styles (if needed) |

---

## License

This project uses [Uptime Kuma](https://github.com/louislam/uptime-kuma), which is licensed under the MIT License.

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

## API Reference

Uptime Kuma provides a REST API for external applications to query monitor status. The API is available at the base URL where your instance is hosted.

### Public Endpoints (No Authentication Required)

These endpoints are publicly accessible for status pages you've configured:

#### Get Status Page Data

```
GET /api/status-page/{slug}
```

Returns the status page configuration and current status for all monitors on that page.

**Example:**
```bash
curl https://monitor.visibleseed.com/api/status-page/visibleseed
```

**Response:**
```json
{
  "config": {
    "slug": "visibleseed",
    "title": "VisibleSeed Status",
    "description": "Service status for VisibleSeed"
  },
  "incident": [],
  "publicGroupList": [
    {
      "id": 1,
      "name": "Services",
      "monitorList": [...]
    }
  ]
}
```

#### Get Status Page Heartbeat

```
GET /api/status-page/heartbeat/{slug}
```

Returns the latest heartbeat data for all monitors on a status page.

**Example:**
```bash
curl https://monitor.visibleseed.com/api/status-page/heartbeat/visibleseed
```

**Response:**
```json
{
  "heartbeatList": {
    "1": [
      {
        "status": 1,
        "time": "2024-01-15 12:00:00",
        "ping": 45,
        "msg": "OK"
      }
    ]
  },
  "uptimeList": {
    "1_24": 99.95,
    "1_720": 99.87
  }
}
```

#### Get Monitor Badge

```
GET /api/badge/{monitorId}/status
GET /api/badge/{monitorId}/uptime
GET /api/badge/{monitorId}/uptime/{duration}
GET /api/badge/{monitorId}/ping
GET /api/badge/{monitorId}/avg-response
GET /api/badge/{monitorId}/cert-exp
GET /api/badge/{monitorId}/response
```

Returns SVG badges for embedding in READMEs or dashboards.

**Example:**
```bash
# Status badge
curl https://monitor.visibleseed.com/api/badge/1/status

# Uptime badge (24h default)
curl https://monitor.visibleseed.com/api/badge/1/uptime

# Uptime badge (30 days)
curl https://monitor.visibleseed.com/api/badge/1/uptime/720

# Response time badge
curl https://monitor.visibleseed.com/api/badge/1/ping
```

**Embed in Markdown:**
```markdown
![Status](https://monitor.visibleseed.com/api/badge/1/status)
![Uptime](https://monitor.visibleseed.com/api/badge/1/uptime/720?label=Uptime%2030d)
```

#### Push Monitor Endpoint

```
GET /api/push/{pushToken}?status=up&msg=OK&ping=
```

Used for push-type monitors where external services report their status.

**Parameters:**
- `status` - `up` or `down`
- `msg` - Status message (optional)
- `ping` - Response time in ms (optional)

**Example:**
```bash
curl "https://monitor.visibleseed.com/api/push/your-push-token?status=up&msg=OK&ping=50"
```

### Health Check Endpoint

```
GET /api/entry
```

Basic health check endpoint to verify the Uptime Kuma instance is running.

**Example:**
```bash
curl https://monitor.visibleseed.com/api/entry
```

**Response (200 OK):**
```json
{
  "ok": true
}
```

### Authenticated API Endpoints

These endpoints require authentication. You can use either API keys or session-based auth.

#### Create an API Key

1. Log in to the Uptime Kuma dashboard
2. Go to Settings > API Keys
3. Click "Add API Key"
4. Copy the generated key

#### Using API Keys

Include the API key in the `Authorization` header:

```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
  https://monitor.visibleseed.com/api/monitors
```

#### Get All Monitors

```
GET /api/monitors
```

Returns a list of all configured monitors.

**Example:**
```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
  https://monitor.visibleseed.com/api/monitors
```

#### Get Single Monitor

```
GET /api/monitors/{id}
```

Returns details for a specific monitor.

**Example:**
```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
  https://monitor.visibleseed.com/api/monitors/1
```

### Status Values

| Value | Meaning |
|-------|---------|
| 0 | Down |
| 1 | Up |
| 2 | Pending |
| 3 | Maintenance |

### Rate Limiting

The API does not have strict rate limiting by default, but excessive requests may be throttled. For production integrations, cache responses appropriately.

### Example: External Health Check Script

```bash
#!/bin/bash
# check-services.sh - Check VisibleSeed service status

STATUS_PAGE="visibleseed"
API_URL="https://monitor.visibleseed.com/api/status-page"

response=$(curl -s "$API_URL/$STATUS_PAGE")
heartbeat=$(curl -s "$API_URL/heartbeat/$STATUS_PAGE")

# Parse with jq
echo "Status Page: $(echo $response | jq -r '.config.title')"
echo "Monitors:"
echo $heartbeat | jq -r '.uptimeList | to_entries[] | "\(.key): \(.value)%"'
```

### Example: Python Integration

```python
import requests

BASE_URL = "https://monitor.visibleseed.com"

def get_status_page(slug):
    """Get status page data"""
    response = requests.get(f"{BASE_URL}/api/status-page/{slug}")
    return response.json()

def get_heartbeats(slug):
    """Get heartbeat data for a status page"""
    response = requests.get(f"{BASE_URL}/api/status-page/heartbeat/{slug}")
    return response.json()

def check_service_health():
    """Check if all services are up"""
    data = get_heartbeats("visibleseed")
    for monitor_id, beats in data.get("heartbeatList", {}).items():
        if beats and beats[0]["status"] != 1:
            return False
    return True

# Usage
if check_service_health():
    print("All services operational")
else:
    print("Some services are down")
```

## Security

### Container Security

The Docker configuration includes several security hardening measures:

- **no-new-privileges** - Prevents privilege escalation inside the container
- **Resource limits** - Memory (512MB) and CPU (1 core) limits prevent resource exhaustion
- **Read-only mounts** - Custom scripts mounted as read-only
- **Localhost binding** - Port 3001 bound to localhost only (use reverse proxy for external access)
- **Health checks** - Built-in container health monitoring

### Network Security

For production deployments:

1. **Use a reverse proxy** (nginx, Caddy, Traefik) to handle HTTPS
2. **Enable authentication** on the Uptime Kuma dashboard
3. **Restrict API access** if not needed publicly
4. **Configure firewall rules** to limit access to port 3001

### Recommended Nginx Configuration

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

## Customization

- `custom-entrypoint.sh` - Rebrands Uptime Kuma to VisibleSeed Monitor
- `custom/custom.css` - Custom CSS styles (if needed)

## Server Location

Running on: `server.visibleseed.com`
Directory: `~/uptime-kuma/`

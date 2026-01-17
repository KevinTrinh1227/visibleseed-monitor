#!/bin/sh

# Minimal CSS - only rebrand text, don't touch colors/theme
cat > /app/dist/custom.css << 'EOF'
/* Minimal VisibleSeed branding - no theme changes */
EOF

# Just change the title
sed -i 's|<title>Uptime Kuma</title>|<title>VisibleSeed Monitor</title>|g' /app/dist/index.html
sed -i 's|Uptime Kuma</title>|VisibleSeed Monitor</title>|g' /app/dist/index.html

exec /usr/bin/dumb-init -- node server/server.js

#!/bin/sh
set -eu

OUTPUT="/usr/share/nginx/html/env.js"

# Ensure directory exists (should already via COPY of dist)
mkdir -p "$(dirname "$OUTPUT")"

# Initialize the runtime env object
printf '%s\n' "window.__ENV__ = window.__ENV__ || {};" > "$OUTPUT"

# Export all environment variables that start with VITE_
# We escape backslashes, single quotes and newlines for safe JS embedding
printenv | awk -F= '/^VITE_/ {print $1}' | while IFS= read -r name; do
  value=$(printenv "$name" || true)
  echo "Processing $name=$value"
  # Escape for JS single-quoted string
  escaped=$(printf "%s" "$value" | tr -d '\r' | sed -e 's/\\\\/\\\\\\\\/g' -e "s/'/\\\\'/g" -e ':a;N;$!ba;s/\n/\\n/g')
  printf '%s\n' "window.__ENV__['$name']='$escaped';" >> "$OUTPUT"
done

# Show what was generated for debugging
# tail -n +1 "$OUTPUT" || true

# Start nginx in the foreground
exec nginx -g 'daemon off;'

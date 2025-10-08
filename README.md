# Traefik Plugin: Clean X-Forwarded-For

Middleware for Traefik that normalizes the `X-Forwarded-For` (XFF) header for backend services.

- Removes noisy proxy chains by keeping only the first IP (client IP) when desired.
- Optionally sets an empty XFF when header is missing, or removes it entirely.
- Mirrors the normalized value to `X-Client-IP` for convenient backend access.

## Options

- `keepEmpty` (bool, default: `true`):
  - `true` — when XFF is missing, set it to an empty string.
  - `false` — when XFF is missing, remove the header.
- `onlyFirst` (bool, default: `true`):
  - `true` — keep only the first IP (before the first comma).
  - `false` — keep the header as received from the client.

## Traefik Plugin Manifest

This repo includes `.traefik.yml`:

```
type: middleware
name: clean-xff
import: github.com/akrasnov-marfatech/traefik-xff-krasnov-plugin
displayName: Krasnov Clean X-Forwarded-For
summary: Passes client's X-Forwarded-For as-is (first IP). Optionally empty or delete if missing.
testData:
  keepEmpty: true
  onlyFirst: true
```

Ensure the `import` path matches the repository location when publishing/using via Traefik Pilot or local plugin config.

## Example (dynamic config)

```
http:
  middlewares:
    clean-xff@plugin:
      plugin:
        clean-xff:
          keepEmpty: true
          onlyFirst: true
  routers:
    my-router:
      rule: Host(`example.com`)
      service: my-svc
      middlewares:
        - clean-xff@plugin
  services:
    my-svc:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:8080
```

## Behavior

- If `X-Forwarded-For` is present and `onlyFirst` is `true`, the middleware keeps the first IP (trimmed) and sets it back to XFF.
- If the header is missing, behavior is controlled by `keepEmpty`.
- The resulting value is also set into `X-Client-IP`.


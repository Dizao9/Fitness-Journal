# Fitness-Journal API — Frontend & Deployment Integration Guide

Full machine-readable contract: [`docs/openapi.yaml`](./openapi.yaml).

---

## 1. Authentication & Auto-Refresh Flow (HttpOnly Cookies)

The backend uses **HttpOnly, `SameSite=Strict` cookies** — JavaScript never touches JWT tokens directly.

- `access_token` — TTL **30 minutes**
- `refresh_token` — TTL **30 days**
- All protected routes expect `access_token` cookie; all requests from the browser must run under the same origin (via Nginx/Vite proxy) or with `credentials: 'include'`.

### Recommended Frontend Client Lifecycle
1. **App boot / session check**: call `GET /api/athlete/profile`.
   - `200 OK` → user is logged in.
   - `401 Unauthorized` → call `POST /api/auth/refresh` once; if `200 OK`, retry `GET /api/athlete/profile`; if `401`, redirect to `/login`.
2. **Interceptor on any `401`**:
   - Avoid infinite loops: never retry `POST /api/auth/refresh` or `POST /api/auth/login` on `401`.
   - Deduplicate concurrent refresh calls with a single shared promise, then replay failed requests.
3. **Registration**: `POST /api/auth/register` creates the user and returns `201 {"id": "..."}`, **without** setting cookies. Follow it up with `POST /api/auth/login` automatically for smooth UX.

---

## 2. Important API Quirks to Respect on Frontend

| Endpoint | Detail to keep in mind |
|---|---|
| **Error responses** | Most errors are returned via Go `http.Error` as **`text/plain`** (e.g. `"workout already finished\n"`). Exception: `POST /auth/register` validation errors return `400` JSON `{"status": "error", "errors": [...]}`. Parse errors defensively (try JSON, fallback to `res.text()`). |
| `POST /exercise` | Returns `{"id": "42"}` — **`id` is a JSON string**, while `exercise_id` everywhere else is a number (`Number(data.id)`). |
| `POST /workout/{id}/set` | Returns `{"set_id": "<uuid>"}` (field name is `set_id`, not `id`). |
| `PATCH /workout/{id}/finish` | Always send a JSON body (at least `{}`), even when omitting `total_time` and `grade_of_training` — an empty HTTP body triggers `400 invalid request body`. |
| `POST /workout` | `date_of_training` must be an RFC3339 timestamp (`new Date().toISOString()`) and **not in the future**. To start a live session, send `{"date_of_training": "...", "status": "in_progress", "sets": []}`. |
| `POST /workout/{id}/set` validation | `weight >= 0` (0 allowed for bodyweight), `reps >= 0` (0 allowed for failed attempts), `rpe` is `null` or `1..10`. |
| `GET /exercises` | `filter` accepts `"my"`, `"system"`, or omit for all. Always returns a JSON array `[]` (never `null`). |
| `GET /workouts` | `limit` clamped to `1..50` (default `10`), `page >= 1` (default `1`). Always returns `[]` (never `null`). |

---

## 3. Reverse Proxy & Docker Network Contract (Separate Frontend Repo)

This backend repository runs `db`, `migrator`, and `app` (container name **`fitness-app`**, port **`8081`**) inside the named Docker network **`fitness-net`**.

### Production (VPS — separate `fitness-frontend` repo)

In the frontend repository's `docker-compose.yaml`, attach Nginx to the external `fitness-net` network:

```yaml
services:
  web:
    build: .
    container_name: fitness-web
    restart: unless-stopped
    ports:
      - "80:80"
    networks:
      - fitness-net

networks:
  fitness-net:
    external: true
```

Reference `nginx.conf` (serves SPA static build + proxies `/api/*` to Go backend):

```nginx
server {
    listen 80;
    server_name _;

    root /usr/share/nginx/html;
    index index.html;

    # Proxy API requests to Go backend (strips /api prefix)
    location /api/ {
        proxy_pass http://fitness-app:8081/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # SPA client-side routing fallback
    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

> **Cookie `Secure` flag note**:
> - Deploying over plain HTTP (IP address without TLS): set `COOKIE_SECURE=false` in the backend `.env`.
> - Once HTTPS / TLS is enabled on a domain: set `COOKIE_SECURE=true` in the backend `.env`.

### Local Frontend Development (Vite Proxy)

In `vite.config.ts`, proxy `/api` to `http://localhost:8081` (stripping `/api`) so browser cookies stay same-origin with zero CORS setup:

```ts
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8081',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, ''),
    },
  },
}
```

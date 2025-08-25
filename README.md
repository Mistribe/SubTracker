# SubTracker — Track Your Subscriptions like a Pro 💸📅

A fun, modern subscription manager that helps you see where your money goes, 
avoid surprise renewals, and get the most out of your services.

Frontend: React + TypeScript + Vite • Styled with Tailwind CSS and shadcn/ui
Backend: Developed in Go (Golang)
Mobile: Flutter app (alpha)

---

## ✨ Features
- Keep a tidy list of your subscriptions with powerful search and filters
- Create, update, and cancel subscriptions in a couple of clicks
- Providers, plans, and prices—modeled cleanly so you stay organized
- Multi-currency support with live exchange rates
- Labels and categories for fast grouping and reporting
- Family accounts and invitations for shared plans
- A clear dashboard with charts and summaries
- Secure authentication with Kinde

> Built to feel fast, accessible, and delightful.

---

## 🧠 Tech Stack
- React 19.1.0 + TypeScript 5.8.3
- Vite 7.0.4 (blazing-fast dev/build)
- Tailwind CSS 4.1.11 + shadcn/ui (Radix UI under the hood) for accessible, consistent UI
- State and data: React Query (@tanstack/react-query)
- Routing: React Router DOM
- Forms & Validation: React Hook Form + Zod
- Charts: Recharts
- Auth: Kinde Auth React
- Package manager: npm
- Backend: Golang with SQL Database

---

## ☁️ Cloud (Hosted) — Coming Soon
The hosted SubTracker Cloud is in the works. Sign up to get notified when the beta launches. Until then, you can self-host with Docker or run locally from source.

---

## 🚀 Quickstart

1) Clone and install

```
git clone <your-repo-url>
cd SubTracker/web
npm install
```

2) Configure environment variables

Create a `.env` file in `web/` with the following keys:

```
# Backend API
VITE_BACKEND_URL=http://localhost:8080

# Auth (Kinde)
VITE_KINDE_CLIENT_ID=your_kinde_client_id
VITE_KINDE_DOMAIN=your_kinde_domain
VITE_KINDE_REDIRECT_URL=http://localhost:5173
VITE_KINDE_LOGOUT_URL=http://localhost:5173
VITE_KINDE_AUDIENCE=your_kinde_audience

# Optional
VITE_TARGET_ENV=development
```

3) Start the app

```
npm run dev
```

4) Build & preview production

```
npm run build
npm run preview
```

> Note: The backend is a Go (Golang) service. Make sure it’s running and that `VITE_BACKEND_URL` points to it.

## 🏠 Self-hosted (Docker Compose)
Prerequisites: Docker Desktop or Docker Engine 20.10+

1) Create a web/.env file (used at build time by Vite):

```
VITE_BACKEND_URL=http://localhost:8080
VITE_KINDE_CLIENT_ID=your_kinde_client_id
VITE_KINDE_DOMAIN=your_kinde_domain
VITE_KINDE_REDIRECT_URL=http://localhost:9080
VITE_KINDE_LOGOUT_URL=http://localhost:9080
VITE_KINDE_AUDIENCE=your_kinde_audience
```

2) From the repository root, start the stack:

```
docker compose up -d --build
```

- Web UI: http://localhost:9080
- API: http://localhost:8080
- Database (Postgres): localhost:5432

Migrations run via the database.migration service. You can re-run migrations with:

```
docker compose run --rm database.migration
```

To stop the stack:

```
docker compose down
```

---

## 🗂️ Monorepo Structure

```
backend/                 # Go API, database, migrations, compose services
  cmd/api/               # API entrypoint and Dockerfile
  database/              # migrations and init SQL
  internal/ pkg/         # services, domain, utilities
  compose.yml            # api, database, migration services
web/                     # React + Vite frontend
  src/                   # app source
  Dockerfile             # multi-stage build (nginx)
  compose.yml            # web service (port 9080->80)
mobile/                  # Flutter app (Android/iOS)
compose.yml              # root compose orchestrating web+api+db
```

---

## 🛠️ Build from Source

Prerequisites:
- Node.js 20+, npm
- Go 1.25+
- Docker Desktop (optional, for self-hosting)
- Flutter SDK (optional, for mobile)

Web (React + Vite)
```
cd web
npm install
npm run dev       # start dev server at http://localhost:5173
npm run build && npm run preview
# Lint
npm run lint
```

Backend (Go)
```
cd backend
# Ensure DATABASE_DSN points to a reachable Postgres instance if running locally
go run ./cmd/api
# Build
go build -o bin/api ./cmd/api
# Test
go test ./...
```

Mobile (Flutter)
```
cd mobile
flutter pub get
flutter run            # run on connected device/emulator
# Builds
flutter build apk      # Android
flutter build ios      # iOS (on macOS)
```

---

## 🔐 Backend (Go)
The backend powering SubTracker is written in Go and lives in `backend/`.

- API port: `8080`
- Database: PostgreSQL (`5432`) when using Docker Compose
- Docker images: multi-stage build producing a static binary (see `backend/cmd/api/Dockerfile`)
- Environment variables (via Docker or local env):
  - `DATABASE_DSN=host=database user=postgres password=postgres dbname=app port=5432`
  - `UPDATER_AT_START=true`
  - `DATA_LABEL=/data/labels.json`
  - `DATA_FAMILY=/data/families.json`
  - `DATA_PROVIDER=/data/providers.json`
  - `DATA_SUBSCRIPTION=/data/subscriptions.json`

Run locally (without Docker):
```
cd backend
# Ensure a PostgreSQL instance is available and DATABASE_DSN points to it
go run ./cmd/api
```

Build binary:
```
cd backend
go build -o bin/api ./cmd/api
```

Run tests:
```
cd backend
go test ./...
```

With Docker Compose, the API and database are started from the repo root; see the Self-hosted section above.

---

## 🧩 Design System & Styling
- Tailwind CSS utilities for layout/spacing/typography
- shadcn/ui components for accessible, consistent UI
- Radix UI primitives under the hood for a11y

We prioritize Tailwind and shadcn/ui across the app.

---

## 🤝 Contributing
We welcome issues and PRs! Please follow these guidelines to keep things smooth:

- Fork the repo and create a feature branch from `main` (e.g., `feat/<topic>` or `fix/<topic>`)
- Keep PRs focused and small; add screenshots/gifs for UI changes
- Write/adjust tests where applicable
- Run checks locally before pushing:
  - Web: `cd web && npm run lint && npm run build`
  - Backend: `cd backend && go test ./...`
  - Mobile (optional): `cd mobile && flutter analyze`
- Commit messages: conventional-ish style is appreciated (e.g., `feat:`, `fix:`, `chore:`)
- Describe your change clearly in the PR and link related issues

First time? Check open issues labeled `good first issue` or propose an idea in a new issue.

---

## 📜 License
No license specified yet. If you plan to use or contribute, consider adding a LICENSE file (e.g., MIT).

---

## 🙌 Acknowledgments
- Tailwind CSS, shadcn/ui, and Radix UI for a delightful developer experience
- Vite team for a fantastic toolchain
- Kinde for simple, secure authentication
- TanStack for React Query

---

## 🗺️ Roadmap
Near-term
- Analytics and budgeting dashboard
- Calendar view for upcoming renewals
- Import/Export (CSV)
- Notifications and reminders

Mid-term
- Mobile apps (Flutter) public beta and feature parity with web
- Cloud (Hosted) beta
- OAuth providers and improved auth flows

Long-term
- Provider integrations and webhooks
- Public API and SDK
- Localization (i18n)
- Multi-tenant features and team management

Follow progress in GitHub Issues and Milestones.

Star the repo if you like it—let’s make subscription tracking painless! 🌟
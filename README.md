# SubTracker — Track Your Subscriptions like a Pro 💸📅

A fun, modern subscription manager that helps you see where your money goes, 
avoid surprise renewals, and get the most out of your services.

Frontend: React + TypeScript + Vite • Styled with Tailwind CSS and shadcn/ui
Backend: Developed in Go (Golang)
Mobile: Coming soon (Flutter)

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
VITE_BACKEND_URL=https://api.your-subtracker.tld

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

---

## 🗂️ Project Structure (Frontend)

```
web/
  src/
    api/                  # API clients (Kiota), endpoints, and models
    components/           # UI and feature components (shadcn/ui-based)
    hooks/                # React hooks (queries, mutations, helpers)
    layouts/              # Page layouts and shells
    pages/                # Route-level pages
    lib/ utils/           # Helpers and shared libs
    assets/               # Static assets
  public/                 # Static public files
  package.json            # Scripts and dependencies
```

---

## 🔐 Backend (Go)
- The backend powering SubTracker is written in Golang.
- It exposes REST endpoints consumed by the React app.
- Repository/URL: coming soon (or replace with your backend repo link).

---

## 🧩 Design System & Styling
- Tailwind CSS utilities for layout/spacing/typography
- shadcn/ui components for accessible, consistent UI
- Radix UI primitives under the hood for a11y

We prioritize Tailwind and shadcn/ui across the app.

---

## 🤝 Contributing
Contributions are welcome! Here’s how to help:
- Fork the repo and create a feature branch
- Keep changes focused and documented
- Run linting before committing: `npm run lint`
- Open a PR with a clear description and screenshots (if UI changes)

Ideas, bug reports, and feature requests are appreciated—please open an issue.

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

## 📣 What’s next?
- More insightful analytics and budgeting tools
- Calendar view for upcoming renewals
- Mobile-friendly enhancements

Star the repo if you like it—let’s make subscription tracking painless! 🌟
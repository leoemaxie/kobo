# Triumph Academy - Kobo Reference Application

This is the "School Fees" reference implementation for Kobo, demonstrating how a third-party integrator can use the Kobo API to build financial infrastructure into their own product. 

This application simulates a school management system (Triumph Academy) where:
1. **Admins** can register students (provisioning a Kobo Identity and Virtual Account) and close accounts.
2. **Parents** can log in, view their linked students' account statements, and see live transaction history powered by Kobo.

> **Note:** This is an isolated, standalone SvelteKit application. It maintains its own local database (`parents`, `students`) and only communicates with Kobo over the public HTTP API, exactly as an external developer would.

## 🚀 Setup & Installation

### Prerequisites
- Node.js & pnpm
- A running PostgreSQL instance
- Kobo Core running locally (or API access to a hosted instance)

### 1. Install Dependencies
```bash
pnpm install
```

### 2. Environment Variables
Copy the `.env.example` file to `.env`:
```bash
cp .env.example .env
```
Fill in the credentials:
- `DATABASE_URL`: Your local Postgres connection string (e.g., `postgres://user:pass@localhost:5432/school_fees`).
- `KOBO_API_KEY`: A valid Secret Key generated from the Kobo Console.
- `KOBO_API_SECRET`: The API Secret used to sign HMAC requests.
- `KOBO_API_URL`: The base URL of the Kobo API (e.g., `http://localhost:8080/v1`).

### 3. Database Setup
We use Drizzle ORM to manage the local schema. Push the schema to your database:
```bash
pnpm drizzle-kit push
```

### 4. Run the Development Server
```bash
pnpm run dev
```
Navigate to `http://localhost:5173`. 
You can create a new Parent account by clicking "Sign up" on the login page.
> **Admin Access:** If you sign up with an email ending in `@triumph.edu`, you will automatically be granted Admin access to the school dashboard to register new students.

## 🏗 Architecture & Design
- **Framework:** SvelteKit 2 + Svelte 5
- **Styling:** Tailwind CSS v4 using Kobo's established dark-mode brand guidelines.
- **Data Layer:** Drizzle ORM + Postgres for local state (user auth, student mappings).
- **API Integration:** All Kobo API calls are strictly executed server-side (`/src/lib/server/kobo-client.ts`), ensuring credentials are never exposed to the browser. It implements standard Kobo HMAC signature requirements.
- **Deployment:** Pre-configured with `@sveltejs/adapter-vercel` for seamless Vercel deployment.

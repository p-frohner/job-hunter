# Job Hunter - React Frontend

This is the React frontend for the job hunter. It uses generated TypeScript types and React Query hooks from the OpenAPI spec to ensure type safety across the stack.

## Tech Stack

- Language: TypeScript
- Framework: React 19
- Build Tool: Vite
- Routing: TanStack Router
- Data Fetching: TanStack Query
- UI Library: Material UI (MUI)
- Codegen: Orval (OpenAPI to TypeScript/React Query hooks)

## Local Development

> If you're using Docker (`make docker-up` from the project root), you can skip this section entirely.

### Prerequisites

- Node.js (20+): [Install Node.js](https://nodejs.org/)

### Setup

1. Install dependencies:
   ```
   npm install
   ```

2. Copy the example env file and set the API URL:
   ```
   cp .env.example .env
   ```
   The default `VITE_API_URL=http://localhost:8080` works if the server is running locally.

### Run

```
make run-client
```

or equivalently:

```
npm run dev
```

The client starts on `http://localhost:5173`.

### API Code Generation

When you modify the OpenAPI spec (`openapi.yaml`), regenerate the types and React Query hooks:

```
make generate
```

or from this directory:

```
npm run api:generate
```

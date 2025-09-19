# Sykell Frontend

React + TypeScript + Vite frontend for the Sykell web crawler application.

## Technology Stack

- **React 19** - UI framework
- **TypeScript** - Type safety
- **Vite** - Fast build tool and dev server
- **Tailwind CSS** - Utility-first CSS framework
- **React Query** - Data fetching and caching
- **React Hook Form** - Form handling with validation
- **React Router** - Client-side routing
- **Zod** - Schema validation
- **Axios** - HTTP client

## Prerequisites

- **Node.js 18+** (recommended: use Node.js 20 LTS)
- **npm** or **yarn** or **pnpm**

## Quick Start

### 1. Install Dependencies

```bash
# Using npm
npm install
```

### 2. Environment Configuration

Create a `.env` file in the frontend directory:

```bash
cp .env.example .env
```

Configure the environment variables in `.env`:

```bash
# API Configuration - Backend server URL
VITE_API_BASE_URL=http://localhost:7070/api/v1
```

### 3. Start Development Server

```bash
# Using npm
npm run dev
```

The application will be available at **http://localhost:5173**

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `VITE_API_BASE_URL` | Backend API base URL | `http://localhost:7070/api/v1` | Yes |

### Notes:
- All environment variables must be prefixed with `VITE_` to be accessible in the frontend
- The backend must be running on the configured URL for the frontend to work properly
- The default configuration assumes the backend is running on `localhost:7070`

## Available Scripts

| Command | Description |
|---------|-------------|
| `npm run dev` | Start development server with hot reload |
| `npm run build` | Build for production |

## Development Workflow

### 1. Start Backend First
Ensure the backend server is running on `http://localhost:7070` before starting the frontend.

### 2. Authentication Flow
The frontend uses JWT tokens for authentication:
- **Bearer tokens** for API requests (stored in localStorage)
- **Cookies** for Server-Sent Events (SSE) connection
- Automatic token refresh and logout on authentication errors

### 3. Real-time Updates
The application uses Server-Sent Events (SSE) for real-time crawl status updates:
- Connection established automatically after login
- Live updates for crawl progress and completion
- Automatic reconnection on connection loss

## Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── AddUrlModal.tsx     # Modal for adding new URLs
│   ├── LoadingSpinner.tsx  # Loading indicator
│   ├── Pagination.tsx      # Table pagination
│   ├── ProtectedRoute.tsx  # Route protection wrapper
│   ├── SearchFilter.tsx    # Search and filter controls
│   └── UrlTable.tsx        # Main data table
├── pages/               # Page-level components
│   ├── Dashboard.tsx       # Main dashboard layout
│   ├── DashboardPage.tsx   # Dashboard page wrapper
│   ├── LoginPage.tsx       # Login form page
│   └── RegisterPage.tsx    # Registration form page
├── services/            # API communication layer
│   ├── api.ts              # Axios instance and interceptors
│   ├── auth.ts             # Authentication API calls
│   └── dashboardApi.ts     # Dashboard data API calls
├── hooks/               # Custom React hooks
│   ├── useCrawlUpdates.ts  # SSE connection for real-time updates
│   └── useDashboard.ts     # Dashboard data fetching
├── context/             # React context providers
│   └── AuthContext.tsx     # Authentication state management
├── types/               # TypeScript type definitions
│   ├── auth.ts             # Authentication types
│   ├── dashboard.ts        # Dashboard data types
│   └── validation.ts       # Form validation schemas
└── assets/              # Static assets
    └── react.svg           # React logo
```

## API Integration

The frontend communicates with the backend through:

### REST API Endpoints
- `POST /auth/register` - User registration
- `POST /auth/login` - User authentication
- `GET /auth/profile` - Get user profile
- `GET /urls` - Fetch URLs with pagination and filters
- `POST /urls` - Add new URL
- `DELETE /urls/:id` - Remove URL
- `POST /crawl/start` - Start crawling selected URLs
- `POST /crawl/stop` - Stop crawling selected URLs

### Real-time Communication
- `GET /crawl/stream` - Server-Sent Events for live updates

## Building for Production

```bash
# Build the application
npm run build
```

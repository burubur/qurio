# Frontend Failed Jobs Feature

Date: 2025-12-28

## Context
Added frontend support for managing failed ingestion jobs.

## Features
1.  **Dashboard Widget**: Added "Failed Jobs" count to the main dashboard stats grid.
2.  **Jobs View**: Dedicated page (`/jobs`) to list all failed jobs with details (error message, payload).
3.  **Retry Action**: "Retry Job" button in the Jobs View triggers the backend retry endpoint.
4.  **Navigation**: Added "Failed Jobs" link to the Sidebar.
5.  **Dev Proxy**: Configured `vite.config.ts` to proxy `/api` requests to `localhost:8081` for local development.

## Files
- `apps/frontend/src/views/JobsView.vue`: UI for listing and retrying jobs.
- `apps/frontend/src/features/jobs/job.store.ts`: State management and API calls.
- `apps/frontend/src/components/layout/Sidebar.vue`: Navigation link.
- `apps/frontend/vite.config.ts`: Proxy configuration.

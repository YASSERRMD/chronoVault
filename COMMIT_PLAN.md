# ChronoVault - Commit Plan

## Branch Structure

| Branch | Description | Status |
|--------|-------------|--------|
| phase-1-backend-db-migrations | Database schema and migrations | Pushed ✓ |
| phase-2-frontend-auth | Frontend setup and Docker config | Pushed ✓ |
| phase-3-api-endpoints | API endpoints implementation | Pending |
| phase-4-obligation-engine | Obligation evaluation engine | Pending |
| phase-5-websocket | WebSocket real-time notifications | Pending |
| phase-6-reports | Financial reports and analytics | Pending |

## Commit Plan Per Branch

### phase-1-backend-db-migrations (COMPLETED)
- feat: add database schema migrations
- docs: add SPEC.md specification

### phase-2-frontend-auth (COMPLETED)
- docs: add README with quick start guide

### phase-3-api-endpoints (NEXT)
- feat: implement auth API endpoints
- feat: implement contracts API
- feat: implement clauses API
- feat: implement obligations API
- feat: implement reports API
- feat: implement audit API

### phase-4-obligation-engine
- feat: add obligation evaluation worker
- feat: implement dependency logic
- feat: add status transition logic

### phase-5-websocket
- feat: implement WebSocket hub
- feat: add real-time notifications
- feat: connect frontend WebSocket

### phase-6-reports
- feat: financial summary endpoint
- feat: penalty tracking endpoint
- feat: risk exposure endpoint
- feat: yearly impact endpoint

## Merge Strategy

1. Each phase = separate branch
2. Create PR after each phase push
3. Squash merge to main after review
4. Delete branch after merge

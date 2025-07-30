# Football Team Management

## Postgres Setup

1. Create a Postgres database and user.
2. Create the teams, players, matches, match_results, and goals tables:

```sql
CREATE TABLE teams (
    name TEXT PRIMARY KEY,
    logo TEXT NOT NULL,
    year_founded INT NOT NULL,
    stadium_addr TEXT NOT NULL,
    city TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE players (
    name TEXT PRIMARY KEY,
    height INT NOT NULL,
    weight INT NOT NULL,
    position TEXT NOT NULL,
    jersey_number INT NOT NULL,
    team_name TEXT NOT NULL REFERENCES teams(name),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    UNIQUE(team_name, jersey_number)
);

CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    match_date DATE NOT NULL,
    match_time TIME NOT NULL,
    home_team TEXT NOT NULL REFERENCES teams(name),
    away_team TEXT NOT NULL REFERENCES teams(name),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CHECK (home_team != away_team)
);

CREATE TABLE match_results (
    id SERIAL PRIMARY KEY,
    match_id INT NOT NULL REFERENCES matches(id),
    home_score INT NOT NULL,
    away_score INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    UNIQUE(match_id)
);

CREATE TABLE goals (
    id SERIAL PRIMARY KEY,
    match_id INT NOT NULL REFERENCES matches(id),
    scorer TEXT NOT NULL,
    goal_time TEXT NOT NULL,
    team TEXT NOT NULL CHECK (team IN ('home', 'away')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);
```

3. Set the environment variables before running:

```bash
export DATABASE_URL="postgres://user:password@localhost:5432/football_db?sslmode=disable"
export JWT_SECRET="your-secret-key-change-in-production"
```

4. Run the server:

```bash
./scripts/run.sh
```

## Authentication

The API uses JWT authentication. All team, player, match, and match result management endpoints require authentication.

### Login
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testusername",
    "password": "pass"
  }'
```

### Using the Token
Include the JWT token in the Authorization header:
```bash
curl -X GET http://localhost:8080/api/v1/teams \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Features

- **JWT Authentication**: Secure API access with role-based authorization
- **Team Management**: CRUD operations for football teams
- **Player Management**: CRUD operations for players with team relationships
- **Match Schedule Management**: CRUD operations for match schedules between teams
- **Match Result Management**: CRUD operations for match results with detailed goal tracking
- **Soft Delete**: Teams, players, matches, and results are not permanently deleted but marked with a `deleted_at` timestamp
- **Timestamps**: Automatic tracking of `created_at` and `updated_at` timestamps
- **Data Integrity**: All information is preserved even after deletion
- **Business Rules**: 
  - One player can only belong to one team
  - One team can have many players
  - Jersey numbers must be unique within a team
  - Matches must be between different teams
  - Teams must exist before creating matches
  - Match results must match the number of goals scored
  - Only one result per match

## API Endpoints

### Public Endpoints
- `POST /api/v1/login` - Login and get JWT token

### Protected Endpoints (Require JWT + Admin Role)

#### Team Management
- `POST /api/v1/teams` - Register a new team
- `PUT /api/v1/teams/:name` - Update a team
- `DELETE /api/v1/teams/:name` - Soft delete a team
- `PATCH /api/v1/teams/:name/restore` - Restore a soft-deleted team

#### Player Management
- `POST /api/v1/players` - Register a new player
- `PUT /api/v1/players/:playerName` - Update a player
- `DELETE /api/v1/players/:playerName` - Soft delete a player
- `PATCH /api/v1/players/:playerName/restore` - Restore a soft-deleted player

#### Match Management
- `POST /api/v1/matches` - Register a new match schedule
- `PUT /api/v1/matches/:id` - Update a match schedule
- `DELETE /api/v1/matches/:id` - Soft delete a match schedule
- `PATCH /api/v1/matches/:id/restore` - Restore a soft-deleted match schedule

#### Match Result Management
- `POST /api/v1/match-results` - Report a match result
- `PUT /api/v1/match-results/:id` - Update a match result
- `DELETE /api/v1/match-results/:id` - Soft delete a match result
- `PATCH /api/v1/match-results/:id/restore` - Restore a soft-deleted match result

### Protected Endpoints (Require JWT Only)
- `GET /api/v1/teams` - List all active teams
- `GET /api/v1/players` - List all active players
- `GET /api/v1/players/team/:teamName` - List players by team
- `GET /api/v1/player/:playerName` - Get player by name
- `GET /api/v1/matches` - List all active matches
- `GET /api/v1/matches/team/:teamName` - List matches by team
- `GET /api/v1/match/:id` - Get match by ID
- `GET /api/v1/match-results` - List all match results
- `GET /api/v1/match-results/match/:matchID` - Get result by match ID
- `GET /api/v1/match-result/:id` - Get result by ID

## Player Positions
- `penyerang` - Forward
- `gelandang` - Midfielder  
- `bertahan` - Defender
- `penjaga gawang` - Goalkeeper

## Match Time Format
- Date: `YYYY-MM-DD` (e.g., "2024-01-15")
- Time: `HH:MM` (e.g., "19:30")

## Goal Time Format
- `MM:SS` (e.g., "45:30" for 45 minutes 30 seconds)
- `HH:MM:SS` (e.g., "01:45:30" for 1 hour 45 minutes 30 seconds)


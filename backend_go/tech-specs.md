## Backend Server Tech Specs 

## Business Goal: 
A language learning school wants to build a prototype of learning portal which will act as three things:
- Inventory of possible vocabulary that can be learned
- Act as a  Learning record store (LRS), providing correct and wrong score on practice vocabulary
- A unified launchpad to launch different learning apps

## Tech Reqs
- The backkend is GO 
- The db is SQLite 
- The API will be built using GIN & return JSON

## Directory Structure
backend_go/
├── cmd/
│   └── server/
├── internal/
│   ├── models/     # Data structures and database operations
│   ├── handlers/   # HTTP handlers organized by feature (dashboard, words, groups, etc.)
│   └── service/    # Business logic
├── db/
│   ├── migrations/
│   └── seeds/      # For initial data population
├── magefile.go
├── go.mod
└── words.db

Database: words.db (sqlite) 

words

    id: integer

    japanese: string

    romaji: string

    english: string

    parts: json

words_groups (many-to-many relationship between words and groups)

    id: integer

    word_id: integer

    group_id: integer

    groups

    id: integer

    name: string

study_sessions

    id: integer

    group_id: integer

    created_at: datetime

    study_activity_id: integer

study_activities (linking a study session to group)

    id: integer

    study_session_id: integer

    group_id: integer

    created_at: datetime

word_review_items (record of word practice)

    word_id: integer

    study_session_id: integer

    correct: boolean

    created_at: datetime


## API Endpoints 
### GET 

GET /api/dashboard/last_study_session

- Returns the most recent study session information.

GET /api/dashboard/study_progress

- Returns study progress statistics.

GET /api/dashboard/quick-stats

- Returns quick overview statistics.

GET /api/study_activities/:id

- Returns a specific study activity.

GET /api/study_activities/:id/study_sessions

- Returns study sessions for a specific study activity (paginated).

GET /api/words

- Returns words .

GET /api/words/:id

- Returns a specific word.

GET /api/groups

- Returns groups (paginated).

GET /api/groups/:id

- Returns a specific group.

GET /api/groups/:id/words

- Returns words in a specific group (paginated).

GET /api/groups/:id/study_sessions

- Returns study sessions for a specific group.

GET /api/study_sessions

- Returns study sessions.

GET /api/study_sessions/:id

- Returns a specific study session.

GET /api/study_sessions/:id/words

- Returns words in a specific study session.

### POST

POST /api/study_activities

- Creates a new study activity.

POST /api/reset_history

- Resets study history.

POST /api/full_reset

- Fully resets the system.

POSt /api/study_sessions/:id/words/:word_id/review

- Records the review of a word in a study session.

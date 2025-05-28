# Gandalf the Budget – PRD (MVP, Local‑Only, **Final Cut**)

A **single‑binary web app** to replace our family spreadsheet. Runs on one laptop, no network, no login.

---
## 1. Must‑Have Goals
| Goal | Done when… |
|------|------------|
| Spreadsheet never opened after go‑live | ✅ |
| Month‑end close ≤ 15 min | ✅ |
| JSON export works for manual backup | ✅ |

---
## 2. High‑Level Architecture
```
+--------------+          GET /  (HTML + JS)
|  net/http    |<------------------------------- Browser (localhost)
|  server.go   |          JSON API (fetch)
+-------+------+            ↑
        |                   |
        | SQLite driver     |
        v                   |
+---------------+           |
|  budget.db    |<----------+
+---------------+
```
*Static assets (`web/dist`) embedded via `embed.FS`; result is a single `gandalf-budget` executable.*

---
## 3. Repository Layout (`gandalf-budget/`)
```
cmd/
└── server/
    └── main.go              # serves dist + JSON API on :8080
internal/
├── app/                     # orchestrates use‑cases
│   ├── dashboard.go         # aggregates & summaries
│   ├── board.go             # monthly board helpers
│   └── backup.go            # JSON export
├── store/                   # SQLite helpers (sqlx)
│   ├── models.go            # structs ↔ tables
│   ├── queries.sql          # raw SQL
│   └── migrations/          # schema file 001_init.sql
└── http/
    ├── router.go            # registers routes
    ├── handlers.go          # controller funcs
    └── middleware.go        # gzip, panic recovery
web/
├── src/
│   ├── components/          # shadcn/ui wrappers
│   ├── pages/               # Dashboard.tsx, Board.tsx, Report.tsx, Manage.tsx, Backup.tsx
│   ├── hooks/               # useApi.ts, useFlash.ts
│   └── lib/                 # api.ts (fetch wrappers), format.ts
└── dist/                    # production build (git‑ignored)
README.md
```

---
## 4. Data Model (SQLite)
```sql
CREATE TABLE months (
  id INTEGER PRIMARY KEY,
  year INT NOT NULL,
  month INT NOT NULL,
  finalized BOOLEAN DEFAULT 0
);
CREATE TABLE categories (
  id INTEGER PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  color TEXT NOT NULL          -- Tailwind colour class
);
CREATE TABLE budget_lines (
  id INTEGER PRIMARY KEY,
  month_id INT NOT NULL REFERENCES months(id) ON DELETE CASCADE,
  category_id INT NOT NULL REFERENCES categories(id),
  label TEXT NOT NULL,
  expected REAL NOT NULL
);
CREATE TABLE actual_lines (
  id INTEGER PRIMARY KEY,
  budget_line_id INT NOT NULL REFERENCES budget_lines(id) ON DELETE CASCADE,
  actual REAL NOT NULL DEFAULT 0
);
CREATE TABLE annual_snaps (
  id INTEGER PRIMARY KEY,
  month_id INT NOT NULL REFERENCES months(id) UNIQUE,
  snap_json TEXT NOT NULL,
  created_at DATETIME NOT NULL
);
```

---
## 5. Use‑Cases (step‑by‑step)
### 5.1 Initial Setup (blank slate)
1. On first run the DB is empty → backend seeds the current **year + month** in `months`.
2. User opens **Manage Items**, adds categories (colour picker) and budget lines with Expected amounts.
3. As each line is added, an `actual_lines` row (0 CLP) is auto‑created.

### 5.2 Daily Edit Payment
1. Open **Monthly Board** (`/board?m=YYYY‑MM`).
2. Enter actual amount; value saved on blur via PUT `/line/:id`.
3. Backend updates row; front‑end marks row green (`actual>0`).

Validation:
* Must be ≥ 0, rounded to 2 decimals.
* Clearing value to 0 reverts row to yellow.

### 5.3 Finalize Month
1. Click **Finalize Month**.
2. App blocks if any yellow rows remain.
3. On confirm, backend transaction:
   * Save dashboard payload into `annual_snaps`.
   * Mark current month `finalized=1`.
   * Clone budget_lines into a new `months` record for next calendar month, with Expected copied and Actual reset to 0.
4. Front‑end redirects to the new (yellow) board.

### 5.4 View Annual Report
1. Choose year on **Annual Report** page.
2. List snapshots; click **View** to open read‑only dashboard render.

### 5.5 Manual Backup (export‑only)
* Click **Export JSON**; GET `/export` downloads `gandalf_YYYYMMDD.json` (pretty‑printed gzip).
* Header shows “Last backup: X days ago” using timestamp stored in localStorage.

### 5.6 Manage Categories & Budget Lines
* Add/Edit/Delete via dialogs.
* Deleting a category with attached lines requires reassignment or cascade delete confirmation.
* Colour picker limited to Tailwind palette classes (`emerald‑500`, `indigo‑600`, etc.).

---
## 6. API Endpoints (all local, JSON)
| Method | Path | Purpose |
|--------|------|---------|
| GET    | /dashboard?m=YYYY‑MM  | Aggregate payload |
| GET    | /board/:monthId       | Budget & actual lines |
| PUT    | /line/:id             | Update actual amount |
| PUT    | /finalize/:monthId    | Finalize month & clone next |
| GET    | /report?year=YYYY     | List snapshots metadata |
| GET    | /snapshot/:id         | Return stored dashboard JSON |
| GET    | /export               | JSON backup download |
| POST   | /category             | Create category |
| PUT    | /category/:id         | Update category |
| DELETE | /category/:id         | Delete/Reassign |
| POST   | /budget-line          | Create budget line |
| PUT    | /budget-line/:id      | Update expected/label |
| DELETE | /budget-line/:id      | Delete line |

All responses `application/json`; errors return `{error:"message"}` with HTTP 4xx.

---
## 7. UI & Brand System
*Colour tokens, ASCII wireframes and component map from earlier revision remain unchanged and valid.*

---
## 8. Build & Run
| Step | Command |
|------|---------|
| Install deps | `go install`, `npm i` in `web/` |
| FE build     | `npm run build` (outputs to `web/dist`) |
| Build binary | `go build -o gandalf-budget ./cmd/server` |
| Run          | `./gandalf-budget` → opens http://localhost:8080 |

---
## 9. Open Questions
_No pending questions._  If new doubts arise, add them here and agree before coding.

---
## 10. Coding Principles
This section outlines the general coding principles to be followed during the development of Gandalf the Budget. Adhering to these principles will help ensure the codebase is maintainable, understandable, and robust.

### 10.1 YAGNI (You Ain't Gonna Need It)
*   **Principle:** Implement features only when they are actively needed and part of the current requirements. Do not add functionality based on speculation that it might be useful in the future.
*   **Rationale:** This approach helps to avoid over-engineering, reduces complexity, and prevents wasted effort on features that may never be used. It keeps the codebase lean and focused on delivering value for the defined scope.

### 10.2 KISS (Keep It Simple, Stupid)
*   **Principle:** Strive for simplicity in both design and implementation. Solutions should be straightforward, easy to understand, and avoid unnecessary complexity.
*   **Rationale:** Simpler solutions are easier to develop, debug, maintain, and reason about. Complexity can often lead to more bugs and increased development time.

### 10.3 DRY (Don't Repeat Yourself)
*   **Principle:** Avoid duplication of code, logic, or data. Identify common patterns or functionalities and encapsulate them into reusable components, functions, modules, or services.
*   **Rationale:** Duplication makes the codebase harder to maintain because changes need to be made in multiple places. This increases the risk of inconsistencies and bugs. DRY code is more maintainable and less error-prone.

### 10.4 General Good Programming Practices
Beyond the core principles above, the following practices are essential for a high-quality codebase:

*   **Readability:**
    *   Write code that is clear and easy for others (and your future self) to understand.
    *   Use meaningful and consistent naming conventions for variables, functions, classes, etc.
    *   Maintain consistent formatting and indentation.
*   **Maintainability:**
    *   Structure code in a way that makes it easy to modify, debug, and extend without unintended side effects.
    *   Favor modular design where components are loosely coupled and have well-defined responsibilities.
*   **Comments:**
    *   Use comments to explain *why* something is done, or to clarify complex or non-obvious logic.
    *   Avoid over-commenting obvious code. Well-written code should be largely self-documenting.
    *   Keep comments up-to-date with code changes.
*   **Testing:**
    *   Write unit tests for individual components and functions to verify their correctness.
    *   Consider integration tests for interactions between components.
    *   Testing helps ensure code quality, prevents regressions, and provides confidence when refactoring or adding new features. (Specific testing strategies and coverage targets might be detailed in a separate document if needed).

---
# END (Original PRD)

---
---

# Progress Checklist & TODOs

## Milestone 1: Project Setup & Foundation (Complete)
- [x] Basic directory structure (`cmd/`, `internal/`, `web/`)
- [x] Go backend: Initial SQLite schema and migrations (`001_init.sql`)
- [x] Go backend: Database connection (`store.go`) and auto-seeding of current month (`setup.go`)
- [x] Go backend: Basic HTTP server setup (`cmd/server/main.go`)
- [x] Go backend: Embedding of static frontend assets (`web/dist/`)
- [x] Frontend: Vite + React + TypeScript project setup in `web/`
- [x] Frontend: Placeholder pages created (`Dashboard.tsx`, `Board.tsx`, etc.)
- [x] Build: Makefile for frontend and backend builds (`Makefile`)
- [x] Git: Initial project commit with all setup files.

## Milestone 2: Core CRUD Functionality - Categories (Complete)
- [x] Backend: Refactor HTTP routing (`internal/http/router.go`) with API versioning (`/api/v1/`)
- [x] Backend: Category model (`internal/store/models.go`)
- [x] Backend: CRUD operations for Categories in `internal/store/categories.go` (Create, ReadAll, GetByID, Update, Delete)
- [x] Backend: HTTP handlers for Category CRUD in `internal/http/category_handlers.go`
- [x] Backend: Register Category API routes (`GET /api/v1/categories`, `POST /api/v1/categories`, `PUT /api/v1/categories/:id`, `DELETE /api/v1/categories/:id`)
- [x] Frontend: `lib/api.ts` for backend communication.
- [x] Frontend: `Manage.tsx` page with UI for listing, adding, editing, and deleting categories (including form, color picker concept using Tailwind classes).
- [x] Backend: Basic tests for category handlers.

## Milestone 3: Core CRUD Functionality - Budget Lines & Actuals
- [ ] Backend: BudgetLine model (`internal/store/models.go`) - (Corresponds to `budget_lines` table)
- [ ] Backend: ActualLine model (`internal/store/models.go`) - (Corresponds to `actual_lines` table)
- [ ] Backend: Store functions for BudgetLine & ActualLine CRUD.
    - [ ] `CreateBudgetLine` should auto-create a linked `ActualLine` with 0 value (Use Case 5.1.3).
- [ ] Backend: HTTP handlers for BudgetLine & ActualLine CRUD.
    - API endpoints needed:
        - `POST /api/v1/budget-lines` (for creating a budget line, should also create its initial actual_line)
        - `PUT /api/v1/budget-lines/:id` (for updating a budget line's label or expected amount)
        - `DELETE /api/v1/budget-lines/:id` (for deleting a budget line and its associated actual_line)
        - `PUT /api/v1/actual-lines/:id` (for updating the actual amount of an actual_line - this is the `/line/:id` from PRD, but more specific) 
        - `GET /api/v1/budget-lines?month_id=:monthId` (To get all budget lines for a month, perhaps for the Manage page or board)
- [ ] Frontend: `Manage.tsx` - UI for adding/editing/deleting Budget Lines (associated with categories and a month).
- [ ] Frontend: `Board.tsx` - UI for displaying budget lines and entering/updating Actual amounts (Use Case 5.2). This will use `PUT /api/v1/actual-lines/:id`.

## Milestone 4: Monthly Board & Finalization Logic
- [ ] Backend: `GET /api/v1/board-data/:monthId` endpoint to fetch combined budget_lines and their corresponding actual_lines for a specific month.
    - This will likely involve a JOIN query.
- [ ] Frontend: `Board.tsx` - Fetch and display data from `/api/v1/board-data/:monthId`. Allow navigation to different months.
- [ ] Backend: `PUT /api/v1/months/:monthId/finalize` endpoint.
    - [ ] Logic to check if any budget lines for the month have an actual_line with amount 0 (or however "yellow rows" are defined). Prevent finalization if so.
    - [ ] Transaction:
        - [ ] Create and save dashboard payload into `annual_snaps` (requires dashboard data structure definition).
        - [ ] Mark current month `finalized=1` in `months` table.
        - [ ] Clone `budget_lines` (and their initial `actual_lines` at 0) into a new `months` record for the next calendar month.
- [ ] Frontend: `Board.tsx` - "Finalize Month" button and UI feedback. Redirect to new month on success.

## Milestone 5: Dashboard & Reporting
- [ ] Backend: Define structure for dashboard aggregate payload.
- [ ] Backend: `GET /api/v1/dashboard?month_id=:monthId` endpoint - Aggregate payload for the dashboard.
- [ ] Frontend: `Dashboard.tsx` - Display aggregated dashboard data.
- [ ] Backend: `GET /api/v1/reports/annual?year=YYYY` endpoint - List `annual_snaps` metadata for a given year.
- [ ] Backend: `GET /api/v1/reports/snapshots/:snapId` endpoint - Return stored dashboard JSON from `annual_snaps` table.
- [ ] Frontend: `Report.tsx` - UI to choose year, list snapshots, and view a read-only render of a selected snapshot.

## Milestone 6: Backup & Miscellaneous
- [ ] Backend: `GET /api/v1/export/json` endpoint - JSON backup download (pretty-printed gzip).
    - Consider what data to include (all tables or a specific selection).
- [ ] Frontend: `Backup.tsx` - "Export JSON" button.
- [ ] Frontend: Display "Last backup: X days ago" using localStorage (after successful export).
- [ ] Validation:
    - [ ] Actual amounts must be ≥ 0, rounded to 2 decimals (backend validation).
    - [ ] Deleting a category with attached budget lines: implement reassign or cascade delete confirmation (currently simple delete).
- [ ] UI/Brand System: Review and ensure adherence to color tokens, wireframes (as per PRD).

```
This detailed checklist should help track progress.

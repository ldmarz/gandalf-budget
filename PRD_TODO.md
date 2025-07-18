# Gandalf the Budget – PRD (MVP, Local‑Only, **Final Cut**)

A **single‑binary web app** to replace our family spreadsheet. Runs on one laptop, no network, no login.

---
## 1. Must‑Have Goals
| Goal | Done when… |
|------|------------|
| Spreadsheet never opened after go‑live | ✅ |
| Month‑end close ≤ 15 min | ✅ |
| JSON export works for manual backup | ❌ |

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

### 10.5 Self-Explanatory Code
*   **Principle:** Code should be self-explanatory. If not, there is a design problem. Avoid adding comments.
*   **Rationale:** Code that requires extensive comments to be understood can indicate overly complex logic, poor naming, or a design that is difficult to follow. Striving for self-documenting code through clear naming, logical structure, and well-defined components makes the codebase easier to understand and maintain directly from the code itself. While comments have their place for explaining *why* something non-obvious is done, the primary goal should be clarity in the code itself.

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
- [x] Backend: BudgetLine model (`internal/store/models.go`) - (Corresponds to `budget_lines` table)
- [x] Backend: ActualLine model (`internal/store/models.go`) - (Corresponds to `actual_lines` table)
- [x] Backend: Store functions for BudgetLine & ActualLine CRUD.
    - [x] `CreateBudgetLine` should auto-create a linked `ActualLine` with 0 value (Use Case 5.1.3).
- [x] Backend: HTTP handlers for BudgetLine & ActualLine CRUD.
    - API endpoints needed:
        - [x] `POST /api/v1/budget-lines` (for creating a budget line, should also create its initial actual_line)
        - [x] `PUT /api/v1/budget-lines/:id` (for updating a budget line's label or expected amount)
        - [x] `DELETE /api/v1/budget-lines/:id` (for deleting a budget line and its associated actual_line)
        - [x] `PUT /api/v1/actual-lines/:id` (for updating the actual amount of an actual_line - this is the `/line/:id` from PRD, but more specific) 
        - [x] `GET /api/v1/budget-lines?month_id=:monthId` (To get all budget lines for a month, perhaps for the Manage page or board)
- [x] Frontend: `Manage.tsx` - UI for adding/editing/deleting Budget Lines (associated with categories and a month).
- [x] Frontend: `Board.tsx` - UI for displaying budget lines and entering/updating Actual amounts (Use Case 5.2). This will use `PUT /api/v1/actual-lines/:id`.

## Milestone 4: Monthly Board & Finalization Logic
- [x] Backend: `GET /api/v1/board-data/:monthId` endpoint to fetch combined budget_lines and their corresponding actual_lines for a specific month.
    - This will likely involve a JOIN query.
- [x] Frontend: `Board.tsx` - Fetch and display data from `/api/v1/board-data/:monthId`. Allow navigation to different months.
- [x] Backend: `PUT /api/v1/months/:monthId/finalize` endpoint.
    - [x] Logic to check if any budget lines for the month have an actual_line with amount 0 (or however "yellow rows" are defined). Prevent finalization if so.
    - [x] Transaction:
        - [x] Create and save dashboard payload into `annual_snaps` (requires dashboard data structure definition).
        - [x] Mark current month `finalized=1` in `months` table.
        - [x] Clone `budget_lines` (and their initial `actual_lines` at 0) into a new `months` record for the next calendar month.
- [x] Frontend: `Board.tsx` - "Finalize Month" button and UI feedback. Redirect to new month on success.

## Milestone 5: Dashboard & Reporting
- [x] Backend: Define structure for dashboard aggregate payload.
- [x] Backend: `GET /api/v1/dashboard?month_id=:monthId` endpoint - Aggregate payload for the dashboard.
- [x] Frontend: `Dashboard.tsx` - Display aggregated dashboard data.
- [x] Backend: `GET /api/v1/reports/annual?year=YYYY` endpoint - List `annual_snaps` metadata for a given year.
- [x] Backend: `GET /api/v1/reports/snapshots/:snapId` endpoint - Return stored dashboard JSON from `annual_snaps` table.
- [x] Frontend: `Report.tsx` - UI to choose year, list snapshots, and view a read-only render of a selected snapshot.

## Milestone 6: Backup & Miscellaneous
- [ ] Backend: `GET /api/v1/export/json` endpoint - JSON backup download (pretty-printed gzip).
    - Consider what data to include (all tables or a specific selection).
- [x] Frontend: `Backup.tsx` - "Export JSON" button.
- [x] Frontend: Display "Last backup: X days ago" using localStorage (after successful export).
- [ ] Validation:
    - [x] Actual amounts must be ≥ 0, rounded to 2 decimals (backend validation).
    - [ ] Deleting a category with attached budget lines: implement reassign or cascade delete confirmation (currently simple delete).
- [ ] UI/Brand System: Review and ensure adherence to color tokens, wireframes (as per PRD).

## Milestone 7: Core Frontend Navigation Setup
- Goal: Implement a functional navigation system and basic routing.
- Tasks:
    - [x] Verify `react-router-dom` installation (currently `^7.6.1`) in `web/package.json`.
    - [x] Create a primary navigation component (e.g., `Navbar.tsx` in `web/src/components/layout/`) displaying links for: Dashboard, Board, Report, Manage, and Backup.
    - [x] Modify `web/src/App.tsx` to include the `Navbar` and set up `react-router-dom` `Routes` (importing `Routes`, `Route`, `Link`) for each page component currently in `web/src/pages/` (Dashboard, Board, Report, Manage, Backup).
    - [x] Ensure basic navigation allows switching between these existing pages.
    - [x] Style the `Navbar` minimally using Tailwind CSS for initial usability.

## Milestone 8: Tailwind CSS Integration and Basic Styling Pass
- Goal: Ensure Tailwind CSS is correctly configured and apply foundational styling.
- Tasks:
    - [ ] Verify Tailwind CSS configuration is active (check `tailwind.config.js`, `postcss.config.js`, and main CSS import in `web/src/index.css` or `web/src/main.tsx`).
    - [ ] Perform an initial styling pass on all existing pages (`Dashboard.tsx`, `Board.tsx`, `Report.tsx`, `Manage.tsx`, `Backup.tsx`) using Tailwind utility classes to improve readability and layout (typography, spacing, containers).
    - [ ] Ensure UI elements like buttons and inputs have consistent basic styling with Tailwind.

## Milestone 9: Shadcn UI Component Adoption
- Goal: Enhance UI consistency and interactivity by adopting Shadcn UI components across all pages.
- General Tasks:
    - [ ] Verify Shadcn UI is installed and `web/src/components/ui/` contains a comprehensive set of base components.
    - [ ] For each page, identify standard HTML elements or basic components that can be replaced with more robust Shadcn UI equivalents.
    - [ ] If a suitable Shadcn component does not exist in `web/src/components/ui/` for a required element, implement a new reusable component adhering to Shadcn principles and styled with Tailwind CSS. Place new general-purpose UI components in `web/src/components/ui/` or page-specific components in `web/src/components/web/` (or a relevant page-specific component folder).
    - [ ] Style Shadcn components using their props and Tailwind utility classes as needed to fit the application's design.
- Page-Specific Tasks:
    - [ ] **Dashboard.tsx (`/src/pages/Dashboard.tsx`)**:
        - [ ] Review and refactor using Shadcn components (e.g., `Card` for sections, `Table` for data display if any, `Button` for actions).
    - [ ] **Board.tsx (`/src/pages/Board.tsx`)**:
        - [ ] Review and refactor using Shadcn components (e.g., `Card` for overall layout, `Table` for budget lines, `Input` for amounts, `Button` for actions like "Finalize Month").
    - [ ] **Report.tsx (`/src/pages/Report.tsx`)**:
        - [ ] Review and refactor using Shadcn components (e.g., `Select` for year choice, `Table` for snapshot list, `Button` for "View").
    - [ ] **Manage.tsx (`/src/pages/Manage.tsx`)**:
        - [ ] Review and refactor using Shadcn components (e.g., `Card` for sections, `Table` for categories/budget lines, `Dialog` or `Modal` for add/edit forms, `Input`, `Button`, `Select` for form elements).
    - [ ] **Backup.tsx (`/src/pages/Backup.tsx`)**:
        - [ ] Review and refactor using Shadcn components (e.g., `Button` for "Export JSON", `Alert` or text display for "Last backup").

## Milestone 10: Frontend Interaction and UX Refinements
- Goal: Improve user experience with better visual feedback and interactive elements.
- Tasks:
    - [ ] Add loading states (e.g., using `LoadingSpinner` from `web/src/components/ui/LoadingSpinner.tsx` or a similar Shadcn component) for data fetching operations in pages like `Board.tsx` or `Dashboard.tsx`.
    - [ ] Implement user feedback messages (e.g., success/error notifications using Shadcn `Alert` or by creating a toast-like system) for actions like saving data or encountering errors.
    - [ ] Review application for areas where UI can be made more intuitive (e.g., clearer calls to action, better visual hierarchy, consistent placement of controls).

## Milestone 11: Frontend Testing for Navigation and Key UI Components
- Goal: Ensure reliability of navigation and critical UI elements.
- Tasks:
    - [ ] Write basic integration tests for the navigation system (e.g., using Vitest or React Testing Library to ensure clicking links in `Navbar.tsx` loads the correct page component).
    - [ ] Write unit tests for any new complex custom components or critical Shadcn component implementations that encapsulate specific application logic.
```
This detailed checklist should help track progress.

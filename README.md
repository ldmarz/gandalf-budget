# Gandalf the Budget: Your Guide to Financial Clarity 🧙‍♂️

Are your finances feeling like a dragon's hoard after a goblin raid – chaotic, confusing, and slightly alarming? Do spreadsheets, once trusted maps, now resemble the Mines of Moria, vast and easy to get lost in? Fear not, for a guiding light has appeared! Gandalf the Budget is here to help you navigate the perilous paths of personal finance and bring order to your treasury.

This application is a **single-binary web app** designed to replace your family spreadsheet. It runs locally on your machine, requires no network connection, and has no login system – your financial counsel is kept private and secure.

## The Quest: Wielding Gandalf's Tools 🛠️

Every great quest requires the right tools. Gandalf the Budget equips you with artifacts of power to manage your coin:

*   **📜 The Scroll of Categories:** Define and manage your realms of spending. Create, edit, and delete spending categories, each with a unique name and a distinct color, much like the banners of different lands.
*   **📒 The Ledger of Lines:** Chart your expected expenditures for each moon (month). For a chosen month, scribe your budget line items with clear labels and their foretold costs, linking each to its rightful category.
*   **🔮 The Crystal of Actuals:** Peer into the crystal to see where your gold truly flowed. Record your actual spending for each budget line upon the great Monthly Board, and watch as the crystal reveals your progress.
*   **📦 The Strongbox of SQLite:** All your financial records are kept safe and sound within a local Strongbox, powered by SQLite. No prying eyes, no network journeys needed.
*   **🖥️ The Seeing-Glass UI:** Interact with your budget through a clear and intuitive Seeing-Glass (a web-based User Interface).
*   **🌬️ The Whispering Winds API:** The Seeing-Glass communicates with the Strongbox via the swift and silent Whispering Winds of a local JSON API.

## The Source of Power: Architecture ✨

What magic fuels Gandalf the Budget? Herein lies the simple yet potent alchemy:

*   **The Go Heart (Backend):** A robust backend forged in the Go language, serving both the Seeing-Glass UI and the Whispering Winds API. It embeds all necessary enchantments (static assets) to operate as a single executable.
*   **The React Visage (Frontend):** The friendly face of your budget, crafted with React, Vite, and TypeScript, providing a dynamic and responsive experience.
*   **The SQLite Foundation (Database):** The very bedrock where your financial runes are stored securely.

The Flow of Magic:
```
+--------------+          GET / (HTML + JS)
|  Go Backend  |<-------------------------- Browser (localhost)
| (net/http)   |          JSON API (fetch)
+-------+------+
        | SQLite Driver
        v
+---------------+
|  budget.db    |
+---------------+
```

## The Ancient Runes: Data Model 📚

The deep knowledge of your finances is inscribed in ancient runes within the Strongbox of SQLite:

*   `runes_of_months`: Chronicles of the fiscal periods (your monthly adventures).
*   `sigils_of_categories`: Marks of your spending domains.
*   `scripts_of_budget_lines`: The written plans for your coin.
*   `glyphs_of_actuals`: The true tale of your expenditures.
*   `tomes_of_annual_snaps`: (Future chronicles) For the grand yearly sagas and reports.

## Summoning Your Guide: Build & Run 🚀

To summon Gandalf the Budget to your aid, follow these sacred incantations:

1.  **Gather the Essences (Install Dependencies):**
    *   For the Go Heart: `go install` (or `go mod tidy`) in the project root.
    *   For the React Visage: `npm install` in the `web/` directory.

2.  **Weave the Frontend Spell (Build Frontend):**
    *   Journey into the `web/` directory: `cd web`
    *   Chant the build command: `npm run build`
    *   This conjures static assets into `web/dist/`.

3.  **Forge the Master Artifact (Build Backend Binary):**
    *   Return to the project root: `cd ..`
    *   Speak the words of power: `go build -o gandalf-budget ./cmd/server`
    *   A single, potent executable named `gandalf-budget` shall appear!

4.  **Awaken the Guide (Run the Application):**
    *   Unleash its power: `./gandalf-budget`
    *   Open your preferred Seeing-Glass (web browser) and navigate to the mystical portal: `http://localhost:8080`.

## The Council's Wisdom: Coding Principles ⚖️

Even the wisest of wizards follows guiding principles. Ours are etched in the `PRD_TODO.md` file (Section 10) for all apprentices to study:

*   **YAGNI (You Ain't Gonna Need It)**
*   **KISS (Keep It Simple, Stupid)**
*   **DRY (Don't Repeat Yourself)**
*   **General Good Programming Practices** (Readability, Maintainability, Comments, Testing)

## A Scribe's Note ✍️

The chronicles of this project note that its initial setup, foundational structure, and the development of key features (including Milestones 1-3 and this very README) were significantly assisted by Jules, a large language model from Google.
```

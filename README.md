# todometrics

<p align="center">
  <img src="https://raw.githubusercontent.com/rinzler69-wastaken/todometrics-pttm/main/server/static/images/appicon-dark.svg" alt="todometrics App Icon" width="120">
</p>

<p align="center">
  <strong>Not your average To-Do list.</strong>
  <br />
  A personal task management system with performance metrics, built with Go.
</p>

<p align="center">
  <a href="https://goreportcard.com/report/github.com/rinzler69-wastaken/todometrics-pttm"><img src="https://goreportcard.com/badge/github.com/rinzler69-wastaken/todometrics-pttm" alt="Go Report Card"></a>
  <a href="https://www.gnu.org/licenses/gpl-3.0"><img src="https://img.shields.io/badge/License-GPLv3-blue.svg" alt="License: GPL v3"></a>
</p>

![todometrics Screenshot](<link-to-your-screenshot.png>)
*Replace the link above with a full screenshot of your application.*

---

## 🧠 Philosophy

**todometrics** is designed for simplicity and discipline. The core principles are **immutable commitments** and **measurable progress**.

* **Deliberate Planning:** Core task details (Name, Due Date, Priority) are set in stone upon creation.
* **Flexible Method:** The commitment is firm, but the method can evolve. Only the 'Description' is editable.
* **Decisive Action:** Completing a task is an irreversible action that moves it to your permanent archive.

---

## ✨ Features

* **Dynamic Overview Dashboard:** Get at-a-glance insights into your performance with statistics for ongoing and completed tasks.
* **Interactive Performance Graph:** Track your discipline over time with a simple scoring system (+2 for ahead of time, +1 for on time, -1 for overdue), visualized in a monthly or yearly graph.
* **Advanced Task Management:** Create, edit, complete, and delete tasks through a clean and responsive modal interface.
* **Completed Task Archive:** An immutable archive of all your accomplishments, with filtering by month/year and client-side search.
* **PDF Export:** Export your filtered completed tasks to a clean PDF report for your records.
* **Persistent Light/Dark Theme:** A sleek, modern UI that respects your system settings and can be toggled manually. The choice is stored locally to prevent FOUC (Flash of Unstyled Content).
* **Client-Side Search:** Instantly find ongoing or completed tasks with a type-to-search function.
* **Responsive Design:** A great user experience on both desktop and mobile, featuring a dedicated bottom navigation bar for smaller screens.

---

## 🛠️ Technology Stack

* **Backend:**
    * **Go:** The core programming language.
    * **Gin Gonic:** A high-performance HTTP web framework for routing.
    * **GORM:** A developer-friendly ORM for database interactions.
    * **SQLite:** The self-contained SQL database engine.
* **Frontend (Current):**
    * **Go `html/template`:** For server-side rendering of the UI.
    * **Vanilla JavaScript (ES6):** For all client-side interactivity, structured into reusable modules.
    * **Bootstrap 5:** For the responsive grid system and core components.
    * **Chart.js:** For rendering the interactive performance graph.
    * **Flatpickr:** For a lightweight and powerful datetime picker.
    * **Custom CSS3:** For advanced styling, theming, and animations.

---

## 🚀 Getting Started

Follow these instructions to get a copy of the project up and running on your local machine.

### Prerequisites

You need to have Go installed on your machine (version 1.18 or newer).

### Installation

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/rinzler69-wastaken/todometrics-pttm.git](https://github.com/rinzler69-wastaken/todometrics-pttm.git)
    cd todometrics-pttm
    ```

2.  **Navigate to the server directory:**
    The Go application lives in the `server` subdirectory.
    ```bash
    cd server
    ```

3.  **Install dependencies:**
    Go will automatically fetch the required modules from `go.mod`.
    ```bash
    go mod tidy
    ```

4.  **Run the application:**
    ```bash
    go run main.go
    ```
    The server will start, and you can access the application at `http://localhost:8080`.

---

## 🗺️ Roadmap

This project is currently a **server-side rendered application**. The next major architectural step is to evolve it into a modern, decoupled web application.

* **Phase 1: API Refactor (In Progress)**
    * Convert all Go handlers to return pure JSON data, transforming the backend into a stateless JSON API.

* **Phase 2: SPA Frontend**
    * Develop a new Single-Page Application (SPA) frontend using **SvelteKit**.
    * This new client will live in the `/client` directory and will consume the Go API for a faster, more fluid user experience.

---

## 📄 License

This project is licensed under the **GNU General Public License v3.0**. See the [LICENSE](https://www.gnu.org/licenses/gpl-3.0.html) file for details.
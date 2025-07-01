
# TrackMyBugs — Full-Stack Issue Tracker for Modern Teams

**TrackMyBugs** is a production-grade, full-stack bug and issue tracking system built with **Go**, **Next.js**, **PostgreSQL**, and **Docker**. Designed for teams and individuals who value performance, security, and a clean, intuitive user experience.

---

## Why TrackMyBugs?

✅ **Robust Architecture** — Clean separation of concerns with dedicated backend, frontend, and database services.  
✅ **Scalable & Secure** — Go-powered REST API with JWT authentication, RBAC, and PostgreSQL at the core.  
✅ **Real-World Functionality** — Issue tracking, project management, advanced filtering, and search.  
✅ **Developer-Focused** — Built with modern tools, Dockerized for easy deployment, designed for extensibility.  

---

## Tech Stack

| Layer             | Technology                          |
|-------------------|-------------------------------------|
| **Frontend**      | Next.js 14 (React), Tailwind CSS    |
| **Backend**       | Go (Gin Framework), JWT Auth, RBAC |
| **Database**      | PostgreSQL                         |
| **Infrastructure**| Docker, Docker Compose             |

---

## Core Features

- **User Authentication:** Secure JWT-based sessions (login/register)  
- **Role-Based Access Control:** Admin/user roles, admin-only management actions  
- **Project Management:** Create, edit, delete, and search projects  
- **Issue Tracking:** Full CRUD for issues with assignment, filtering, and prioritization  
- **Comment System:** Discuss issues with threaded comments, edit/delete support  
- **Advanced Filtering & Search:** Filter issues by status, priority, assignee, and more  
- **Pagination:** Optimized queries for large datasets (projects, issues, comments)  
- **Responsive UI:** Clean, accessible design with Tailwind CSS  
- **Containerized Deployment:** Docker-powered setup for consistent environments

---

## Quick Start

### Prerequisites

- [Node.js](https://nodejs.org/) 18+  
- [Go](https://golang.org/) 1.21+  
- [Docker](https://www.docker.com/) & Docker Compose  

---

### Run with Docker (Recommended)

```bash
git clone https://github.com/yourusername/trackmybugs.git
cd trackmybugs
docker-compose up --build
```
-   Frontend: [http://localhost:3000](http://localhost:3000/)
    
-   Backend API: [http://localhost:8080](http://localhost:8080/)
    
-   Database: `localhost:5432`
    

----------

### Manual Setup

1.  **Database**  
    Run the SQL in `db/init.sql` on your PostgreSQL instance.
    
2.  **Backend**
    
    ```bash
    cd backend
    go run .
    ```
    
3.  **Frontend**
    
    ```bash
    cd frontend
    npm install
    npm run dev
    ```
    

----------

## Project Structure

```
trackmybugs/
├── frontend/         # Next.js frontend
├── backend/          # Go REST API
├── db/               # SQL schema and migrations
├── docker-compose.yml
└── README.md

```
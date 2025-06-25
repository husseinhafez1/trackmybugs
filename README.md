# TrackMyBugs 🐛

A full-stack issue tracking system built with Next.js frontend, Go REST API, and PostgreSQL.

## 🎯 Project Overview

TrackMyBugs is a simple but powerful bug tracking application that helps teams manage and track issues efficiently. Built with modern technologies and best practices.

## 🏗️ Tech Stack

- **Frontend**: Next.js 14 (React)
- **Backend**: Go (Gin framework)
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Styling**: Tailwind CSS
- **Deployment**: Docker Compose for local development

## 📁 Project Structure

```
trackmybugs/
├── frontend/          # Next.js application
├── backend/           # Go REST API
├── db/               # Database scripts and migrations
├── docker-compose.yml # Local development setup
└── README.md         # This file
```

## 🚀 Getting Started

### Prerequisites

- Node.js 18+
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (or use Docker)

### Quick Start

1. Clone the repository
2. Run `docker-compose up` for local development
3. Access the application at `http://localhost:3000`

## 📋 Development Roadmap

### Phase 1 - Setup ✅
- [x] Project structure
- [x] Docker Compose configuration
- [x] Basic README

### Phase 2 - Database & Backend
- [x] Database schema design
- [x] Go API setup with Gin
- [x] Basic CRUD operations
- [x] JWT authentication
- [ ] Role-based permissions

### Phase 3 - Frontend
- [ ] Next.js project setup
- [ ] Authentication pages
- [ ] Dashboard and project views
- [ ] Issue management interface

### Phase 4 - Core Features
- [ ] Issue status transitions
- [ ] Comment system
- [ ] User assignment
- [ ] Filtering and sorting

### Phase 5 - Polish & Extras
- [ ] UI improvements
- [ ] Error handling
- [ ] Pagination
- [ ] Deployment setup

### Phase 6 - Bonus Features
- [ ] Real-time updates
- [ ] Activity logging
- [ ] Testing suite
- [ ] CI/CD pipeline

## 🤝 Contributing

This is a learning project. Feel free to explore the code and learn from it!

## 📄 License

MIT License - feel free to use this code for your own projects. 
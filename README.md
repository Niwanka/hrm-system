# Next-Gen HRM System - Core Auth & Employee Directory

A modern Human Resource Management (HRM) system built with a **Go (Fiber + GORM)** backend and a **Vue 3 (Vite + Pinia + Vue Router + Tailwind CSS)** frontend with dual Light & Dark theme support in a vibrant Emerald Green aesthetic.

---

## 🚀 Tech Stack

- **Backend:** Go (Golang) with [Fiber v2](https://gofiber.io/) web framework
- **Database & ORM:** PostgreSQL using [GORM](https://gorm.io/) (with auto-table migrations and SQLite fallback)
- **Security:** HTTP-Only JWT authentication & bcrypt password hashing
- **Frontend:** Vue 3 (Composition API / `<script setup>`) built with Vite
- **State Management:** Pinia (`authStore`, `themeStore`)
- **Routing:** Vue Router with RBAC navigation guards
- **Styling:** Tailwind CSS with Emerald Green palette (`#10b981`) & Light/Dark mode support

---

## 🔑 Demo Accounts & Pre-Seeded Roles

The system auto-migrates and seeds sample users into the database upon launch:

| Role | Work Email | Password | Access Level | Permissions |
|---|---|---|---|---|
| **Admin** | `admin@company.com` | `password123` | **100** | Full System Access, Payroll, Directory & Employee Insertion |
| **HR** | `hr@company.com` | `password123` | **80** | Company Payroll, Directory & Employee Insertion |
| **Manager** | `manager@company.com` | `password123` | **50** | Direct Reports, Directory & Employee Insertion |
| **Employee** | `employee@company.com` | `password123` | **10** | Employee Directory & Org Hierarchy View |

---

## ⚙️ Quick Start Guide

### 1. Start Go Backend
```bash
cd backend
go run main.go
```
*Backend runs on `http://localhost:8080` (or dynamic `PORT` env var).*

### 2. Start Vue 3 Frontend
```bash
cd frontend
npm install
npm run dev
```
*Client runs on `http://localhost:5173`.*

---

## 📦 Project Structure

```
hrm-system/
├── backend/
│   ├── main.go
│   ├── config/          # GORM DB connection & auto-seeder
│   ├── controllers/     # Auth & Employee CRUD controllers
│   ├── middleware/      # Fiber JWT auth & RBAC enforcement
│   ├── models/          # Role & Employee GORM models
│   └── routes/          # API route definitions
└── frontend/
    ├── src/
    │   ├── stores/      # Pinia stores (authStore, themeStore)
    │   ├── router/      # Vue Router with RBAC guards
    │   ├── components/  # Layout, Sidebar, OrgTreeNode
    │   └── views/       # Login, Dashboard, Directory, Hierarchy, Payroll
    ├── vite.config.js
    └── tailwind.config.js
```

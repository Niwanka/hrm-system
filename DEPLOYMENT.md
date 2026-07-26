# 100% FREE Production Deployment Guide

This guide walks you through deploying the **Next-Gen HRM System** for **$0/month** using 100% free production services with SSL certificates, managed database, and automated CI/CD!

---

## 🌟 The 100% Free Hosting Stack

| Layer | Provider | Free Features | Cost |
|---|---|---|---|
| **Database** | **Supabase** or **Neon.tech** | Free Managed PostgreSQL (500MB storage, connection pooler) | **$0/mo** |
| **Backend API** | **Render.com** or **Koyeb** | Free Web Service (Docker container execution, automatic SSL) | **$0/mo** |
| **Frontend CDN** | **Vercel** or **Netlify** | Free Global CDN, unlimited builds, free SSL domain (`.vercel.app`) | **$0/mo** |
| **CI/CD Pipeline**| **GitHub Actions** | 2,000 build minutes/month for automated tests & deploys | **$0/mo** |

---

## Step 1: Create Free PostgreSQL Database on Supabase

1. Sign up for a free account at **[supabase.com](https://supabase.com)** (No credit card needed).
2. Click **New Project** -> set project name to `hrm-db` and set a strong database password.
3. Once created, go to **Project Settings** -> **Database** -> **Connection String**.
4. Copy the URI string (looks like `postgres://postgres.[ref]:[password]@aws-0-us-east-1.pooler.supabase.com:6543/postgres`).

---

## Step 2: Deploy Go Fiber Backend for Free on Render.com

1. Sign up for a free account at **[render.com](https://render.com)**.
2. Click **New +** -> **Web Service** -> Connect your GitHub repository (`hrm-system`).
3. Select **Docker** environment (Render will automatically detect `backend/Dockerfile`).
4. Set **Root Directory** to `backend`.
5. Under **Environment Variables**, add:
   - `DB_HOST`: *Supabase host address*
   - `DB_USER`: `postgres`
   - `DB_PASSWORD`: *Your Supabase DB Password*
   - `DB_NAME`: `postgres`
   - `DB_PORT`: `6543` (or `5432`)
   - `PORT`: `8080`
   - `JWT_SECRET`: `hrm-super-secret-production-key-2026`
6. Click **Create Web Service**.
7. Render will build your Go container and assign a free HTTPS URL (e.g. `https://hrm-backend.onrender.com`).

---

## Step 3: Deploy Vue 3 Frontend for Free on Vercel

1. Sign up for a free account at **[vercel.com](https://vercel.com)** using your GitHub account.
2. Click **Add New...** -> **Project** -> Select your `hrm-system` repo.
3. Set **Root Directory** to `frontend`.
4. Vercel automatically detects **Vite + Vue.js**.
5. Under **Environment Variables**, add:
   - `VITE_API_BASE_URL`: `https://hrm-backend.onrender.com` (Your backend Render URL)
6. Click **Deploy**.
7. Vercel will build and publish your app to a free HTTPS domain (e.g. `https://hrm-system.vercel.app`).

---

## Step 4: Enable Automated GitHub Actions CI/CD

Every time you run `git push origin main`:

1. **GitHub Actions CI (`ci.yml`)** runs automatically:
   - Runs `go test ./...` & verifies backend compilation.
   - Runs `npm run build` & verifies Vue 3 bundle compilation.
2. **GitHub Actions CD (`cd.yml`)** triggers:
   - Render automatically pulls the latest commit and re-deploys the Go API.
   - Vercel automatically updates the live Vue 3 frontend CDN!

---

## 🎉 Congratulations!
Your full-stack enterprise HRM system is now live on the internet with:
- 🔒 HTTPS SSL encryption on all domains
- 🐘 Production-grade PostgreSQL database
- ⚡ Global CDN edge performance
- 🤖 Automated CI/CD deployment pipelines
- 💰 **Total Cost: $0.00 / month!**

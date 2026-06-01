# Contributing to goPanel 🤝

First off, thank you for considering contributing to `goPanel`! It's people like you who make the self-hosted and home lab community so amazing.

Here are the guidelines to ensure a smooth, professional, and efficient contribution process.

---

## Code of Conduct

By participating in this project, you agree to maintain a respectful, welcoming, and constructive environment. Please be professional and supportive in all issues and pull requests.

## How Can I Contribute?

### 1. Reporting Bugs 🐛
- Search the open issues list to make sure the bug hasn't already been reported.
- If it's a new issue, use our **Bug Report Template** and provide:
  - A clear, descriptive title.
  - Steps to reproduce the bug.
  - Expected behavior vs. actual behavior.
  - Telemetries or console logs (errors from backend/frontend).

### 2. Suggesting Enhancements 💡
- Create a **Feature Request** issue describing:
  - What problem the enhancement solves.
  - How you envision the feature working.
  - Any architectural changes needed.

### 3. Submitting Pull Requests (PRs) 🚀
- **Branch Naming**: Use clean, descriptive branch names:
  - `feat/some-new-feature`
  - `fix/some-annoying-bug`
  - `docs/clarify-instructions`
- **Coding Standards**:
  - Keep Go packages focused, clean, and well-commented.
  - Leverage standard libraries where possible to prevent package inflation.
  - Maintain the glassmorphic dark-mode obsidian visual theme when editing Vue/CSS files.
  - Run type checks (`vue-tsc --noEmit`) and lints before committing!
- **Commit Messages**: Write meaningful, conventional commit messages:
  - `feat: add 2FA authentications challenge`
  - `fix: resolve file zip slip path vulnerabilities`

---

## Dev Setup Guide

To set up a local development workspace:

1. **Fork and Clone** the repository.
2. **Launch backend dev server** using default port `3636`:
   ```bash
   go run main.go
   ```
3. **Launch Vue development server** (runs on `localhost:5173`):
   ```bash
   cd web
   npm install
   npm run dev
   ```
   *Note: Update the proxy target in `vite.config.ts` if your backend runs on a different port.*

4. **Verify your changes** build cleanly:
   ```bash
   # Build frontend static assets
   npm run build
   
   # Build final Go executable
   go build -o gopanel .
   ```

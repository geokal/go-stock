# Repository Guidelines

go-stock is a Wails-based (Go + Vue.js) desktop application for AI-powered stock analysis, supporting A-shares, HK, and US markets.

## Project Structure & Module Organization

```
go-stock/
├── backend/              # Go business logic
│   ├── agent/           # AI agent system (react/plan_execute modes, 9+ LLM providers)
│   ├── data/            # Stock data APIs (EastMoney, Sina, TDX, Tushare)
│   └── db/              # SQLite via GORM
├── frontend/            # Vue.js + NaiveUI + TDesign frontend
│   ├── src/components/  # 40+ Vue components (KLine charts, AI assistants, etc.)
│   ├── src/locales/    # i18n (en.json, zh-CN.json)
│   └── wailsjs/        # Auto-generated Go-to-JS bindings
├── ai-assistant-web/    # Embedded HTTP service (separate Vue frontend)
├── scripts/             # Build scripts (build-windows.sh, build-linux.sh, etc.)
└── docs/                # Chinese documentation
```

**Key architectural points:**
- Wails bridges Go backend → Vue.js frontend. Use `wailsjs/go/main/App.js` to call Go from JS.
- Agent mode supports two execution strategies: `react` (fast, single-step) and `plan_execute` (multi-step reasoning with forced tool calling).
- Multi-provider LLM factory in `backend/agent/chat_model_factory.go` handles OpenAI, DeepSeek, Ollama, LMStudio, AnythingLLM, Volcano Ark, DashScope, Qianfan, Anthropic, Gemini, OpenRouter.
- ai-assistant-web runs as a goroutine inside the main app, serving its own Vue frontend on a configurable port (`AI_ASSISTANT_WEB_ADDR`).
- Stock data comes from EastMoney (primary), Sina, TongDaXin, and Tushare — crawlers use `chromedp` for JS-rendered pages.
- VIP/sponsor verification is encrypted and built into the app.

## Build, Test, and Development Commands

**Frontend (main app):**
```bash
cd frontend
npm run dev          # Start Vite dev server (http://localhost:5173)
npm run build        # Production build
npm run preview      # Preview production build
```

**ai-assistant-web frontend:**
```bash
cd ai-assistant-web/frontend
npm run dev          # Vite dev server
npm run build        # TypeScript check + production build
```

**Full Wails app build:**
```bash
# From project root
wails build --clean  # Clean build
wails dev            # Dev mode with live reload

# Linux-specific
./scripts/build-linux.sh
```

**GitHub Actions CI** (triggers on `*-release`, `*-dev`, `*-pre` tags):
```bash
wails build           # Uses Go 1.26, Node 20.19.3
```

**Quick test** (single component): No isolated test runner configured. Manual testing via `npm run dev`.

## Coding Style & Naming Conventions

**No enforced conventions** — project has no ESLint, Prettier, EditorConfig, or pre-commit hooks.

That said, the codebase follows these patterns organically:
- Go code uses `camelCase` for functions/variables, `PascalCase` for types/structs
- Vue components use `PascalCase` for component names and `kebab-case` for file names
- Vue i18n keys use `section.keyName` format (e.g., `t('stock.dayK')`), stored in `src/locales/*.json`
- Chinese PR template in `.github/pull_request_template.md` — describe changes, link issues, note file changes

## Testing Guidelines

**No formal test suite.** The project relies on manual testing and GitHub Actions build verification (which only runs `wails build`, no unit/integration tests).

## Commit & Pull Request Guidelines

**Commit format** (conventional Chinese messages):
```
<type>(<scope>): <description>
feat(prompt): 新增提示词广场功能
fix(uplimitLadder): 修复日期选择器默认值问题
chore(deps): 升级项目依赖库版本
```

**Types:** `feat`, `fix`, `chore`, `docs`, `refactor`, `perf`, `test`

**PR template** (`.github/pull_request_template.md` — Chinese): Include description of changes, linked issues, per-file change notes, test results, and deployment notes.
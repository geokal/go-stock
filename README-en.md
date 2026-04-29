# go-stock: AI-Powered Stock Analysis Tool

## ![go-stock](./build/appicon.png)
[![GitHub Release](https://img.shields.io/github/v/release/ArvinLovegood/go-stock)](https://github.com/ArvinLovegood/go-stock/releases)
[![GitHub Repo stars](https://img.shields.io/github/stars/ArvinLovegood/go-stock)](https://github.com/ArvinLovegood/go-stock)
[![star](https://gitee.com/arvinlovegood_admin/go-stock/badge/star.svg?theme=dark)](https://gitee.com/arvinlovegood_admin/go-stock)

### 🌟 WeChat Official Account
![QR code](./build/screenshot/扫码_搜索联合传播样式-白色版.png)

### 📈 Community

- QQ Group: [Join go-stock discussion group](http://qm.qq.com/cgi-bin/qm/qr?_wv=1027&k=0YQ8qD3exahsD4YLNhzQTWe5ssstWC89&authKey=usOMMRFtIQDC%2FYcatHYapcxQbJ7PwXPHK9OypTXWzNjAq%2FRVvQu9bj2lRgb%2BSZ3p&noverify=0&group_code=491605333) (periodically cleaned up)
- SiliconFlow: Register and get 20M tokens free — [Sign up](https://cloud.siliconflow.cn/i/foufCerk)

### ✨ About

- Built with Wails and NaiveUI, powered by AI large language models for stock analysis.
- Supports A-shares, HK stocks, and US stocks. Fund/ETF support planned.
- Features: market-wide and individual stock sentiment analysis, K-line technical indicator analysis, and more.
- **Disclaimer**: This project is for entertainment purposes only. AI analysis results are for reference only. Trading involves risk — use at your own discretion.
- Primary development environment is Windows 10+. Other platforms may have limited or untested features.

### 📃 User Manual
[go-stock User Manual](docs/go-stock使用手册.md)

### 📦 Downloads
- Portable (Windows): [go-stock-windows-amd64.exe](https://github.com/ArvinLovegood/go-stock/releases)
- Portable (macOS): [go-stock-darwin-universal](https://github.com/ArvinLovegood/go-stock/releases)

### 💬 Supported AI Models / Platforms

| Model | Status | Notes |
| --- | --- | --- |
| [OpenAI](https://platform.openai.com/) | ✅ | OpenAI-compatible endpoints |
| [Ollama](https://ollama.com/) | ✅ | Local LLM platform |
| [LMStudio](https://lmstudio.ai/) | ✅ | Local LLM platform |
| [AnythingLLM](https://anythingllm.com/) | ✅ | Local knowledge base |
| [DeepSeek](https://www.deepseek.com/) | ✅ | deepseek-reasoner, deepseek-chat |
| [Aggregated Platforms](https://cloud.siliconflow.cn/i/foufCerk) | ✅ | SiliconFlow, Volcano Ark |

**If you find this project useful, please give it a star — thank you!** ⭐

- Volcano Ark: 500K tokens free for new users — [Sign up](https://www.volcengine.com/experience/ark?utm_term=202502dsinvite&ac=DSASUQY5&rc=IJSE43PZ)

This project is actively developed. Please use the latest release when possible. Issues and PRs are welcome! And of course, [sponsorship is greatly appreciated](#-sponsor-me).

### 💕 Sponsor Plans

| Plan | Tier | Benefits |
|:---|:---|:---|
| Free | vip0 | Full features, auto-update from GitHub (you handle network issues) |
| ¥18.8/month or ¥120/year | vip1 | Full features, auto-update via CDN (fast), AI config guidance, prompt templates |
| ¥28.8/month or ¥240/year | vip2 | All vip1 features, plus auto-sync last 24h market news on startup, go-stock AI Assistant |
| Custom | vipX | More plans depending on project growth... |

## 🧩 Major Feature Roadmap

| Feature | Status | Notes |
|---|---|---|
| Stock Analysis Knowledge Base | 🚧 | Planned (available in Research Center) |
| AI Stock Screening | ✅ | AI intelligent stock screening (Market → AI Summary / AI Agent) |
| ETF Support | 🚧 | ETF data (can view NAV and estimates) |
| US Stock Support | ✅ | |
| HK Stock Support | ✅ | |
| Multi-turn Conversation | ✅ | Continue asking questions after AI analysis |
| Custom AI Analysis Prompt Templates | ✅ | Configurable prompt templates [v2025.2.12.7-alpha](https://github.com/ArvinLovegood/go-stock/releases/tag/v2025.2.12.7-alpha) |
| No longer requires Chrome | ✅ | Uses Edge by default for news crawling |

## 👀 Changelog

### 2026
- **04.12** — Added AgentMode support for React and PlanExecute agent modes
- **04.11** — Added MCP tool-calling support
- **04.04** — Added holiday tool support
- **04.03** — Added stock change data (A-shares intraday)
- **03.10** — New AI Assistant feature
- **02.08** — Improved AI recommended stocks and data query features
- **01.25** — AI analysis reports and AI stock recommendation history

### 2025
- **12.16** — New AI thinking mode and hot stock strategy features
- **11.21** — Frequency-weighted sentiment analysis
- **10.30** — AI Agent feature toggle (off by default due to UX issues), removed page watermark
- **09.27** — Institutional/broker research report AI tool functions
- **08.09** — AI Agent chat feature
- **07.08** — Auto-update feature
- **07.07** — Mini intraday chart on cards
- **07.05** — macOS support
- **07.01** — AI analysis integrated with tool functions
- **06.30** — Indicator-based stock screening
- **06.27** — Financial calendar and major events timeline
- **06.25** — Hot stocks, events and topics
- **06.18** — Updated built-in stock data, real-time market news alerts, added industry research
- **06.15** — Company announcement search/view; individual stock research reports popup
- **06.13** — Individual stock research reports
- **06.12** — Dragon Tiger List, new industry ranking categories
- **05.30** — Optimized stock intraday chart display
- **05.20** — Fixed CaiLianShe telegraph retrieval issue
- **05.16** — Improved capital trend chart component
- **05.15** — Refactored app loading and data initialization, added stock capital trend, added main force net inflow data
- **05.14** — Individual stock capital flow, added K-line chart popup to rankings
- **05.13** — Industry ranking
- **05.09** — A-share order book data parsing and display
- **05.07** — Optimized intraday chart display
- **04.29** — Complete HK/US stock basic data, improved HK stock price delay, optimized initialization
- **04.25** — Market news supports AI analysis and summary: let AI read the market for you!
- **04.24** — New market module: real-time global market info — no more visiting financial websites!
- **04.22** — Improved K-line chart display with zoom support
- **04.21** — Optimized HK and US stock K-line data retrieval
- **04.01** — Optimized some settings to avoid restart
- **03.31** — Improved data crawling
- **03.30** — AI auto-scheduled analysis
- **03.29** — Multiple prompt template management, support selecting different templates for AI analysis
- **03.28** — AI analysis results saved as markdown with customizable save location
- **03.15** — Custom browser path for crawler
- **03.14** — Optimized build, significantly reduced output file size
- **03.09** — Fund NAV and estimate monitoring
- **03.06** — Community sharing feature
- **02.28** — US stock data support
- **02.23** — Danmu (bullet comments) feature for live monitoring
- **02.22** — HK stock data support (with delay)

## 🦄 Major Updates

### 2026.03.10 — New AI Assistant
![img_1.png](build/screenshot/img16.png)

### 2025.11.21 — Frequency-Weighted Sentiment Analysis
![img_1.png](build/screenshot/img15.png)

### 2025.04.25 — AI Market News Summary
![img.png](img.png)

### 2025.04.24 — Market Module
![img.png](build/screenshot/img13.png)
![img_13.png](build/screenshot/img_13.png)
![img_14.png](build/screenshot/img_14.png)

### 2025.01.17 — AI Stock Analysis
![img_5.png](build/screenshot/img.png)

## 📸 Screenshots

![img_1.png](build/screenshot/img_6.png)

### Settings
![img_12.png](build/screenshot/img_4.png)

### Cost Settings
![img.png](build/screenshot/img_7.png)

### Daily K-Line
![img_12.png](build/screenshot/img_12.png)

### Intraday Chart
![img_3.png](build/screenshot/img_9.png)

### DingTalk Alerts
![img_4.png](build/screenshot/img_5.png)

### AI Stock Analysis
![img_5.png](build/screenshot/img.png)

### Version Info
![img_11.png](build/screenshot/img_11.png)

## 💕 Acknowledgments

Built with:
- [NaiveUI](https://www.naiveui.com/)
- [Wails](https://wails.io/)
- [Vue](https://vuejs.org/)
- [Vite](https://vitejs.dev/)
- [Tushare](https://tushare.pro/register?reg=701944)

## 😘 Sponsor Me

If this project has been helpful, please consider sponsoring — thank you!

| Alipay | WeChat Pay |
|---|---|
| ![alipay.jpg](build/screenshot/alipay.jpg) | ![wxpay.jpg](build/screenshot/wxpay.jpg) |

## ⭐ Star History

[![Star History Chart](https://api.star-history.com/svg?repos=ArvinLovegood/go-stock&type=Date)](https://star-history.com/#ArvinLovegood/go-stock&Date)

## 🤖 CI Status

![Alt](https://repobeats.axiom.co/api/embed/40b07d415a42c2264a18c4fe1b6f182ff1470687.svg "Repobeats analytics image")

## 🐳 Technical Support Disclaimer

- This software is built on open-source technologies: Wails, NaiveUI, Vue, AI LLMs, etc. For technical questions, please seek help from the respective open-source communities first.
- Open-source is hard and my time is limited. For one-on-one technical support, please sponsor first. Contact QQ (note Technical Support): 506808970

| Support Type | Cost |
|:---|---:|
| QQ consultation | ¥100/time |
| Long-term technical support (unlimited, new features priority, etc.) | ¥5000 |

## License

[GNU GPLv3](LICENSE)
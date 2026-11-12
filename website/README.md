# Halo 官网 / Landing Site

Halo 的产品落地页，纯静态（HTML + CSS + 原生 JS），无构建步骤、无外部依赖。

## 本地预览

```bash
cd website
python3 -m http.server 8788
# 打开 http://127.0.0.1:8788/
```

## 部署到 GitHub Pages

1. 仓库 `Settings → Pages`
2. Source 选择 `Deploy from a branch`
3. Branch 选择 `main`，目录选择 `/website`
4. 保存后即通过 `https://<user>.github.io/halo/` 访问

也可直接使用任意静态托管（Netlify / Vercel / Cloudflare Pages），发布目录指向 `website/`。

## 结构

```
website/
├── index.html      # 页面结构与中英双语文案
├── styles.css      # iOS 设计系统（含深色模式、响应式）
├── main.js         # 语言切换、复制、Tab、滚动动画
└── assets/         # 产品截图（graph/events/approval/detail）
```

设计语言：SF 字体栈、系统色（#007AFF 等）、毛玻璃、圆角卡片、`prefers-color-scheme` 自适应。

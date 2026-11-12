# 接入 Cline

Cline（VS Code 扩展）支持「OpenAI 兼容」自定义端点，可直接指向 Halo。

## 1. 配置上游

在 `halo.yaml` 中设置 OpenAI 兼容上游：

```yaml
mode: observe
upstream:
  baseUrl: "https://api.openai.com"
```

启动 Halo：

```sh
halo
```

## 2. 在 Cline 中配置

1. 打开 Cline 面板，点击顶部设置图标。
2. `API Provider` 选择 **OpenAI Compatible**。
3. `Base URL` 填入：
   ```
   http://127.0.0.1:8787
   ```
4. `Model ID` 填入上游支持的模型名（如 `gpt-4o`）。
5. `API Key` 填入你的真实 Key，由 Halo 转发到上游。

## 3. 验证

Cline 会大量调用工具，接入后在 `http://127.0.0.1:8787/console/` 的「调用图谱」中能清晰看到 Agent 与各工具的调用链，风险工具会高亮。

## 启用管控

确认观察结果无误后，把 `halo.yaml` 的 `mode` 改为 `guard` 并配置策略文件，敏感调用将被拦截或转人工审批。

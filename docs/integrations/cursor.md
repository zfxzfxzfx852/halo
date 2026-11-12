# 接入 Cursor

Cursor 支持覆盖 OpenAI 的 Base URL，在图形设置中即可完成接入。

## 1. 配置上游

在 `halo.yaml` 中设置一个 OpenAI 兼容的上游地址：

```yaml
mode: observe
upstream:
  baseUrl: "https://api.openai.com"
```

启动 Halo：

```sh
halo
```

## 2. 在 Cursor 中指向 Halo

1. 打开 `Settings`（设置）→ `Models`（模型）。
2. 开启 `OpenAI API Key`，填入你的 OpenAI Key（或任意占位，Halo 会转发到真实上游）。
3. 找到 `Override OpenAI Base URL`，填入：
   ```
   http://127.0.0.1:8787
   ```
4. 保存设置。

## 3. 验证

在 Cursor 的聊天或 Agent 模式中发起一次请求，打开 `http://127.0.0.1:8787/console/`，应能看到实时事件；切换到「调用图谱」可查看工具调用关系。

## 恢复默认

清空 `Override OpenAI Base URL` 即可恢复 Cursor 官方端点。

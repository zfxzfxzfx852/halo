# Halo 产品演示视频 · 分镜脚本

> 目标时长 **75 秒**，16:9（1920×1080），同时输出 MP4（X / B站 / YouTube / PH）与 GIF（README / Product Hunt）。
> 每个镜头给出：画面、操作、字幕（中/英）、旁白（中/英）。录屏在干净的演示环境进行，分辨率 1920×1080，鼠标移动放慢、点击有停顿。

---

## 镜头 1 · 开场标题卡（0:00–0:05）

- **画面**：深色背景，Halo 品牌标（渐变光环）居中，柔和光晕缓慢呼吸
- **字幕**
  - 中：`Halo — 看见并守护 Agent 的每一次调用`
  - 英：`Halo — See and secure every AI agent call`
- **旁白**
  - 中：你的 AI Agent，正在自己调用工具、访问网络、读写文件。
  - 英：Your AI agent is calling tools, hitting the network, reading and writing files — on its own.

## 镜头 2 · 痛点（0:05–0:13）

- **画面**：终端里 Agent 自动执行一连串命令，画面略快、略带压迫感
- **字幕**
  - 中：`但它到底做了什么？你看不见，也拦不住。`
  - 英：`But what is it actually doing? You can't see it — and you can't stop it.`
- **旁白**
  - 中：但它究竟做了什么，你看不见，也拦不住。
  - 英：But what it actually does is invisible — and unstoppable.

## 镜头 3 · 一行接入（0:13–0:21）

- **画面**：终端特写，逐字敲出安装命令并运行，`halo init` 向导问答
- **屏幕操作**：
  ```bash
  curl -fsSL https://raw.githubusercontent.com/zfxzfxzfx852/halo/main/scripts/install.sh | sh
  halo init
  halo
  ```
- **字幕**
  - 中：`不改一行 Agent 代码，一行命令接入。`
  - 英：`Zero code changes. One command to get started.`
- **旁白**
  - 中：不改 Agent 一行代码，Halo 透明接入，本地运行。
  - 英：Without touching agent code, Halo plugs in transparently and runs locally.

## 镜头 4 · 实时事件流（0:21–0:33）

- **画面**：控制台「事件流」，新事件实时逐条滑入，风险标签红/橙/灰
- **屏幕操作**：保持事件流视图，触发若干工具调用，停顿展示风险颜色
- **字幕**
  - 中：`每一次调用，实时可见，风险一目了然。`
  - 英：`Every call, in real time. Risk visible at a glance.`
- **旁白**
  - 中：每一次调用实时出现在事件流里，风险等级一目了然。
  - 英：Every call streams into the console in real time, with risk levels at a glance.

## 镜头 5 · 调用图谱（0:33–0:45）

- **画面**：切到「调用图谱」，Agent 节点连出各工具，风险工具红色高亮，轻微缩放
- **屏幕操作**：切图谱，鼠标悬停风险节点，缓慢放大
- **字幕**
  - 中：`交互式调用图谱，像 Wireshark 一样读懂 Agent。`
  - 英：`An interactive call graph — Wireshark for your agents.`
- **旁白**
  - 中：交互式调用图谱，让你像用 Wireshark 一样读懂 Agent 的行为。
  - 英：An interactive call graph lets you read agent behavior like Wireshark.

## 镜头 6 · 详情抽屉（0:45–0:53）

- **画面**：点一个风险节点，右侧圆角抽屉滑出，展示完整 JSON（参数、时间、风险）
- **屏幕操作**：点击 `read_file` 行，展示 `/etc/passwd` 参数
- **字幕**
  - 中：`点开任意调用，查看完整参数与上下文。`
  - 英：`Open any call to inspect full arguments and context.`
- **旁白**
  - 中：点开任意调用，完整参数和上下文都在这里。
  - 英：Open any call and every argument, every context is right there.

## 镜头 7 · 策略 / 审批 / 沙箱（0:53–1:05）

- **画面**：危险调用触发策略，进入「审批」卡片，手动点「拒绝 / 允许」
- **屏幕操作**：触发一条高危调用 → 审批挂起 → 点拒绝；再触发一条 → 允许
- **字幕**
  - 中：`危险操作，先问过你；高危命令，沙箱里执行。`
  - 英：`Risky actions ask you first. Dangerous commands run sandboxed.`
- **旁白**
  - 中：危险操作先经过你的审批，高危命令在沙箱中隔离执行。
  - 英：Risky actions ask for your approval, and dangerous commands run sandboxed.

## 镜头 8 · 审计 + 结尾（1:05–1:15）

- **画面**：展示审计校验 `valid: true`，随后切回品牌结尾卡，GitHub 地址与 CTA
- **字幕**
  - 中：`防篡改审计，全程可回放。开源 MIT，本地优先。`
  - 英：`Tamper-proof audit, fully replayable. Open source, local-first.`
- **旁白**
  - 中：防篡改哈希链让全程可回放、可验证。Halo，开源 MIT，本地优先。现在就去 GitHub 加星。
  - 英：A tamper-proof hash chain makes everything replayable and verifiable. Halo, open source under MIT, local-first. Star it on GitHub.
- **结尾文字**：`github.com/zfxzfxzfx852/halo`

---

## 制作备注

- 配色与控制台一致：iOS 系统色，蓝 `#007AFF`、红 `#FF3B30`、橙 `#FF9500`、绿 `#34C759`
- 字体：SF Pro / 系统字体；代码用 SF Mono
- 转场：150–250ms 交叉淡化，避免花哨特效
- 背景音乐：低音量、轻快科技感；旁白处自动压低（ducking）
- 输出：
  - MP4 H.264，1920×1080，30fps
  - GIF 压缩版（≤8MB，宽度 720）用于 README / Product Hunt
  - 竖版 9:16 剪辑（镜头 3/4/7）用于 X / 抖音 / 视频号

/* Halo landing — interactions */

/* ---------- i18n ---------- */
const html = document.documentElement;
const langBtn = document.getElementById('langToggle');
let lang = 'zh';

function applyLang(next) {
  lang = next;
  html.lang = lang === 'en' ? 'en' : 'zh-CN';
  langBtn.textContent = lang === 'en' ? '中文' : 'EN';

  document.querySelectorAll('[data-en]').forEach((el) => {
    if (!el.dataset.zh) el.dataset.zh = el.textContent.trim();
    el.textContent = lang === 'en' ? el.dataset.en : el.dataset.zh;
  });

  document.querySelectorAll('img[data-en-alt]').forEach((el) => {
    if (!el.dataset.zhAlt) el.dataset.zhAlt = el.alt;
    el.alt = lang === 'en' ? el.dataset.enAlt : el.dataset.zhAlt;
  });

  document.title = lang === 'en'
    ? 'Halo — Local-first observability & security plane for AI agents'
    : 'Halo — 本地优先的 AI Agent 可观测与安全控制平面';
}
langBtn.addEventListener('click', () => applyLang(lang === 'zh' ? 'en' : 'zh'));

/* ---------- copy install command ---------- */
const copyBtn = document.getElementById('copyBtn');
const installCmd = document.getElementById('installCmd');
copyBtn.addEventListener('click', async () => {
  try {
    await navigator.clipboard.writeText(installCmd.textContent);
  } catch (e) {
    const ta = document.createElement('textarea');
    ta.value = installCmd.textContent;
    document.body.appendChild(ta);
    ta.select();
    document.execCommand('copy');
    ta.remove();
  }
  const original = copyBtn.textContent;
  copyBtn.textContent = lang === 'en' ? 'Copied' : '已复制';
  setTimeout(() => (copyBtn.textContent = original), 1600);
});

/* ---------- code tabs ---------- */
document.querySelectorAll('.code-tab').forEach((tab) => {
  tab.addEventListener('click', () => {
    const target = tab.dataset.tab;
    document.querySelectorAll('.code-tab').forEach((t) => t.classList.toggle('active', t === tab));
    document.querySelectorAll('.code-pane').forEach((p) =>
      p.classList.toggle('active', p.dataset.pane === target)
    );
  });
});

/* ---------- reveal on scroll ---------- */
const io = new IntersectionObserver(
  (entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        entry.target.classList.add('in');
        io.unobserve(entry.target);
      }
    });
  },
  { threshold: 0.12 }
);
document.querySelectorAll('.reveal').forEach((el) => io.observe(el));

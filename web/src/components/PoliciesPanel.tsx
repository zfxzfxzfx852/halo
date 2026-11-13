import type { PolicyRule } from '../types'
import { Icon } from './Icon'
import { decisionLabel } from '../lib/ui'

interface Props {
  defaultDecision: string
  rules: PolicyRule[]
  mode: string
  onCopy: (text: string, label: string) => void
}

const EXAMPLE = `version: "1"
default: allow
rules:
  - name: block-rm
    command: "rm *"
    decision: deny
  - name: review-writes
    tool: write_file
    decision: ask`

export function PoliciesPanel({ defaultDecision, rules, mode, onCopy }: Props) {
  const guard = mode === 'guard'
  return (
    <section className="policies-view">
      <article className={`card glass mode-card ${guard ? 'guard' : 'observe'}`}>
        <div className="mode-copy">
          <span className="eyebrow">
            <Icon name={guard ? 'shield' : 'pulse'} size={16} />
            当前模式
          </span>
          <h2>{guard ? 'Guard · 主动管控' : 'Observe · 只观测'}</h2>
          <p>
            {guard
              ? '策略命中 deny 直接 403；命中 ask 挂起等待人工审批；其余放行并完整记录。'
              : '所有调用透明转发，只做记录、检测与可视化，不阻断任何请求。可随时切换为 guard。'}
          </p>
        </div>
        <div className="mode-facts">
          <div>
            <strong>{rules.length}</strong>
            <span>条规则</span>
          </div>
          <div>
            <strong>{decisionLabel(defaultDecision === 'deny' ? 'deny' : defaultDecision === 'ask' ? 'ask' : 'allow')}</strong>
            <span>默认决策</span>
          </div>
        </div>
      </article>

      {rules.length === 0 ? (
        <article className="card glass">
          <header className="card-head">
            <h2>
              <Icon name="shield" size={18} />
              还没有策略规则
            </h2>
            <span className="subtle">fail-closed</span>
          </header>
          <p className="card-desc">
            在 <code>halo.yaml</code> 的 <code>policyFile</code> 指向一个 YAML 文件即可启用。
            规则按顺序匹配，首条命中生效；多个条件是 AND 关系。
          </p>
          <pre className="code-block">{EXAMPLE}</pre>
          <button className="btn ghost" onClick={() => onCopy(EXAMPLE, '策略示例已复制')}>
            <Icon name="copy" size={16} />
            复制策略模板
          </button>
        </article>
      ) : (
        <div className="rule-list">
          {rules.map((r, i) => (
            <article className="rule-card glass" key={`${r.name}-${i}`}>
              <div className="rule-top">
                <span className={`decision-pill ${r.decision}`}>{decisionLabel(r.decision)}</span>
                <h3>{r.name || `规则 ${i + 1}`}</h3>
                {r.maxRisk && <span className="subtle">≤ {r.maxRisk}</span>}
              </div>
              <div className="rule-conditions">
                {r.tool && <code>tool: {r.tool}</code>}
                {r.method && <code>method: {r.method}</code>}
                {r.command && <code>command: {r.command}</code>}
                {r.paths?.map((p) => <code key={p}>path: {p}</code>)}
                {r.args &&
                  Object.entries(r.args).map(([k, v]) => (
                    <code key={k}>args.{k}: {v}</code>
                  ))}
                {!r.tool && !r.method && !r.command && !r.paths?.length && !r.args && (
                  <span className="subtle">匹配全部调用</span>
                )}
              </div>
            </article>
          ))}
        </div>
      )}
    </section>
  )
}

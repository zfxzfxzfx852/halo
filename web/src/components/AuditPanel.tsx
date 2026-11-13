import { useState } from 'react'
import type { AuditReport } from '../types'
import { downloadAudit, verifyAudit } from '../lib/api'
import { Icon } from './Icon'

interface Props {
  onToast: (msg: string) => void
}

export function AuditPanel({ onToast }: Props) {
  const [report, setReport] = useState<AuditReport>()
  const [loading, setLoading] = useState(false)
  const [exporting, setExporting] = useState(false)

  const runVerify = async () => {
    setLoading(true)
    try {
      const r = await verifyAudit()
      setReport(r)
      onToast(r.valid ? '审计链校验通过' : '审计链校验失败，请查看定位信息')
    } catch {
      onToast('审计校验请求失败')
    } finally {
      setLoading(false)
    }
  }

  const runExport = async () => {
    setExporting(true)
    try {
      const lines = await downloadAudit()
      onToast(`已导出 ${lines} 条审计记录`)
    } catch {
      onToast('导出失败')
    } finally {
      setExporting(false)
    }
  }

  return (
    <section className="audit-view">
      <div className="two-col">
        <article className="card glass audit-card">
          <header className="card-head">
            <h2>
              <Icon name="shield" size={18} />
              审计链
            </h2>
            <span className="subtle">SHA-256 哈希链</span>
          </header>
          <p className="card-desc">
            每条事件写入后按 <code>sha256(prevHash + payload)</code> 串联。
            篡改、删除或断链都会被识别；保留策略清理的旧事件不会误报，只计入缺失。
          </p>
          <div className="audit-actions">
            <button className="btn primary" disabled={loading} onClick={runVerify}>
              <Icon name="refresh" size={16} />
              {loading ? '校验中…' : '立即校验'}
            </button>
            <button className="btn ghost" disabled={exporting} onClick={runExport}>
              <Icon name="download" size={16} />
              {exporting ? '导出中…' : '导出 JSONL'}
            </button>
          </div>
        </article>

        <article className={`card glass audit-result ${report ? (report.valid ? 'ok' : 'bad') : ''}`}>
          {!report ? (
            <div className="audit-idle">
              <div className="empty-orb small">
                <Icon name="seal" size={22} />
              </div>
              <h3>尚未校验</h3>
              <p>点击左侧"立即校验"，从创世记录开始重算整条链。</p>
            </div>
          ) : (
            <>
              <div className={`audit-badge ${report.valid ? 'ok' : 'bad'}`}>
                <Icon name={report.valid ? 'check' : 'alert'} size={24} />
                {report.valid ? '链完整' : '发现篡改或断链'}
              </div>
              <div className="audit-stats">
                <div>
                  <strong>{report.checked}</strong>
                  <span>已校验</span>
                </div>
                <div>
                  <strong>{report.missingEvents}</strong>
                  <span>已清理</span>
                </div>
                <div>
                  <strong>{report.valid ? '0' : report.brokenSeq ?? '—'}</strong>
                  <span>断点序号</span>
                </div>
              </div>
              {report.reason && <p className="audit-reason">{report.reason}</p>}
            </>
          )}
        </article>
      </div>

      <article className="card glass">
        <header className="card-head">
          <h2>
            <Icon name="book" size={18} />
            怎么用这条链
          </h2>
        </header>
        <div className="guide-steps">
          <div>
            <b>01</b>
            <p>启动时从最后一条记录恢复链头，崩溃重启不会断链。</p>
          </div>
          <div>
            <b>02</b>
            <p>事件落库后立即追加审计记录，审批结果同样进入链。</p>
          </div>
          <div>
            <b>03</b>
            <p>导出 JSONL 交给合规或第三方复核，可独立重算。</p>
          </div>
        </div>
      </article>
    </section>
  )
}

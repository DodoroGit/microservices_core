import React, { useState } from 'react'
import { aiAPI } from '../api'

const CARDS = [
  {
    key: 'background',
    title: '個人背景介紹',
    subtitle: '根據「履歷」筆記生成',
    desc: '整合工作經歷與職涯背景，生成適合放在履歷或 LinkedIn 的自我介紹。',
    apiFn: () => aiAPI.generateBackground(),
    render: (data) => <BackgroundResult data={data} />,
  },
  {
    key: 'projects',
    title: '專案介紹',
    subtitle: '根據「專案」筆記生成',
    desc: '將所有專案筆記歸納整理，每個專案各自輸出一段技術描述。',
    apiFn: () => aiAPI.generateProjects(),
    render: (data) => <ProjectsResult data={data} />,
  },
  {
    key: 'skills',
    title: '技術能力總覽',
    subtitle: '根據「技術筆記」生成',
    desc: '從技術學習筆記中萃取技能，依語言、框架、資料庫、工具分類整理。',
    apiFn: () => aiAPI.generateSkills(),
    render: (data) => <SkillsResult data={data} />,
  },
  {
    key: 'daily',
    title: '生活雜談',
    subtitle: '根據「日常隨筆」生成',
    desc: '從日常記錄中提煉出輕鬆自然的個人生活介紹，適合個人網站。',
    apiFn: () => aiAPI.generateDaily(),
    render: (data) => <DailyResult data={data} />,
  },
]

// ── 各類型結果元件 ──────────────────────────────────────────────────────────

function TagList({ items }) {
  return (
    <div className="ailab-tags">
      {items.map((item, i) => (
        <span key={i} className="ailab-tag">{item}</span>
      ))}
    </div>
  )
}

function BackgroundResult({ data }) {
  return (
    <>
      <div className="ailab-result-section">
        <div className="ailab-section-label">背景介紹</div>
        <p className="ailab-section-content">{data.background_intro}</p>
      </div>
      <div className="ailab-result-section">
        <div className="ailab-section-label">職位 / 角色</div>
        <TagList items={data.key_roles} />
      </div>
      <div className="ailab-result-section last">
        <div className="ailab-section-label">年資摘要</div>
        <p className="ailab-section-content">{data.years_summary}</p>
      </div>
    </>
  )
}

function ProjectsResult({ data }) {
  return (
    <>
      {data.projects.map((proj, i) => (
        <div key={i} className={`ailab-result-section ${i === data.projects.length - 1 ? 'last' : ''}`}>
          <div className="ailab-section-label">{proj.project_name}</div>
          <p className="ailab-section-content">{proj.description}</p>
          <div style={{ marginTop: 10 }}>
            <TagList items={proj.tech_keywords} />
          </div>
        </div>
      ))}
    </>
  )
}

function SkillsResult({ data }) {
  const categories = [
    { label: '程式語言', items: data.languages },
    { label: '框架 / 函式庫', items: data.frameworks },
    { label: '資料庫', items: data.databases },
    { label: '工具 / 平台', items: data.tools },
  ]
  return (
    <>
      <div className="ailab-result-section">
        <div className="ailab-section-label">技術概述</div>
        <p className="ailab-section-content">{data.summary}</p>
      </div>
      {categories.map(({ label, items }) =>
        items.length > 0 ? (
          <div key={label} className="ailab-result-section">
            <div className="ailab-section-label">{label}</div>
            <TagList items={items} />
          </div>
        ) : null
      )}
    </>
  )
}

function DailyResult({ data }) {
  return (
    <>
      <div className="ailab-result-section">
        <div className="ailab-section-label">生活介紹</div>
        <p className="ailab-section-content">{data.introduction}</p>
      </div>
      <div className="ailab-result-section last">
        <div className="ailab-section-label">興趣 / 關鍵詞</div>
        <TagList items={data.interests} />
      </div>
    </>
  )
}

// ── 單張功能卡片 ────────────────────────────────────────────────────────────

function AICard({ card }) {
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState(null)
  const [error, setError] = useState('')

  const handleGenerate = async () => {
    if (loading) return
    setError('')
    setResult(null)
    setLoading(true)
    try {
      const res = await card.apiFn()
      setResult(res.data)
    } catch (err) {
      setError(err.response?.data?.detail || '生成失敗，請稍後再試')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="aicard">
      <div className="aicard-header">
        <div>
          <h3 className="aicard-title">{card.title}</h3>
          <span className="aicard-subtitle">{card.subtitle}</span>
        </div>
        <button
          className="btn-primary aicard-btn"
          onClick={handleGenerate}
          disabled={loading}
        >
          {loading ? '生成中...' : '生成'}
        </button>
      </div>

      <p className="aicard-desc">{card.desc}</p>

      {loading && (
        <div className="ailab-loading">
          <div className="ailab-loading-dots">
            <span /><span /><span />
          </div>
        </div>
      )}

      {error && <div className="chat-error" style={{ margin: '12px 0 0' }}>{error}</div>}

      {result && (
        <div className="aicard-result">
          <div className="aicard-result-meta">來源筆記：{result.source_notes_count} 篇</div>
          {card.render(result)}
        </div>
      )}
    </div>
  )
}

// ── 主頁面 ──────────────────────────────────────────────────────────────────

function AILab({ user }) {
  return (
    <div className="ailab-page">
      <div className="page-header dark-header">
        <div>
          <h1>AI Lab</h1>
          <p>從筆記自動生成個人介紹素材 · 由 Claude 驅動</p>
        </div>
      </div>

      <div className="ailab-grid">
        {CARDS.map(card => (
          <AICard key={card.key} card={card} />
        ))}
      </div>
    </div>
  )
}

export default AILab

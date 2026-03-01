import React from 'react'
import { Link } from 'react-router-dom'

const SERVICES = [
  { name: 'API Gateway',   status: 'online' },
  { name: 'User Service',  status: 'online' },
  { name: 'PostgreSQL',    status: 'online' },
  { name: 'Notes Service', status: 'soon' },
  { name: 'AI Lab',        status: 'soon' },
]

const STACK = ['Go', 'React', 'PostgreSQL', 'Docker', 'Microservices', 'API Gateway', 'JWT', 'Layer Architecture', 'Unit Test', 'Integration Test']

function Home({ user }) {
  return (
    <div className="home-page">

      {/* ── Hero ── */}
      <section className="hero-section">
        <div className="hero-glow hero-glow-purple" />
        <div className="hero-glow hero-glow-cyan" />
        <div className="hero-dots" />

        <div className="hero-inner">
          <div className="hero-content">
            <div className="hero-badge">
              <span className="hero-badge-dot" />
              Personal Lab
            </div>
            <h1 className="hero-title">Rick's Lab</h1>
            <p className="hero-subtitle">
              記錄學習、實驗技術、探索 AI 的個人空間
            </p>
            <div className="hero-actions">
              {user ? (
                <Link to="/dashboard" className="btn-gradient">進入 Dashboard →</Link>
              ) : (
                <>
                  <Link to="/register" className="btn-gradient">開始使用 →</Link>
                  <Link to="/login" className="btn-glass">登入</Link>
                </>
              )}
            </div>
          </div>

          <div className="status-card">
            <div className="status-card-header">
              <span className="status-card-title">Service Status</span>
              <span className="status-card-dot" />
            </div>
            <div className="status-card-list">
              {SERVICES.map(svc => (
                <div className="status-row" key={svc.name}>
                  <span className={`status-indicator ${svc.status === 'online' ? 'status-online' : 'status-soon'}`} />
                  <span className="status-name">{svc.name}</span>
                  <span className={`status-badge ${svc.status === 'online' ? 'badge-online' : 'badge-soon'}`}>
                    {svc.status}
                  </span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* ── Features ── */}
      <section className="features-section">
        <div className="features-header">
          <span className="section-label">WHAT'S INSIDE</span>
          <h2 className="section-title">Lab 功能介紹</h2>
          <p className="section-desc">微服務架構的個人實驗場，持續建構中</p>
        </div>
        <div className="features-grid">
          <div className="feature-card">
            <div className="feature-icon-wrap icon-blue">
              <span>👤</span>
            </div>
            <h3>使用者系統</h3>
            <p>JWT 身份驗證，API Gateway 統一入口，保護後端微服務不被直接存取。</p>
            <div className="feature-tag">上線中</div>
          </div>
          <div className="feature-card">
            <div className="feature-icon-wrap icon-purple">
              <span>📝</span>
            </div>
            <h3>技術筆記</h3>
            <p>記錄每天學習的技術筆記、做過的專案，以及每日閱讀的內容整理。</p>
            <div className="feature-tag tag-soon">Coming Soon</div>
          </div>
          <div className="feature-card">
            <div className="feature-icon-wrap icon-cyan">
              <span>🤖</span>
            </div>
            <h3>AI Lab</h3>
            <p>整合本地 Ollama 模型，在自己的 server 上跑 LLM，完全掌控資料。</p>
            <div className="feature-tag tag-soon">Coming Soon</div>
          </div>
        </div>
      </section>

      {/* ── Stack ── */}
      <section className="stack-section">
        <div className="stack-glow" />
        <span className="section-label light">TECH STACK</span>
        <h2 className="section-title light">Rick's Lab 使用技術一覽</h2>
        <div className="stack-grid">
          {STACK.map(tech => (
            <div className="stack-chip" key={tech}>{tech}</div>
          ))}
        </div>
      </section>

    </div>
  )
}

export default Home

import React from 'react'
import { Link } from 'react-router-dom'

const codeSnippet = [
  { type: 'comment', text: '// Rick\'s Lab — microservices' },
  { type: 'keyword', text: 'const ', suffix: { type: 'var', text: 'lab' }, end: ' = {' },
  { type: 'indent', key: 'service', value: '"api-gateway"' },
  { type: 'indent', key: 'auth', value: '"JWT + middleware"' },
  { type: 'indent', key: 'stack', value: '["Go", "React", "Docker"]' },
  { type: 'indent', key: 'next', value: '"Ollama integration"' },
  { type: 'plain', text: '}' },
]

const STACK = ['Go', 'React', 'PostgreSQL', 'Redis', 'Docker', 'API Gateway', 'JWT', 'Ollama']

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

          <div className="hero-terminal">
            <div className="terminal-header">
              <span className="terminal-dot dot-red" />
              <span className="terminal-dot dot-yellow" />
              <span className="terminal-dot dot-green" />
              <span className="terminal-title">lab.ts</span>
            </div>
            <div className="terminal-body">
              {codeSnippet.map((line, i) => (
                <div className="code-line" key={i}>
                  {line.type === 'comment' && <span className="c-comment">{line.text}</span>}
                  {line.type === 'keyword' && (
                    <>
                      <span className="c-keyword">{line.text}</span>
                      <span className="c-var">{line.suffix.text}</span>
                      <span className="c-plain">{line.end}</span>
                    </>
                  )}
                  {line.type === 'indent' && (
                    <span className="c-indent">
                      {'  '}<span className="c-key">{line.key}</span>
                      <span className="c-plain">: </span>
                      <span className="c-string">{line.value}</span>
                      <span className="c-plain">,</span>
                    </span>
                  )}
                  {line.type === 'plain' && <span className="c-plain">{line.text}</span>}
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
          <h2 className="section-title">Lab 裡有什麼</h2>
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
        <h2 className="section-title light">建構這個 Lab 用的東西</h2>
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

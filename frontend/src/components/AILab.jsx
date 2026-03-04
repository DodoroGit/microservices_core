import React, { useState, useEffect, useRef } from 'react'
import { aiAPI } from '../api'

function AILab({ user }) {
  const [messages, setMessages] = useState([])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const bottomRef = useRef(null)

  useEffect(() => {
    if (user?.id) fetchHistory()
  }, [user])

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const fetchHistory = async () => {
    try {
      const res = await aiAPI.getHistory(user.id)
      setMessages(res.data.history || [])
    } catch {
      // 歷史載入失敗不阻斷主流程
    }
  }

  const handleSend = async () => {
    const text = input.trim()
    if (!text || loading) return

    setInput('')
    setError('')
    setMessages(prev => [...prev, { role: 'user', content: text }])
    setLoading(true)

    try {
      const res = await aiAPI.sendMessage(text, user.id)
      setMessages(res.data.history)
    } catch (err) {
      setError(err.response?.data?.detail || '傳送失敗，請稍後再試')
      setMessages(prev => prev.slice(0, -1))
    } finally {
      setLoading(false)
    }
  }

  const handleClear = async () => {
    try {
      await aiAPI.clearHistory(user.id)
      setMessages([])
    } catch {
      setError('清除失敗')
    }
  }

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  return (
    <div className="ailab-page">
      <div className="page-header dark-header">
        <div>
          <h1>AI Lab</h1>
          <p>本地 LLM 對話介面 · llama3.2:3b</p>
        </div>
        {messages.length > 0 && (
          <button className="btn-ghost" onClick={handleClear}>清除紀錄</button>
        )}
      </div>

      <div className="chat-container">
        <div className="chat-messages">
          {messages.length === 0 && !loading && (
            <div className="chat-empty">
              <div className="chat-empty-icon">🤖</div>
              <p style={{ color: '#000' }}>開始和 AI 對話吧</p>
              <span style={{ color: '#000' }}>支援中文 · 對話紀錄會自動儲存 24 小時</span>
            </div>
          )}

          {messages.map((msg, i) => (
            <div key={i} className={`chat-bubble-wrap ${msg.role === 'user' ? 'bubble-user' : 'bubble-ai'}`}>
              <div className="chat-bubble">
                <p>{msg.content}</p>
              </div>
            </div>
          ))}

          {loading && (
            <div className="chat-bubble-wrap bubble-ai">
              <div className="chat-bubble bubble-typing">
                <span /><span /><span />
              </div>
            </div>
          )}

          {error && (
            <div className="chat-error">{error}</div>
          )}

          <div ref={bottomRef} />
        </div>

        <div className="chat-input-bar">
          <textarea
            className="chat-input"
            value={input}
            onChange={e => setInput(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="輸入訊息，Enter 傳送，Shift+Enter 換行"
            rows={1}
            disabled={loading}
          />
          <button
            className="btn-primary chat-send-btn"
            onClick={handleSend}
            disabled={loading || !input.trim()}
          >
            傳送
          </button>
        </div>
      </div>
    </div>
  )
}

export default AILab

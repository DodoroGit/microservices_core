import React, { useState } from 'react'

const CATEGORIES = ['技術筆記', '專案', '每日閱讀']

function Notes({ user }) {
  const [showModal, setShowModal] = useState(false)
  const [form, setForm] = useState({ title: '', category: '技術筆記', content: '' })

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  const handleSubmit = (e) => {
    e.preventDefault()
    // TODO: 後端 API 完成後串接
    alert('後端功能開發中，Coming Soon！')
    setShowModal(false)
    setForm({ title: '', category: '技術筆記', content: '' })
  }

  return (
    <div className="notes-page">
      <div className="page-header dark-header">
        <div>
          <h1>技術筆記</h1>
          <p>記錄你的學習軌跡</p>
        </div>
        {user?.role === 'admin' && (
          <button className="btn-primary" onClick={() => setShowModal(true)}>+ 新增內容</button>
        )}
      </div>

      <div className="coming-soon-section">
        <div className="coming-soon-icon">🚧</div>
        <h2>後端功能開發中</h2>
        <p>筆記儲存功能即將上線，目前可以點選「新增內容」預覽新增流程。</p>
        {user?.role === 'admin' && (
          <button className="btn-outline-light" onClick={() => setShowModal(true)}>預覽新增功能</button>
        )}
      </div>

      {showModal && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>新增內容</h3>
              <button className="modal-close" onClick={() => setShowModal(false)}>✕</button>
            </div>
            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label>類別</label>
                <select name="category" value={form.category} onChange={handleChange}>
                  {CATEGORIES.map(c => <option key={c}>{c}</option>)}
                </select>
              </div>
              <div className="form-group">
                <label>標題</label>
                <input
                  name="title"
                  value={form.title}
                  onChange={handleChange}
                  required
                  placeholder="輸入標題"
                />
              </div>
              <div className="form-group">
                <label>內容</label>
                <textarea
                  name="content"
                  value={form.content}
                  onChange={handleChange}
                  rows={6}
                  required
                  placeholder="輸入內容..."
                />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn-ghost" onClick={() => setShowModal(false)}>取消</button>
                <button type="submit" className="btn-primary">新增</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

export default Notes

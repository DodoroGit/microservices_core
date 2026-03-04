import React, { useState, useEffect } from 'react'
import { noteAPI } from '../api'

const CATEGORIES = ['專案', '技術筆記', '日常隨筆']

function Notes({ user }) {
  const [notes, setNotes] = useState([])
  const [activeCategory, setActiveCategory] = useState(null)
  const [showModal, setShowModal] = useState(false)
  const [form, setForm] = useState({ title: '', category: '專案', content: '' })
  const [error, setError] = useState('')
  const [selectedNote, setSelectedNote] = useState(null)
  const [editNote, setEditNote] = useState(null)
  const [editForm, setEditForm] = useState({ title: '', category: '專案', content: '' })
  const [deleteConfirmId, setDeleteConfirmId] = useState(null)

  useEffect(() => {
    fetchNotes(activeCategory)
  }, [activeCategory])

  const fetchNotes = async (category) => {
    try {
      const res = await noteAPI.getNotes(category)
      setNotes(res.data)
    } catch {
      setError('載入筆記失敗')
    }
  }

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value })
  }

  const handleEditChange = (e) => {
    setEditForm({ ...editForm, [e.target.name]: e.target.value })
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      await noteAPI.createNote(form)
      setShowModal(false)
      setForm({ title: '', category: '專案', content: '' })
      fetchNotes(activeCategory)
    } catch {
      setError('新增失敗，請稍後再試')
    }
  }

  const openEdit = (note) => {
    setEditNote(note)
    setEditForm({ title: note.title, category: note.category, content: note.content })
    setSelectedNote(null)
  }

  const handleUpdate = async (e) => {
    e.preventDefault()
    try {
      await noteAPI.updateNote(editNote.id, editForm)
      setEditNote(null)
      fetchNotes(activeCategory)
    } catch {
      setError('更新失敗，請稍後再試')
    }
  }

  const handleDelete = async () => {
    try {
      await noteAPI.deleteNote(deleteConfirmId)
      setDeleteConfirmId(null)
      setSelectedNote(null)
      fetchNotes(activeCategory)
    } catch {
      setError('刪除失敗')
    }
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

      <div className="notes-filter">
        <button
          className={activeCategory === null ? 'filter-btn active' : 'filter-btn'}
          onClick={() => setActiveCategory(null)}
        >
          全部
        </button>
        {CATEGORIES.map(c => (
          <button
            key={c}
            className={activeCategory === c ? 'filter-btn active' : 'filter-btn'}
            onClick={() => setActiveCategory(c)}
          >
            {c}
          </button>
        ))}
      </div>

      {error && <div className="chat-error">{error}</div>}

      <div className="notes-grid">
        {notes.length === 0 ? (
          <div className="notes-empty">
            <p>目前沒有筆記</p>
          </div>
        ) : (
          notes.map(note => (
            <div key={note.id} className="note-card" onClick={() => setSelectedNote(note)}>
              <div className="note-card-header">
                <span className="note-category">{note.category}</span>
              </div>
              <h3 className="note-title">{note.title}</h3>
              <p className="note-content">{note.content}</p>
              <span className="note-date">
                {new Date(note.created_at).toLocaleDateString('zh-TW')}
              </span>
            </div>
          ))
        )}
      </div>

      {/* 詳細內容 modal */}
      {selectedNote && (
        <div className="modal-overlay" onClick={() => setSelectedNote(null)}>
          <div className="modal note-detail-modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <div className="note-detail-meta">
                <span className="note-category">{selectedNote.category}</span>
                <span className="note-detail-date">
                  {new Date(selectedNote.created_at).toLocaleDateString('zh-TW')}
                </span>
              </div>
              <button className="modal-close" onClick={() => setSelectedNote(null)}>✕</button>
            </div>
            <h2 className="note-detail-title">{selectedNote.title}</h2>
            <div className="note-detail-content">{selectedNote.content}</div>
            {user?.role === 'admin' && (
              <div className="note-detail-actions">
                <button className="btn-ghost" onClick={() => openEdit(selectedNote)}>編輯</button>
                <button className="btn-danger" onClick={() => { setDeleteConfirmId(selectedNote.id); setSelectedNote(null) }}>刪除</button>
              </div>
            )}
          </div>
        </div>
      )}

      {/* 編輯 modal */}
      {editNote && (
        <div className="modal-overlay" onClick={() => setEditNote(null)}>
          <div className="modal" onClick={e => e.stopPropagation()}>
            <div className="modal-header">
              <h3>編輯筆記</h3>
              <button className="modal-close" onClick={() => setEditNote(null)}>✕</button>
            </div>
            <form onSubmit={handleUpdate}>
              <div className="form-group">
                <label>類別</label>
                <select name="category" value={editForm.category} onChange={handleEditChange}>
                  {CATEGORIES.map(c => <option key={c}>{c}</option>)}
                </select>
              </div>
              <div className="form-group">
                <label>標題</label>
                <input
                  name="title"
                  value={editForm.title}
                  onChange={handleEditChange}
                  required
                  placeholder="輸入標題"
                />
              </div>
              <div className="form-group">
                <label>內容</label>
                <textarea
                  name="content"
                  value={editForm.content}
                  onChange={handleEditChange}
                  rows={6}
                  required
                  placeholder="輸入內容..."
                />
              </div>
              <div className="modal-actions">
                <button type="button" className="btn-ghost" onClick={() => setEditNote(null)}>取消</button>
                <button type="submit" className="btn-primary">儲存</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* 刪除確認 dialog */}
      {deleteConfirmId && (
        <div className="modal-overlay" onClick={() => setDeleteConfirmId(null)}>
          <div className="modal confirm-modal" onClick={e => e.stopPropagation()}>
            <div className="confirm-icon">🗑</div>
            <h3 className="confirm-title">確定要刪除嗎？</h3>
            <p className="confirm-desc">刪除後無法復原</p>
            <div className="modal-actions">
              <button className="btn-ghost" onClick={() => setDeleteConfirmId(null)}>取消</button>
              <button className="btn-danger" onClick={handleDelete}>確定刪除</button>
            </div>
          </div>
        </div>
      )}

      {/* 新增 modal */}
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

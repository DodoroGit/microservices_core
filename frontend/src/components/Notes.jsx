import React, { useState, useEffect } from 'react'
import { noteAPI } from '../api'

const CATEGORIES = ['專案', '技術筆記', '履歷', '日常隨筆']

const EMPTY_FORM = {
  category: '專案',
  // 技術筆記 & 日常隨筆
  title: '', content: '',
  // 技術筆記
  tags: '',
  // 專案
  project_name: '', description: '', tech_stack: '', github_url: '', year: '',
  // 履歷
  company: '', position: '', period: '', projects: '',
}

// 把逗號分隔字串轉為陣列（供 tech_stack / tags 使用）
function toArray(str) {
  if (!str) return []
  return str.split(',').map(s => s.trim()).filter(Boolean)
}

function formToPayload(form) {
  const base = { category: form.category }
  switch (form.category) {
    case '專案':
      return { ...base, project_name: form.project_name, description: form.description,
        tech_stack: toArray(form.tech_stack), github_url: form.github_url,
        year: parseInt(form.year, 10) || 0 }
    case '技術筆記':
      return { ...base, title: form.title, content: form.content, tags: toArray(form.tags) }
    case '履歷':
      return { ...base, company: form.company, position: form.position,
        period: form.period, projects: form.projects, tech_stack: toArray(form.tech_stack) }
    case '日常隨筆':
      return { ...base, title: form.title, content: form.content }
    default:
      return base
  }
}

// ── 卡片內容（依類別） ───────────────────────────────────────────────────────

function TagRow({ items }) {
  if (!items?.length) return null
  return (
    <div className="note-tag-row">
      {items.map((t, i) => <span key={i} className="note-tag">{t}</span>)}
    </div>
  )
}

function NoteCardBody({ note }) {
  switch (note.category) {
    case '專案':
      return (
        <>
          <div className="note-card-title-row">
            <h3 className="note-title">{note.project_name}</h3>
            {note.year > 0 && <span className="note-year">{note.year}</span>}
          </div>
          <p className="note-content">{note.description}</p>
          <TagRow items={note.tech_stack} />
          {note.github_url && (
            <a className="note-github" href={note.github_url} target="_blank" rel="noreferrer"
              onClick={e => e.stopPropagation()}>
              GitHub →
            </a>
          )}
        </>
      )
    case '技術筆記':
      return (
        <>
          <h3 className="note-title">{note.title}</h3>
          <p className="note-content">{note.content}</p>
          <TagRow items={note.tags} />
        </>
      )
    case '履歷':
      return (
        <>
          <h3 className="note-title">{note.position}</h3>
          <p className="note-subtitle">{note.company}</p>
          {note.period && <p className="note-period">{note.period}</p>}
          <TagRow items={note.tech_stack} />
        </>
      )
    case '日常隨筆':
      return (
        <>
          <h3 className="note-title">{note.title}</h3>
          <p className="note-content">{note.content}</p>
        </>
      )
    default:
      return <h3 className="note-title">{note.title || note.project_name}</h3>
  }
}

// ── 詳細 modal 內容（依類別） ───────────────────────────────────────────────

function NoteDetailBody({ note }) {
  switch (note.category) {
    case '專案':
      return (
        <div className="note-detail-fields">
          <Field label="專案名稱" value={note.project_name} />
          <Field label="建置年份" value={note.year > 0 ? String(note.year) : ''} />
          <Field label="專案介紹" value={note.description} multiline />
          {note.tech_stack?.length > 0 && (
            <div className="note-detail-field">
              <span className="field-label">技術棧</span>
              <TagRow items={note.tech_stack} />
            </div>
          )}
          {note.github_url && (
            <div className="note-detail-field">
              <span className="field-label">GitHub</span>
              <a href={note.github_url} target="_blank" rel="noreferrer">{note.github_url}</a>
            </div>
          )}
        </div>
      )
    case '技術筆記':
      return (
        <div className="note-detail-fields">
          <Field label="標題" value={note.title} />
          <Field label="內容" value={note.content} multiline />
          {note.tags?.length > 0 && (
            <div className="note-detail-field">
              <span className="field-label">技術標籤</span>
              <TagRow items={note.tags} />
            </div>
          )}
        </div>
      )
    case '履歷':
      return (
        <div className="note-detail-fields">
          <Field label="任職公司" value={note.company} />
          <Field label="任職職位" value={note.position} />
          <Field label="任職日期" value={note.period} />
          <Field label="負責專案" value={note.projects} multiline />
          {note.tech_stack?.length > 0 && (
            <div className="note-detail-field">
              <span className="field-label">技術棧</span>
              <TagRow items={note.tech_stack} />
            </div>
          )}
        </div>
      )
    case '日常隨筆':
      return (
        <div className="note-detail-fields">
          <Field label="標題" value={note.title} />
          <Field label="內容" value={note.content} multiline />
        </div>
      )
    default:
      return null
  }
}

function Field({ label, value, multiline }) {
  if (!value) return null
  return (
    <div className="note-detail-field">
      <span className="field-label">{label}</span>
      {multiline
        ? <p className="field-value">{value}</p>
        : <span className="field-value">{value}</span>}
    </div>
  )
}

// ── 表單欄位（依類別） ──────────────────────────────────────────────────────

function CategoryFields({ form, onChange }) {
  switch (form.category) {
    case '專案':
      return (
        <>
          <FormInput label="專案名稱 *" name="project_name" value={form.project_name} onChange={onChange} placeholder="例：microservices_core" />
          <FormTextarea label="專案介紹 *" name="description" value={form.description} onChange={onChange} placeholder="簡介這個專案做了什麼..." />
          <FormInput label="技術棧" name="tech_stack" value={form.tech_stack} onChange={onChange} placeholder="用逗號分隔，例：Go, Docker, PostgreSQL" />
          <FormInput label="GitHub 連結" name="github_url" value={form.github_url} onChange={onChange} placeholder="https://github.com/..." />
          <FormInput label="建置年份 *" name="year" value={form.year} onChange={onChange} placeholder="例：2024" type="number" />
        </>
      )
    case '技術筆記':
      return (
        <>
          <FormInput label="標題 *" name="title" value={form.title} onChange={onChange} placeholder="輸入標題" />
          <FormTextarea label="內容 *" name="content" value={form.content} onChange={onChange} placeholder="記錄技術心得..." />
          <FormInput label="技術標籤" name="tags" value={form.tags} onChange={onChange} placeholder="用逗號分隔，例：React, Hooks, Context" />
        </>
      )
    case '履歷':
      return (
        <>
          <FormInput label="任職公司 *" name="company" value={form.company} onChange={onChange} placeholder="例：Acme Corp" />
          <FormInput label="任職職位 *" name="position" value={form.position} onChange={onChange} placeholder="例：後端工程師" />
          <FormInput label="任職日期 *" name="period" value={form.period} onChange={onChange} placeholder="例：2022/03 - 2024/06" />
          <FormTextarea label="負責專案" name="projects" value={form.projects} onChange={onChange} placeholder="描述在這份工作中負責的專案..." />
          <FormInput label="技術棧" name="tech_stack" value={form.tech_stack} onChange={onChange} placeholder="用逗號分隔，例：Go, Kubernetes, GCP" />
        </>
      )
    case '日常隨筆':
      return (
        <>
          <FormInput label="標題 *" name="title" value={form.title} onChange={onChange} placeholder="輸入標題" />
          <FormTextarea label="內容 *" name="content" value={form.content} onChange={onChange} placeholder="記錄生活點滴..." />
        </>
      )
    default:
      return null
  }
}

function FormInput({ label, name, value, onChange, placeholder, type = 'text' }) {
  return (
    <div className="form-group">
      <label>{label}</label>
      <input name={name} value={value} onChange={onChange} placeholder={placeholder} type={type} />
    </div>
  )
}

function FormTextarea({ label, name, value, onChange, placeholder }) {
  return (
    <div className="form-group">
      <label>{label}</label>
      <textarea name={name} value={value} onChange={onChange} placeholder={placeholder} rows={5} />
    </div>
  )
}

// ── 把 note 欄位轉回 form 狀態（用於編輯） ──────────────────────────────────

function noteToEditForm(note) {
  return {
    category: note.category,
    title: note.title || '',
    content: note.content || '',
    tags: (note.tags || []).join(', '),
    project_name: note.project_name || '',
    description: note.description || '',
    tech_stack: (note.tech_stack || []).join(', '),
    github_url: note.github_url || '',
    year: note.year > 0 ? String(note.year) : '',
    company: note.company || '',
    position: note.position || '',
    period: note.period || '',
    projects: note.projects || '',
  }
}

// ── 主元件 ──────────────────────────────────────────────────────────────────

function Notes({ user }) {
  const [notes, setNotes] = useState([])
  const [activeCategory, setActiveCategory] = useState(null)
  const [showModal, setShowModal] = useState(false)
  const [form, setForm] = useState(EMPTY_FORM)
  const [error, setError] = useState('')
  const [selectedNote, setSelectedNote] = useState(null)
  const [editNote, setEditNote] = useState(null)
  const [editForm, setEditForm] = useState(EMPTY_FORM)
  const [deleteConfirmId, setDeleteConfirmId] = useState(null)

  useEffect(() => { fetchNotes(activeCategory) }, [activeCategory])

  const fetchNotes = async (category) => {
    try {
      const res = await noteAPI.getNotes(category)
      setNotes(res.data)
    } catch {
      setError('載入筆記失敗')
    }
  }

  const handleChange = (e) => setForm({ ...form, [e.target.name]: e.target.value })
  const handleEditChange = (e) => setEditForm({ ...editForm, [e.target.name]: e.target.value })

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      await noteAPI.createNote(formToPayload(form))
      setShowModal(false)
      setForm(EMPTY_FORM)
      fetchNotes(activeCategory)
    } catch (err) {
      setError(err.response?.data?.error || '新增失敗，請稍後再試')
    }
  }

  const openEdit = (note) => {
    setEditNote(note)
    setEditForm(noteToEditForm(note))
    setSelectedNote(null)
  }

  const handleUpdate = async (e) => {
    e.preventDefault()
    try {
      await noteAPI.updateNote(editNote.id, formToPayload(editForm))
      setEditNote(null)
      fetchNotes(activeCategory)
    } catch (err) {
      setError(err.response?.data?.error || '更新失敗，請稍後再試')
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
          <button className="btn-primary" onClick={() => { setForm(EMPTY_FORM); setShowModal(true) }}>
            + 新增內容
          </button>
        )}
      </div>

      <div className="notes-filter">
        <button className={activeCategory === null ? 'filter-btn active' : 'filter-btn'} onClick={() => setActiveCategory(null)}>全部</button>
        {CATEGORIES.map(c => (
          <button key={c} className={activeCategory === c ? 'filter-btn active' : 'filter-btn'} onClick={() => setActiveCategory(c)}>{c}</button>
        ))}
      </div>

      {error && <div className="chat-error" style={{ margin: '0 24px' }}>{error}</div>}

      <div className="notes-grid">
        {notes.length === 0 ? (
          <div className="notes-empty"><p>目前沒有筆記</p></div>
        ) : (
          notes.map(note => (
            <div key={note.id} className="note-card" onClick={() => setSelectedNote(note)}>
              <div className="note-card-header">
                <span className="note-category">{note.category}</span>
                <span className="note-date">{new Date(note.created_at).toLocaleDateString('zh-TW')}</span>
              </div>
              <NoteCardBody note={note} />
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
                <span className="note-detail-date">{new Date(selectedNote.created_at).toLocaleDateString('zh-TW')}</span>
              </div>
              <button className="modal-close" onClick={() => setSelectedNote(null)}>✕</button>
            </div>
            <NoteDetailBody note={selectedNote} />
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
              <CategoryFields form={editForm} onChange={handleEditChange} />
              <div className="modal-actions">
                <button type="button" className="btn-ghost" onClick={() => setEditNote(null)}>取消</button>
                <button type="submit" className="btn-primary">儲存</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* 刪除確認 */}
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
              <CategoryFields form={form} onChange={handleChange} />
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

import React from 'react'
import { NavLink, useNavigate } from 'react-router-dom'

function Sidebar({ user, onLogout }) {
  const navigate = useNavigate()

  const handleLogout = async () => {
    await onLogout()
    navigate('/')
  }

  return (
    <aside className="sidebar">
      <div className="sidebar-brand">
        <span className="sidebar-brand-icon">⚗️</span>
        <span className="sidebar-brand-name">Rick's Lab</span>
      </div>

      <nav className="sidebar-nav">
        <NavLink to="/" className={({ isActive }) => isActive ? 'nav-item active' : 'nav-item'} end>
          <span className="nav-icon">🏠</span>
          <span>Home</span>
        </NavLink>
        <NavLink to="/dashboard" className={({ isActive }) => isActive ? 'nav-item active' : 'nav-item'}>
          <span className="nav-icon">👤</span>
          <span>Dashboard</span>
        </NavLink>
        <NavLink to="/notes" className={({ isActive }) => isActive ? 'nav-item active' : 'nav-item'}>
          <span className="nav-icon">📝</span>
          <span>Notes</span>
        </NavLink>
        <NavLink to="/ai" className={({ isActive }) => isActive ? 'nav-item active' : 'nav-item'}>
          <span className="nav-icon">🤖</span>
          <span>AI Lab</span>
        </NavLink>
      </nav>

      <div className="sidebar-footer">
        {user ? (
          <>
            <div className="sidebar-user">
              <div className="sidebar-user-avatar">{user.username?.[0]?.toUpperCase()}</div>
              <div className="sidebar-user-info">
                <span className="sidebar-user-name">{user.username}</span>
                <span className="sidebar-user-email">{user.email}</span>
              </div>
            </div>
            <button className="sidebar-logout" onClick={handleLogout}>登出</button>
          </>
        ) : (
          <div className="sidebar-auth-links">
            <NavLink to="/login" className="sidebar-auth-btn primary">登入</NavLink>
            <NavLink to="/register" className="sidebar-auth-btn secondary">註冊</NavLink>
          </div>
        )}
      </div>
    </aside>
  )
}

export default Sidebar

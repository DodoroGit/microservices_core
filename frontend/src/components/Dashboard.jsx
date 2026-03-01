import React, { useState, useEffect } from 'react'
import { userAPI } from '../api'

function Dashboard({ user }) {
  const [users, setUsers] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    fetchUsers()
  }, [])

  const fetchUsers = async () => {
    try {
      const response = await userAPI.getUsers()
      setUsers(response.data || [])
      setLoading(false)
    } catch (err) {
      setError('無法載入用戶列表')
      setLoading(false)
    }
  }

  return (
    <div className="dashboard-page">
      <div className="page-header dark-header">
        <div>
          <h1>Dashboard</h1>
          <p>歡迎回來，{user.username}</p>
        </div>
      </div>

      <div className="dashboard-content">
        <div className="profile-card">
          <div className="profile-avatar">{user.username?.[0]?.toUpperCase()}</div>
          <div className="profile-info">
            <h2>{user.username}</h2>
            <p className="profile-email">{user.email}</p>
            <p className="profile-date">加入時間：{new Date(user.created_at).toLocaleString('zh-TW')}</p>
          </div>
        </div>

        <div className="users-section">
          <div className="section-header">
            <h3>所有用戶</h3>
            <span className="count-badge">{users.length} 人</span>
          </div>
          {loading && <div className="loading-state">載入中...</div>}
          {error && <div className="error-message" style={{ margin: '16px 24px' }}>{error}</div>}
          {!loading && !error && (
            <div className="users-table-wrapper">
              <table className="users-table">
                <thead>
                  <tr>
                    <th>用戶名稱</th>
                    <th>Email</th>
                    <th>建立時間</th>
                  </tr>
                </thead>
                <tbody>
                  {users.length === 0 ? (
                    <tr>
                      <td colSpan={3} className="empty-state">暫無用戶</td>
                    </tr>
                  ) : (
                    users.map((u) => (
                      <tr key={u.id}>
                        <td>
                          <div className="user-cell">
                            <div className="user-avatar-sm">{u.username?.[0]?.toUpperCase()}</div>
                            {u.username}
                          </div>
                        </td>
                        <td>{u.email}</td>
                        <td>{new Date(u.created_at).toLocaleDateString('zh-TW')}</td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default Dashboard

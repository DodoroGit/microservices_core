import React, { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { userAPI } from '../api'

function Register() {
  const [formData, setFormData] = useState({ email: '', username: '', password: '', confirmPassword: '' })
  const [error, setError] = useState('')
  const [success, setSuccess] = useState('')
  const navigate = useNavigate()

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value })
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setSuccess('')

    if (formData.password !== formData.confirmPassword) {
      setError('密碼不一致')
      return
    }
    if (formData.password.length < 6) {
      setError('密碼長度至少需要 6 個字元')
      return
    }

    try {
      await userAPI.register({
        email: formData.email,
        username: formData.username,
        password: formData.password,
      })
      setSuccess('註冊成功！即將跳轉到登入頁面...')
      setTimeout(() => navigate('/login'), 2000)
    } catch (err) {
      setError(err.response?.data?.error || '註冊失敗，請稍後再試')
    }
  }

  return (
    <div className="auth-page">
      <div className="auth-card">
        <h2>建立帳號</h2>
        <p className="auth-subtitle">加入 Rick's Lab</p>
        {error && <div className="error-message" style={{ marginBottom: '16px' }}>{error}</div>}
        {success && <div className="success-message" style={{ marginBottom: '16px' }}>{success}</div>}
        <form className="auth-form" onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Email</label>
            <input
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              required
              placeholder="your@email.com"
            />
          </div>
          <div className="form-group">
            <label>用戶名稱</label>
            <input
              type="text"
              name="username"
              value={formData.username}
              onChange={handleChange}
              required
              placeholder="username"
            />
          </div>
          <div className="form-group">
            <label>密碼</label>
            <input
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              required
              placeholder="至少 6 個字元"
            />
          </div>
          <div className="form-group">
            <label>確認密碼</label>
            <input
              type="password"
              name="confirmPassword"
              value={formData.confirmPassword}
              onChange={handleChange}
              required
              placeholder="再次輸入密碼"
            />
          </div>
          <button type="submit" className="auth-submit">註冊</button>
        </form>
        <div className="auth-footer">
          已有帳號？<Link to="/login">立即登入</Link>
        </div>
      </div>
    </div>
  )
}

export default Register

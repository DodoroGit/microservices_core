import React, { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { userAPI } from '../api'

function Login({ setUser }) {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const navigate = useNavigate()

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    try {
      const response = await userAPI.login({ email, password })
      const userData = response.data.user
      const token = response.data.token
      localStorage.setItem('user', JSON.stringify(userData))
      localStorage.setItem('token', token)
      setUser(userData)
      navigate('/dashboard')
    } catch (err) {
      setError(err.response?.data?.error || '登入失敗，請檢查您的帳號密碼')
    }
  }

  return (
    <div className="auth-page">
      <div className="auth-card">
        <h2>歡迎回來</h2>
        <p className="auth-subtitle">登入 Rick's Lab</p>
        {error && <div className="error-message" style={{ marginBottom: '16px' }}>{error}</div>}
        <form className="auth-form" onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Email</label>
            <input
              type="email"
              value={email}
              onChange={e => setEmail(e.target.value)}
              required
              placeholder="your@email.com"
            />
          </div>
          <div className="form-group">
            <label>密碼</label>
            <input
              type="password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              required
              placeholder="••••••••"
            />
          </div>
          <button type="submit" className="auth-submit">登入</button>
        </form>
        <div className="auth-footer">
          還沒有帳號？<Link to="/register">立即註冊</Link>
        </div>
      </div>
    </div>
  )
}

export default Login

import React, { useState, useEffect } from 'react'
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom'
import './App.css'
import Login from './components/Login.jsx'
import Register from './components/Register.jsx'
import Dashboard from './components/Dashboard.jsx'

function App() {
  const [user, setUser] = useState(null)

  useEffect(() => {
    // Check if user is logged in
    const savedUser = localStorage.getItem('user')
    if (savedUser) {
      setUser(JSON.parse(savedUser))
    }
  }, [])

  const handleLogout = () => {
    localStorage.removeItem('user')
    setUser(null)
  }

  return (
    <Router>
      <div className="App">
        <header className="App-header">
          <nav>
            <h1>Microservices App</h1>
            <div className="nav-links">
              {user ? (
                <>
                  <Link to="/dashboard">Dashboard</Link>
                  <button onClick={handleLogout}>登出</button>
                  <span>歡迎, {user.username}</span>
                </>
              ) : (
                <>
                  <Link to="/login">登入</Link>
                  <Link to="/register">註冊</Link>
                </>
              )}
            </div>
          </nav>
        </header>

        <main className="App-main">
          <Routes>
            <Route path="/" element={
              <div className="home">
                <h2>歡迎來到微服務架構個人網站</h2>
                <p>這是一個使用 React + Golang + PostgreSQL 的微服務架構示例</p>
              </div>
            } />
            <Route path="/login" element={<Login setUser={setUser} />} />
            <Route path="/register" element={<Register />} />
            <Route path="/dashboard" element={
              user ? <Dashboard user={user} /> : <Login setUser={setUser} />
            } />
          </Routes>
        </main>
      </div>
    </Router>
  )
}

export default App

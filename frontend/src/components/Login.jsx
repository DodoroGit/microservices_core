import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { userAPI } from '../api';

function Login({ setUser }) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    try {
      const response = await userAPI.login({ email, password });
      const userData = response.data.user;

      localStorage.setItem('user', JSON.stringify(userData));
      setUser(userData);
      navigate('/dashboard');
    } catch (err) {
      setError(err.response?.data?.error || '登入失敗，請檢查您的帳號密碼');
    }
  };

  return (
    <div className="form-container">
      <h2>登入</h2>
      {error && <div className="error-message">{error}</div>}
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Email</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label>密碼</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <button type="submit" className="submit-btn">登入</button>
      </form>
    </div>
  );
}

export default Login;

import React, { useState, useEffect } from 'react';
import { userAPI } from '../api';

function Dashboard({ user }) {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    try {
      const response = await userAPI.getUsers();
      setUsers(response.data || []);
      setLoading(false);
    } catch (err) {
      setError('無法載入用戶列表');
      setLoading(false);
    }
  };

  return (
    <div className="dashboard">
      <h2>控制台</h2>

      <div className="user-info">
        <h3>個人資訊</h3>
        <p><strong>ID:</strong> {user.id}</p>
        <p><strong>Email:</strong> {user.email}</p>
        <p><strong>用戶名稱:</strong> {user.username}</p>
        <p><strong>建立時間:</strong> {new Date(user.created_at).toLocaleString('zh-TW')}</p>
      </div>

      <div style={{ marginTop: '2rem' }}>
        <h3>所有用戶</h3>
        {loading && <p>載入中...</p>}
        {error && <div className="error-message">{error}</div>}
        {!loading && !error && (
          <div style={{ background: 'white', padding: '1rem', borderRadius: '8px', marginTop: '1rem' }}>
            {users.length === 0 ? (
              <p>暫無用戶</p>
            ) : (
              <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                <thead>
                  <tr style={{ borderBottom: '2px solid #ddd' }}>
                    <th style={{ padding: '0.75rem', textAlign: 'left' }}>用戶名稱</th>
                    <th style={{ padding: '0.75rem', textAlign: 'left' }}>Email</th>
                    <th style={{ padding: '0.75rem', textAlign: 'left' }}>建立時間</th>
                  </tr>
                </thead>
                <tbody>
                  {users.map((u) => (
                    <tr key={u.id} style={{ borderBottom: '1px solid #eee' }}>
                      <td style={{ padding: '0.75rem' }}>{u.username}</td>
                      <td style={{ padding: '0.75rem' }}>{u.email}</td>
                      <td style={{ padding: '0.75rem' }}>
                        {new Date(u.created_at).toLocaleDateString('zh-TW')}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
        )}
      </div>
    </div>
  );
}

export default Dashboard;

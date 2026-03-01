import React from 'react'

function AILab() {
  return (
    <div className="ailab-page">
      <div className="page-header dark-header">
        <div>
          <h1>AI Lab</h1>
          <p>本地 LLM 整合實驗室</p>
        </div>
      </div>

      <div className="coming-soon-section">
        <div className="coming-soon-icon">🤖</div>
        <h2>Ollama 整合開發中</h2>
        <p>即將整合本地 Ollama 模型，讓你在自己的 server 上運行 LLM，對話資料完全掌控在自己手中。</p>
        <div className="ailab-plan">
          <h3>規劃中的功能</h3>
          <ul>
            <li>✦ 對話介面</li>
            <li>✦ 模型選擇（llama3、mistral 等）</li>
            <li>✦ 對話歷史記錄</li>
            <li>✦ 筆記摘要生成</li>
          </ul>
        </div>
      </div>
    </div>
  )
}

export default AILab

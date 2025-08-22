import { useState, useEffect } from 'react';

export default function SimpleTest() {
  const [mounted, setMounted] = useState(false);
  const [time, setTime] = useState('未設置');

  useEffect(() => {
    console.log('🔄 SimpleTest useEffect 執行');
    setMounted(true);
    setTime(new Date().toLocaleTimeString());
    
    const timer = setInterval(() => {
      setTime(new Date().toLocaleTimeString());
    }, 1000);
    
    return () => clearInterval(timer);
  }, []);

  console.log('🎨 SimpleTest 渲染，mounted:', mounted);

  return (
    <div style={{ padding: '20px', fontFamily: 'Arial, sans-serif' }}>
      <h1>Pages Router 測試</h1>
      <p>掛載狀態: <strong style={{ color: mounted ? 'green' : 'red' }}>
        {mounted ? '✅ 已掛載' : '❌ 未掛載'}
      </strong></p>
      <p>時間: <strong>{time}</strong></p>
      <p>環境: <strong>
        {typeof window !== 'undefined' ? '瀏覽器' : 'Node.js'}
      </strong></p>
      <div style={{ marginTop: '20px', padding: '10px', backgroundColor: '#f0f0f0' }}>
        <p><strong>說明：</strong>這是使用 Pages Router 的測試頁面</p>
        <p>如果看到"已掛載"和時間更新，說明 React 正常工作</p>
      </div>
    </div>
  );
}

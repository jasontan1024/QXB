const fetch = require('node-fetch');
(async () => {
  try {
    const response = await fetch('http://localhost:8080/api/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email: 'test3@test.com', password: 'Test123!' })
    });
    const data = await response.json();
    console.log('API 连接成功:', data.success);
  } catch (error) {
    console.log('API 连接失败:', error.message);
  }
})();

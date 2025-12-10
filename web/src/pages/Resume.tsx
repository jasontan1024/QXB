import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../api';
import { marked } from 'marked';
import './Resume.css';

// 配置 marked 选项
marked.setOptions({
  breaks: true, // 支持换行
  gfm: true, // 启用 GitHub Flavored Markdown
});

export default function Resume() {
  const [content, setContent] = useState<string>('');
  const [html, setHtml] = useState<string>('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    const load = async () => {
      try {
        setLoading(true);
        const res = await api.getResume();
        if (res.success && res.data) {
          const markdownContent = res.data.content || '';
          setContent(markdownContent);
          // marked.parse 返回 Promise<string>
          const htmlContent = await marked.parse(markdownContent);
          setHtml(htmlContent);
        } else {
          setError(res.error || '加载失败');
        }
      } catch (e: any) {
        setError(e?.message || '加载失败');
      } finally {
        setLoading(false);
      }
    };
    load();
  }, []);

  if (loading) {
    return <div className="resume-page"><div className="loading">加载中...</div></div>;
  }

  if (error) {
    return <div className="resume-page"><div className="error">{error}</div></div>;
  }

  return (
    <div className="resume-page">
      <div style={{ marginBottom: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1>作者简历</h1>
        <Link to="/dashboard" style={{ fontSize: '14px', color: '#007bff', textDecoration: 'none' }}>
          返回 Dashboard
        </Link>
      </div>
      <div className="resume-render">
        <div className="markdown-body" dangerouslySetInnerHTML={{ __html: html }} />
      </div>
    </div>
  );
}


import React, { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { api } from '../api';
import { auth } from '../auth';
import './Dashboard.css';

export default function Dashboard() {
  const [user, setUser] = useState<{ email: string; address: string } | null>(null);
  const [balance, setBalance] = useState<string>('0');
  const [rewardStatus, setRewardStatus] = useState<{ canClaim: boolean; lastClaimDay: number } | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [claimLoading, setClaimLoading] = useState(false);
  const [transferLoading, setTransferLoading] = useState(false);
  const [claimPassword, setClaimPassword] = useState('');
  const [transferPassword, setTransferPassword] = useState('');
  const [showPasswordInput, setShowPasswordInput] = useState(false);
  const [transferTo, setTransferTo] = useState('');
  const [transferAmount, setTransferAmount] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    if (!auth.isAuthenticated()) {
      navigate('/login');
      return;
    }
    loadData();
  }, [navigate]);

  const loadData = async () => {
    try {
      setLoading(true);
      const userRes = await api.getMe();
      
      // 如果获取用户信息失败（token 无效），清除 token 并重定向到登录页
      if (!userRes.success || !userRes.data) {
        auth.removeToken();
        navigate('/login');
        return;
      }

      const user = userRes.data;
      setUser(user);

      // 获取余额和奖励状态
      const [balanceRes, rewardRes] = await Promise.all([
        api.getBalance(user.address),
        api.getRewardStatus(user.address),
      ]);

      if (balanceRes?.success && balanceRes.data) {
        setBalance(balanceRes.data.balance);
      }

      if (rewardRes?.success && rewardRes.data) {
        setRewardStatus({
          canClaim: rewardRes.data.canClaim,
          lastClaimDay: rewardRes.data.lastClaimDay,
        });
      }
    } catch (err: any) {
      // 如果是认证错误，清除 token 并重定向
      if (err.message?.includes('401') || err.message?.includes('Unauthorized')) {
        auth.removeToken();
        navigate('/login');
        return;
      }
      setError(err.message || '加载数据失败');
    } finally {
      setLoading(false);
    }
  };

  const handleClaim = async () => {
    if (!user) return;

    if (!showPasswordInput) {
      setShowPasswordInput(true);
      return;
    }

    if (!claimPassword) {
      setError('请输入密码');
      return;
    }

    setClaimLoading(true);
    setError('');
    try {
      const response = await api.claimReward({ password: claimPassword });
      if (response.success && response.data) {
        alert(`奖励领取成功！交易哈希: ${response.data.txHash}`);
        setClaimPassword('');
        setShowPasswordInput(false);
        await loadData();
      }
    } catch (err: any) {
      setError(err.message || '领取失败');
    } finally {
      setClaimLoading(false);
    }
  };

  const handleTransfer = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user || !transferTo || !transferAmount || !transferPassword) {
      setError('请填写完整信息');
      return;
    }

    // 将金额转换为 wei（18 位小数）
    const amountInWei = (parseFloat(transferAmount) * 1e18).toString();
    if (isNaN(parseFloat(transferAmount)) || parseFloat(transferAmount) <= 0) {
      setError('请输入有效的金额');
      return;
    }

    setTransferLoading(true);
    setError('');
    try {
      const response = await api.transfer({
        to: transferTo,
        amount: amountInWei,
        password: transferPassword,
      });
      if (response.success && response.data) {
        alert(`转账成功！交易哈希: ${response.data.txHash}`);
        setTransferTo('');
        setTransferAmount('');
        setTransferPassword('');
        await loadData();
      }
    } catch (err: any) {
      setError(err.message || '转账失败');
    } finally {
      setTransferLoading(false);
    }
  };

  const handleLogout = () => {
    auth.removeToken();
    navigate('/login');
  };

  if (loading) {
    return (
      <div className="dashboard-container">
        <div className="loading">加载中...</div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="dashboard-container">
        <div className="error">无法加载用户信息</div>
      </div>
    );
  }

  return (
    <div className="dashboard-container">
      <div className="dashboard-header">
        <h1>QXB 代币管理</h1>
        <div className="user-info">
          <span>{user.email}</span>
          <Link to="/resume" className="author-link">作者简历</Link>
          <button onClick={handleLogout} className="logout-button">
            退出
          </button>
        </div>
      </div>

      <div className="dashboard-content">
        <div className="info-card">
          <h2>钱包地址</h2>
          <p className="address">{user.address}</p>
        </div>

        <div className="info-card">
          <h2>代币余额</h2>
          <p className="balance">{balance} QXB</p>
        </div>

        <div className="info-card">
          <h2>每日奖励</h2>
          {rewardStatus ? (
            <div>
              <p className={rewardStatus.canClaim ? 'can-claim' : 'cannot-claim'}>
                {rewardStatus.canClaim ? '✅ 可以领取' : '⏳ 今日已领取'}
              </p>
              {rewardStatus.lastClaimDay > 0 && (
                <p className="claim-info">上次领取日期: {rewardStatus.lastClaimDay}</p>
              )}
            </div>
          ) : (
            <p>加载中...</p>
          )}
          {showPasswordInput && (
            <div className="password-input-group">
              <input
                type="password"
                value={claimPassword}
                onChange={(e) => setClaimPassword(e.target.value)}
                placeholder="请输入密码"
                className="password-input"
              />
            </div>
          )}
          <button
            onClick={handleClaim}
            disabled={claimLoading || (showPasswordInput && !claimPassword)}
            className="claim-button"
          >
            {claimLoading ? '领取中...' : showPasswordInput ? '确认领取' : '领取奖励'}
          </button>
        </div>

        <div className="info-card">
          <h2>转账</h2>
          <form onSubmit={handleTransfer}>
            <div className="form-group">
              <label>接收地址</label>
              <input
                type="text"
                value={transferTo}
                onChange={(e) => setTransferTo(e.target.value)}
                placeholder="0x..."
                required
              />
            </div>
            <div className="form-group">
              <label>金额</label>
              <input
                type="text"
                value={transferAmount}
                onChange={(e) => setTransferAmount(e.target.value)}
                placeholder="例如: 100"
                required
              />
            </div>
            <div className="form-group">
              <label>密码</label>
              <input
                type="password"
                value={transferPassword}
                onChange={(e) => setTransferPassword(e.target.value)}
                placeholder="请输入密码"
                required
              />
            </div>
            {error && <div className="error-message">{error}</div>}
            <button type="submit" disabled={transferLoading} className="transfer-button">
              {transferLoading ? '转账中...' : '转账'}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}


// 直接使用后端 URL（在测试环境中）
// 注意：在生产环境中应该使用代理或配置 CORS
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterResponse {
  user_id: number;
  email: string;
  address: string;
  token: string;
}

export interface LoginResponse {
  user_id: number;
  email: string;
  address: string;
  token: string;
}

export interface UserInfo {
  user_id: number;
  email: string;
  address: string;
}

export interface TokenInfo {
  name: string;
  symbol: string;
  decimals: number;
  totalSupply: string;
  version?: string;
}

export interface BalanceInfo {
  address: string;
  balance: string;
  symbol: string;
}

export interface RewardStatus {
  address: string;
  canClaim: boolean;
  lastClaimDay: number;
  nextClaimDay: number;
}

export interface ClaimRequest {
  privateKey?: string;
  password?: string;
}

export interface TransferRequest {
  to: string;
  amount: string;
  password: string;
}

export interface ClaimResponse {
  txHash: string;
  status: string;
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  const token = localStorage.getItem('token');
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string> || {}),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    let errorMessage = '请求失败';
    try {
      const errorData = await response.json();
      errorMessage = errorData.error || errorMessage;
    } catch (e) {
      const text = await response.text();
      errorMessage = text || errorMessage;
    }
    throw new Error(errorMessage);
  }
  
  const data = await response.json();
  return data;
}

export const api = {
  // 认证相关
  async register(req: RegisterRequest): Promise<ApiResponse<RegisterResponse>> {
    return request<RegisterResponse>('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify(req),
    });
  },

  async login(req: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    return request<LoginResponse>('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify(req),
    });
  },

  async getMe(): Promise<ApiResponse<UserInfo>> {
    return request<UserInfo>('/api/auth/me');
  },

  // 代币相关
  async getTokenInfo(): Promise<ApiResponse<TokenInfo>> {
    return request<TokenInfo>('/api/token/info');
  },

  async getBalance(address: string): Promise<ApiResponse<BalanceInfo>> {
    return request<BalanceInfo>(`/api/token/balance/${address}`);
  },

  async transfer(req: TransferRequest): Promise<ApiResponse<ClaimResponse>> {
    return request<ClaimResponse>('/api/token/transfer', {
      method: 'POST',
      body: JSON.stringify(req),
    });
  },

  // 奖励相关
  async getRewardStatus(address: string): Promise<ApiResponse<RewardStatus>> {
    return request<RewardStatus>(`/api/reward/status/${address}`);
  },

  async claimReward(req: ClaimRequest): Promise<ApiResponse<ClaimResponse>> {
    return request<ClaimResponse>('/api/reward/claim', {
      method: 'POST',
      body: JSON.stringify(req),
    });
  },
};


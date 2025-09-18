export interface User {
  id: number;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
  expires_at: number;
}

export interface RegisterResponse {
  message: string;
  user: User;
}

export interface ApiError {
  error: string;
}
import api from './api';
import type { User, LoginResponse, RegisterResponse } from '../types/auth';
import type { LoginFormData } from '../types/validation';

export const authAPI = {
  login: async (credentials: LoginFormData): Promise<LoginResponse> => {
    const response = await api.post('/auth/login', credentials);
    return response.data;
  },

  register: async (userData: { email: string; password: string }): Promise<RegisterResponse> => {
    const response = await api.post('/auth/register', userData);
    return response.data;
  },

  getCurrentUser: async (): Promise<User> => {
    const response = await api.get('/auth/me');
    return response.data;
  },

  logout: async (): Promise<void> => {
    await api.post('/auth/logout');
  },
};
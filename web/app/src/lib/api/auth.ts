import * as request from '@/lib/api'

import { UserModel } from '@/lib/models'

export const whoami = async (): Promise<UserModel> => {
  type WhoamiResponse = {
    success: boolean
    user: UserModel
  }

  const { body } = await request.GET<WhoamiResponse>('/api/v1/auth/whoami')
  return body.user
}

export const login = async (username: string, password: string): Promise<void> => {
  type LoginRequest = {
    username: string
    password: string
  }

  type LoginResponse = {
    success: boolean
    token: string
  }

  try {
    const { body } = await request.POST<LoginRequest, LoginResponse>(
      '/api/v1/auth/login',
      { username, password },
    )
    localStorage.setItem('token', body.token)
  }
  catch (error) {
    localStorage.removeItem('token')
    throw error
  }
}

export const logout = async (): Promise<void> => {
  localStorage.removeItem('token')
}

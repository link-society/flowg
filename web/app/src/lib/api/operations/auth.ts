import * as request from '@/lib/api/request'

import { ProfileModel, UserModel, PermissionsModel } from '@/lib/models'

export const whoami = async (): Promise<ProfileModel> => {
  type WhoamiResponse = {
    success: boolean
    user: UserModel
    permissions: PermissionsModel
  }

  const { body } = await request.GET<WhoamiResponse>('/api/v1/auth/whoami')
  return { user: body.user, permissions: body.permissions }
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

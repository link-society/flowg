import * as request from '@/lib/api/request'
import { PermissionsModel, ProfileModel, UserModel } from '@/lib/models/auth'

export const whoami = async (): Promise<ProfileModel> => {
  type WhoamiResponse = {
    success: boolean
    user: UserModel
    permissions: PermissionsModel
  }

  const { body } = await request.GET<WhoamiResponse>({
    path: '/api/v1/auth/whoami',
  })
  return { user: body.user, permissions: body.permissions }
}

export const login = async (
  username: string,
  password: string
): Promise<void> => {
  type LoginRequest = {
    username: string
    password: string
  }

  type LoginResponse = {
    success: boolean
    token: string
  }

  try {
    const { body } = await request.POST<LoginRequest, LoginResponse>({
      path: '/api/v1/auth/login',
      body: { username, password },
    })
    localStorage.setItem('token', body.token)
  } catch (error) {
    localStorage.removeItem('token')
    throw error
  }
}

export const logout = async (): Promise<void> => {
  localStorage.removeItem('token')
}

export const changePassword = async (
  oldPassword: string,
  newPassword: string
): Promise<void> => {
  type ChangePasswordRequest = {
    old_password: string
    new_password: string
  }

  type ChangePasswordResponse = {
    success: boolean
  }

  await request.POST<ChangePasswordRequest, ChangePasswordResponse>({
    path: '/api/v1/auth/change-password',
    body: {
      old_password: oldPassword,
      new_password: newPassword,
    },
  })
}

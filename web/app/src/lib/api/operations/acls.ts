import * as request from '@/lib/api/request'

import { UserModel, RoleModel } from '@/lib/models/auth'

export const listUsers = async (): Promise<UserModel[]> => {
  type ListUsersResponse = {
    success: boolean
    users: UserModel[]
  }

  const { body } = await request.GET<ListUsersResponse>({
    path: '/api/v1/users',
  })
  return body.users
}

export const saveUser = async (user: UserModel, password: string): Promise<void> => {
  type SaveUserRequest = {
    password: string
    roles: string[]
  }

  type SaveUserResponse = {
    success: boolean
  }

  await request.PUT<SaveUserRequest, SaveUserResponse>({
    path: `/api/v1/users/${user.name}`,
    body: {
      password,
      roles: user.roles,
    },
  })
}

export const deleteUser = async (username: string): Promise<void> => {
  await request.DELETE({
    path: `/api/v1/users/${username}`,
  })
}

export const listRoles = async (): Promise<RoleModel[]> => {
  type ListRolesResponse = {
    success: boolean
    roles: RoleModel[]
  }

  const { body } = await request.GET<ListRolesResponse>({
    path: '/api/v1/roles',
  })
  return body.roles
}

export const saveRole = async (role: RoleModel): Promise<void> => {
  type SaveRoleRequest = {
    scopes: string[]
  }

  type SaveRoleResponse = {
    success: boolean
  }

  await request.PUT<SaveRoleRequest, SaveRoleResponse>({
    path: `/api/v1/roles/${role.name}`,
    body: {
      scopes: role.scopes,
    },
  })
}

export const deleteRole = async (roleName: string): Promise<void> => {
  await request.DELETE({
    path: `/api/v1/roles/${roleName}`,
  })
}

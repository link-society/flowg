import { LoaderFunction } from 'react-router-dom'

import { loginRequired } from '@/lib/decorators/loaders'
import * as aclApi from '@/lib/api/operations/acls'

import { RoleModel, UserModel } from '@/lib/models'

export type LoaderData = {
  roles: RoleModel[]
  users: UserModel[]
}

export const loader: LoaderFunction = loginRequired(
  async (): Promise<LoaderData> => {
    const [roles, users] = await Promise.all([
      aclApi.listRoles(),
      aclApi.listUsers(),
    ])

    return { roles, users }
  },
)
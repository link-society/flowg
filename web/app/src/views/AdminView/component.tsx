import { LoaderFunction, useLoaderData } from 'react-router'

import * as aclApi from '@/lib/api/operations/acls'
import { getSystemConfiguration } from '@/lib/api/operations/config'

import { loginRequired } from '@/lib/decorators/loaders'

import RoleTable from '@/components/RoleTable/component'
import UserTable from '@/components/UserTable/component'

import { AdminViewContainer, AdminViewPanel } from './styles'
import { LoaderData } from './types'

export const loader: LoaderFunction = loginRequired(
  async (): Promise<LoaderData> => {
    const [roles, users, systemConfig] = await Promise.all([
      aclApi.listRoles(),
      aclApi.listUsers(),
      getSystemConfiguration(),
    ])

    return { roles, users, defaultRoles: systemConfig.default_roles ?? [] }
  }
)

const AdminView = () => {
  const { roles, users, defaultRoles } = useLoaderData() as LoaderData

  return (
    <AdminViewContainer variant="page">
      <AdminViewPanel>
        <RoleTable roles={roles} />
      </AdminViewPanel>
      <AdminViewPanel>
        <UserTable roles={roles.map((role) => role.name)} users={users} defaultRoles={defaultRoles} />
      </AdminViewPanel>
    </AdminViewContainer>
  )
}

export default AdminView

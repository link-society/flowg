import { LoaderFunction, useLoaderData } from 'react-router'

import Grid from '@mui/material/Grid'

import * as aclApi from '@/lib/api/operations/acls'

import RoleModel from '@/lib/models/RoleModel'
import UserModel from '@/lib/models/UserModel'

import { loginRequired } from '@/lib/decorators/loaders'

import RoleTable from '@/components/RoleTable'
import UserTable from '@/components/UserTable'

type LoaderData = {
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
  }
)

const AdminView = () => {
  const { roles, users } = useLoaderData() as LoaderData

  return (
    <Grid container spacing={2} className="p-3 h-full max-lg:overflow-auto">
      <Grid size={{ xs: 12, md: 6 }} className="lg:h-full">
        <RoleTable roles={roles} />
      </Grid>
      <Grid size={{ xs: 12, md: 6 }} className="lg:h-full">
        <UserTable roles={roles.map((role) => role.name)} users={users} />
      </Grid>
    </Grid>
  )
}

export default AdminView

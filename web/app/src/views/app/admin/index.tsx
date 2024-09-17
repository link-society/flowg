import { useLoaderData } from 'react-router-dom'

import Grid from '@mui/material/Grid2'

import { RoleList } from './role-list'
import { UserList } from './user-list'

import { loginRequired } from '@/lib/decorators/loaders'
import * as aclApi from '@/lib/api/operations/acls'

import { RoleModel, UserModel } from '@/lib/models'

export const loader = async () => {
  const [roles, users] = await Promise.all([
    loginRequired(aclApi.listRoles)(),
    loginRequired(aclApi.listUsers)(),
  ])

  return { roles, users }
}

export const AdminView = () => {
  const { roles, users } = useLoaderData() as { roles: RoleModel[], users: UserModel[] }

  return (
    <>
      <Grid
        container
        spacing={2}
        className="p-3 h-full max-lg:overflow-auto"
      >
        <Grid size={{ xs: 12, md: 6 }} className="lg:h-full">
          <RoleList roles={roles} />
        </Grid>
        <Grid size={{ xs: 12, md: 6 }} className="lg:h-full">
          <UserList
            roles={roles.map((role) => role.name)}
            users={users}
          />
        </Grid>
      </Grid>
    </>
  )
}

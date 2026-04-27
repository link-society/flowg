import RoleModel from '@/lib/models/RoleModel'

export type ButtonNewRoleProps = Readonly<{
  onRoleCreated: (role: RoleModel) => void
}>

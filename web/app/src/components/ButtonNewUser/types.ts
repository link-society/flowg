import UserModel from '@/lib/models/UserModel'

export type ButtonNewUserProps = Readonly<{
  roles: string[]
  defaultRoles: string[]
  onUserCreated: (user: UserModel) => void
}>

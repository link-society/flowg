import UserModel from '@/lib/models/UserModel'

export type UserTableProps = Readonly<{
  roles: string[]
  users: UserModel[]
}>

import RoleModel from '@/lib/models/RoleModel'
import UserModel from '@/lib/models/UserModel'

export type LoaderData = {
  roles: RoleModel[]
  users: UserModel[]
}

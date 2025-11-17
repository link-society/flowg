import PermissionsModel from '@/lib/models/PermissionsModel'
import UserModel from '@/lib/models/UserModel'

type ProfileModel = {
  user: UserModel
  permissions: PermissionsModel
}

export default ProfileModel

import Card from '@mui/material/Card'
import CardContent from '@mui/material/CardContent'
import CardHeader from '@mui/material/CardHeader'
import Divider from '@mui/material/Divider'

import AccountCircleIcon from '@mui/icons-material/AccountCircle'

import PasswordChange from '@/components/PasswordChange'
import PermissionDisplay from '@/components/PermissionDisplay'
import RolesDisplay from '@/components/RolesDisplay'
import UsernameDisplay from '@/components/UsernameDisplay'

const ProfileInfo = () => (
  <Card className="lg:h-full lg:flex lg:flex-col lg:items-stretch">
    <CardHeader
      title={
        <div className="flex items-center gap-3">
          <AccountCircleIcon />
          <span>Account Information</span>
        </div>
      }
      className="bg-blue-400 text-white shadow-lg"
    />
    <CardContent
      className="
        lg:grow lg:shrink lg:h-0 lg:overflow-auto
        flex flex-col gap-3 items-stretch
      "
    >
      <UsernameDisplay />
      <RolesDisplay />
      <PermissionDisplay />
      <Divider />
      <PasswordChange />
    </CardContent>
  </Card>
)

export default ProfileInfo

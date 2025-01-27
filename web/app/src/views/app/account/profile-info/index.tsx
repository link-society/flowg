import AccountCircleIcon from '@mui/icons-material/AccountCircle'

import Card from '@mui/material/Card'
import CardHeader from '@mui/material/CardHeader'
import CardContent from '@mui/material/CardContent'
import Divider from '@mui/material/Divider'

import { Username } from './username'
import { RoleList } from './role-list'
import { Permissions } from './permissions'
import { PasswordChange } from './password-change'

export const ProfileInfo = () => (
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
      <Username />
      <RoleList />
      <Permissions />
      <Divider />
      <PasswordChange />
    </CardContent>
  </Card>
)

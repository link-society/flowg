import { useTranslation } from 'react-i18next'

import Divider from '@mui/material/Divider'
import Typography from '@mui/material/Typography'

import AccountCircleIcon from '@mui/icons-material/AccountCircle'

import PasswordChange from '@/components/PasswordChange/component'
import PermissionDisplay from '@/components/PermissionDisplay/component'
import RolesDisplay from '@/components/RolesDisplay/component'
import UsernameDisplay from '@/components/UsernameDisplay/component'

import {
  ProfileInfoCard,
  ProfileInfoCardContent,
  ProfileInfoCardHeader,
  ProfileInfoCardHeaderTitle,
} from './styles'

const ProfileInfo = () => {
  const { t } = useTranslation()

  return (
    <ProfileInfoCard>
      <ProfileInfoCardHeader
        title={
          <ProfileInfoCardHeaderTitle>
            <AccountCircleIcon />
            <Typography variant="titleSm">
              {t('components.profileInfo.title')}
            </Typography>
          </ProfileInfoCardHeaderTitle>
        }
      />
      <ProfileInfoCardContent>
        <UsernameDisplay />
        <RolesDisplay />
        <PermissionDisplay />
        <Divider />
        <PasswordChange />
      </ProfileInfoCardContent>
    </ProfileInfoCard>
  )
}

export default ProfileInfo

import Box from '@mui/material/Box'
import Chip from '@mui/material/Chip'
import ListItem from '@mui/material/ListItem'

import { useProfile } from '@/lib/hooks/profile'

import { Label, RolesPaper } from './styles'

const RolesDisplay = () => {
  const { user } = useProfile()

  return (
    <Box>
      <Label variant="text">Roles:</Label>

      <RolesPaper variant="outlined" component="ul">
        {user.roles.map((role) => (
          <ListItem key={role}>
            <Chip label={role} size="small" />
          </ListItem>
        ))}
      </RolesPaper>
    </Box>
  )
}

export default RolesDisplay

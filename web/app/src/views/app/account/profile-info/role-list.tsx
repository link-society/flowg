import Paper from '@mui/material/Paper'
import Chip from '@mui/material/Chip'
import ListItem from '@mui/material/ListItem'

import { useProfile } from '@/lib/context/profile'

export const RoleList = () => {
  const { user } = useProfile()

  return (
    <div>
      <span className="font-semibold mb-1">
        Roles:
      </span>

      <Paper
        className="
          p-1
          flex flex-row justify-center flex-wrap
          list-none
        "
        variant="outlined"
        component="ul"
      >
        {user.roles.map((role) => (
          <ListItem key={role}>
            <Chip label={role} size="small" />
          </ListItem>
        ))}
      </Paper>
    </div>
  )
}

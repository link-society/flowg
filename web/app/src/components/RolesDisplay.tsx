import Chip from '@mui/material/Chip'
import ListItem from '@mui/material/ListItem'
import Paper from '@mui/material/Paper'

import { useProfile } from '@/lib/hooks/profile'

const RolesDisplay = () => {
  const { user } = useProfile()

  return (
    <div>
      <span className="font-semibold mb-1">Roles:</span>

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

export default RolesDisplay

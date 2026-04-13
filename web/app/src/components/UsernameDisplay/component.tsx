import TextField from '@mui/material/TextField'

import { useProfile } from '@/lib/hooks/profile'

const UsernameDisplay = () => {
  const { user } = useProfile()

  return (
    <TextField
      label="Username"
      value={user.name}
      type="text"
      variant="outlined"
      slotProps={{
        input: {
          readOnly: true,
        },
      }}
    />
  )
}

export default UsernameDisplay

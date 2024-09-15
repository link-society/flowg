import TextField from '@mui/material/TextField'

import { useProfile } from '@/lib/context/profile'

export const Username = () => {
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

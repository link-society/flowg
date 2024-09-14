import { useState } from 'react'
import { Form, useNavigate } from 'react-router-dom'
import { useSnackbar } from 'notistack'

import { AccountCircle, Lock } from '@mui/icons-material'

import {
  Box,
  Button,
  Card,
  CircularProgress,
  Grid2,
  TextField,
} from '@mui/material'

import * as api from '@/lib/api'
import * as authApi from '@/lib/api/auth'

export default function LoginView() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)

  const navigate = useNavigate()
  const { enqueueSnackbar } = useSnackbar()

  const handleLogin = async () => {
    setLoading(true)

    try {
      await authApi.login(username, password)
      navigate('/web/')
    }
    catch (error) {
      if (error instanceof api.UnauthenticatedError) {
        enqueueSnackbar({
          message: 'Invalid credentials',
          variant: 'error'
        })
      }
      else {
        enqueueSnackbar({
          message: 'Unknown error',
          variant: 'error'
        })
      }

      console.error(error)
    }

    setLoading(false)
  }

  return (
    <div className="py-6">
      <Grid2 container>
        <Grid2 size={{ sm: 12, md: 6 }} offset={{ sm: 0, md: 3 }}>
          <Card className="p-3">
            <Form
              className="flex flex-col items-stretch gap-3"
              onSubmit={e => {
                e.preventDefault()
                handleLogin()
              }}
            >
              <header>
                <h1 className="text-2xl text-center">Sign In with your account</h1>
              </header>

              <hr />

              <section className="flex flex-col items-stretch gap-3">
                <Box className="flex flex-row items-end">
                  <AccountCircle sx={{ color: 'action.active', mr: 1, my: 0.5 }} />
                  <TextField
                    label="Username"
                    value={username}
                    type="text"
                    onChange={e => setUsername(e.target.value)}
                    variant="standard"
                    className="flex-grow"
                  />
                </Box>

                <Box className="flex flex-row items-end">
                  <Lock sx={{ color: 'action.active', mr: 1, my: 0.5 }} />
                  <TextField
                    label="Password"
                    value={password}
                    type="password"
                    onChange={e => setPassword(e.target.value)}
                    variant="standard"
                    className="flex-grow"
                  />
                </Box>
              </section>

              <hr />

              <Button
                variant="contained"
                color="secondary"
                className="w-full"
                type="submit"
              >
                {loading
                  ? <CircularProgress color="inherit" />
                  : <>Sign In</>
                }
              </Button>
            </Form>
          </Card>
        </Grid2>
      </Grid2>
    </div>
  )
}

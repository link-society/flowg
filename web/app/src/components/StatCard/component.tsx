import { Button, Card, CardActions, CardHeader, Divider } from '@mui/material'

import { useNavigate } from 'react-router'

import { StatCardContent, StatCardHeaderWrapper, StatCardValue } from './styles'
import { StatCardProps } from './types'

const StatCard = ({ icon, title, value, to }: StatCardProps) => {
  const navigate = useNavigate()

  return (
    <Card>
      <CardHeader
        title={
          <StatCardHeaderWrapper>
            {icon}
            {title}
          </StatCardHeaderWrapper>
        }
      />

      <StatCardContent>
        <StatCardValue variant="titleLg">{value}</StatCardValue>
        <Divider />
      </StatCardContent>

      <CardActions>
        <Button onClick={() => navigate(to)} fullWidth>
          View More
        </Button>
      </CardActions>
    </Card>
  )
}

export default StatCard

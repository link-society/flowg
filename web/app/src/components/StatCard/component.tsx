import {
  Button,
  Card,
  CardActions,
  CardHeader,
  Divider,
  Typography,
} from '@mui/material'

import { StatCardContent, StatCardHeaderWrapper } from './styles'
import { StatCardProps } from './types'

const StatCard = ({ icon, title, value, href }: StatCardProps) => (
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
      <Typography variant="titleLg" fontWeight={700}>
        {value}
      </Typography>
      <Divider />
    </StatCardContent>

    <CardActions>
      <Button href={href} fullWidth>
        View More
      </Button>
    </CardActions>
  </Card>
)

export default StatCard

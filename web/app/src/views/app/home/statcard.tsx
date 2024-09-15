import { ReactNode } from 'react'

import Divider from '@mui/material/Divider'
import Card from '@mui/material/Card'
import CardHeader from '@mui/material/CardHeader'
import CardContent from '@mui/material/CardContent'
import CardActions from '@mui/material/CardActions'
import Button from '@mui/material/Button'

type StatCardProps = {
  icon: ReactNode
  title: ReactNode
  value: ReactNode
  href: string
}

export const StatCard = ({ icon, title, value, href }: StatCardProps) => (
  <Card>
    <CardHeader
      title={
        <div
          className="
            flex items-center justify-center gap-3
            text-2xl font-semibold
          "
        >
          {icon}
          {title}
        </div>
      }
    />
    <CardContent className="!p-0 text-center text-3xl font-bold">
      <div className="mb-3">{value}</div>
      <Divider />
    </CardContent>
    <CardActions>
      <Button href={href} className="w-full">
        View More
      </Button>
    </CardActions>
  </Card>
)

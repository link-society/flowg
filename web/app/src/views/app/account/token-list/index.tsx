import { useState } from 'react'

import KeyIcon from '@mui/icons-material/Key'
import AddIcon from '@mui/icons-material/Add'

import Card from '@mui/material/Card'
import CardHeader from '@mui/material/CardHeader'
import CardContent from '@mui/material/CardContent'
import { Button } from '@mui/material'

type TokenListProps = {
  tokens: string[]
}

export const TokenList = (props: TokenListProps) => {
  const [tokens] = useState(props.tokens)

  return (
    <Card className="lg:h-full lg:flex lg:flex-col lg:items-stretch">
      <CardHeader
        title={
          <div className="flex items-center gap-3">
            <KeyIcon />
            <span className="flex-grow">API Tokens</span>
            <Button
              variant="contained"
              color="secondary"
              size="small"
              startIcon={<AddIcon />}
            >
              New Token
            </Button>
          </div>
        }
        className="bg-blue-400 text-white shadow-lg"
      />
      <CardContent className="!p-0 lg:flex-grow lg:flex-shrink lg:h-0 lg:overflow-auto">
        {tokens}
      </CardContent>
    </Card>
  )
}

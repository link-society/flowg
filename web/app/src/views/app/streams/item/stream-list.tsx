import List from '@mui/material/List'
import ListItemButton from '@mui/material/ListItemButton'
import ListItemText from '@mui/material/ListItemText'

type StreamListProps = {
  streams: string[]
  currentStream: string
}

export const StreamList = (props: StreamListProps) => {
  return (
    <List component="nav" className="!p-0">
      {props.streams.map((stream, index) => (
        <ListItemButton
          key={index}
          component="a"
          href={`/web/streams/${stream}`}
          sx={stream !== props.currentStream
            ? {
              color: 'secondary.main',
            }
            : {
              backgroundColor: 'secondary.main',
              '&:hover': {
                backgroundColor: 'secondary.main',
              },
              color: 'white',
            }
          }
        >
          <ListItemText primary={stream} />
        </ListItemButton>
      ))}
    </List>
  )
}

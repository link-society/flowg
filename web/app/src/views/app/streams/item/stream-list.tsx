import List from '@mui/material/List'
import ListItemButton from '@mui/material/ListItemButton'
import ListItemText from '@mui/material/ListItemText'

type StreamListProps = Readonly<{
  streams: string[]
  currentStream: string
}>

export const StreamList = (props: StreamListProps) => {
  return (
    <List component="nav" className="p-0!">
      {props.streams.map((stream) => (
        <ListItemButton
          key={stream}
          component="a"
          href={`/web/streams/${stream}`}
          sx={
            stream === props.currentStream
              ? {
                  backgroundColor: 'secondary.main',
                  '&:hover': {
                    backgroundColor: 'secondary.main',
                  },
                  color: 'white',
                }
              : {
                  color: 'secondary.main',
                }
          }
        >
          <ListItemText primary={stream} />
        </ListItemButton>
      ))}
    </List>
  )
}

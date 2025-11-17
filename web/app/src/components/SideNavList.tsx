import List from '@mui/material/List'
import ListItemButton from '@mui/material/ListItemButton'
import ListItemText from '@mui/material/ListItemText'
import Paper from '@mui/material/Paper'

type SideNavListProps = Readonly<{
  namespace: string
  urlPrefix: string
  items: string[]
  currentItem: string
}>

const SideNavList = (props: SideNavListProps) => {
  return (
    <Paper className="h-full overflow-auto">
      <List component="nav" className="p-0!">
        {props.items.map((item) => (
          <ListItemButton
            key={item}
            component="a"
            href={`${props.urlPrefix}/${item}`}
            sx={
              item === props.currentItem
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
            <ListItemText
              id={`label:${props.namespace}.list-item.${item}`}
              primary={item}
            />
          </ListItemButton>
        ))}
      </List>
    </Paper>
  )
}

export default SideNavList

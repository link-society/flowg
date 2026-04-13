import ListItemButton from '@mui/material/ListItemButton'
import ListItemText from '@mui/material/ListItemText'

import { SideNavListContainer, SideNavListNav } from './styles'
import { SideNavListProps } from './types'

const SideNavList = (props: SideNavListProps) => {
  return (
    <SideNavListContainer>
      <SideNavListNav>
        {props.items.map((item) => (
          <ListItemButton
            key={item}
            component="a"
            href={`${props.urlPrefix}/${item}`}
            className={item === props.currentItem ? 'active' : undefined}
          >
            <ListItemText
              id={`label:${props.namespace}.list-item.${item}`}
              primary={item}
            />
          </ListItemButton>
        ))}
      </SideNavListNav>
    </SideNavListContainer>
  )
}

export default SideNavList

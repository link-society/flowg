import { useNavigate } from 'react-router'

import ListItemButton from '@mui/material/ListItemButton'
import ListItemText from '@mui/material/ListItemText'

import { SideNavListContainer, SideNavListNav } from './styles'
import { SideNavListProps } from './types'

const SideNavList = (props: SideNavListProps) => {
  const navigate = useNavigate()

  return (
    <SideNavListContainer>
      <SideNavListNav>
        {props.items.map((item) => (
          <ListItemButton
            key={item}
            onClick={() => navigate(`${props.urlPrefix}/${item}`)}
            className={item === props.currentItem ? 'active' : undefined}
          >
            <ListItemText
              id={`label:${props.namespace}.list-item.${item}`}
              primary={item}
              className={item === props.currentItem ? 'active' : undefined}
            />
          </ListItemButton>
        ))}
      </SideNavListNav>
    </SideNavListContainer>
  )
}

export default SideNavList

import LinkSocietyLogoUrl from '@site/static/img/users/link-society.png'
import SaashupLogoUrl from '@site/static/img/users/saashup.svg'
import StatStreamLogoUrl from '@site/static/img/users/statstream.png'

import styles from './styles.module.css'

export default [
  {
    key: 'saashup',
    image: <SaashupLogoUrl />,
    href: 'https://saashup.cloud',
  },
  {
    key: 'statstream',
    image: <img src={StatStreamLogoUrl} alt="StatStream" />,
    href: 'https://statstream.ai',
  },
  {
    key: 'linksociety',
    image: (
      <>
        <img src={LinkSocietyLogoUrl} alt="" />
        <span className={styles.userName}>Link Society</span>
      </>
    ),
    href: 'https://link-society.com',
  },
]

import LinkSocietyLogo from '@site/static/img/users/link-society.png'
import SaashupLogo from '@site/static/img/users/saashup.svg'
import StatStreamLogo from '@site/static/img/users/statstream.png'

const users = [
  {
    name: 'SaaShup',
    href: 'https://saashup.cloud',
    logo: <SaashupLogo aria-hidden="true" />,
  },
  {
    name: 'StatStream',
    href: 'https://statstream.ai',
    logo: <img src={StatStreamLogo} alt="" />,
  },
  {
    name: 'Link Society',
    href: 'https://link-society.com',
    logo: (
      <>
        <img src={LinkSocietyLogo} alt="" />
        <span>Link Society</span>
      </>
    ),
  },
]

export default users

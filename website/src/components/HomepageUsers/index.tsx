import clsx from 'clsx'

import Heading from '@theme/Heading'

import styles from './styles.module.css'

import users from './users'

function UserList({ hidden = false }: { hidden?: boolean }) {
  return (
    <div className={styles.carouselGroup} aria-hidden={hidden}>
      {users.map((user) => (
        <a
          key={user.key}
          className={styles.userCard}
          href={user.href}
          target="_blank"
          rel="noopener noreferrer"
          tabIndex={hidden ? -1 : undefined}
        >
          {user.image}
        </a>
      ))}
    </div>
  )
}

export default function HomepageUsers() {
  return (
    <section className="container">
      <Heading as="h2" className={clsx('text--center', styles.sectionTitle)}>
        Already running in production
      </Heading>

      <div className={styles.carousel}>
        <div className={styles.carouselTrack}>
          <UserList />
          {/* three clones keep the strip filled on wide screens; the -25% shift equals one copy */}
          <UserList hidden />
          <UserList hidden />
          <UserList hidden />
        </div>
      </div>
    </section>
  )
}

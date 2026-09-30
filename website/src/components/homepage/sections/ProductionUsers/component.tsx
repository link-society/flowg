import layoutStyles from '@site/src/components/homepage/styles.module.css'

import users from './users'

import styles from './styles.module.css'

const ProductionUsers = () => {
  return (
    <section className={styles.proof} aria-labelledby="fg-proof-title">
      <div className={layoutStyles.shell}>
        <h2 id="fg-proof-title">
          Already running in{' '}
          <span className={layoutStyles.emphasisBlue}>production</span>
        </h2>
        <div className={styles.carousel}>
          <div className={styles.carouselTrack}>
            {[0, 1, 2, 3].map((copy) => (
              <div
                className={styles.carouselGroup}
                key={copy}
                aria-hidden={copy !== 0 ? true : undefined}
              >
                {users.map((user) => (
                  <a
                    key={user.name}
                    href={user.href}
                    aria-label={user.name}
                    tabIndex={copy !== 0 ? -1 : undefined}
                  >
                    {user.logo}
                  </a>
                ))}
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}

export default ProductionUsers

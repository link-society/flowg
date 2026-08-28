import clsx from 'clsx'

import Link from '@docusaurus/Link'
import useDocusaurusContext from '@docusaurus/useDocusaurusContext'

import Heading from '@theme/Heading'

import LogoUrl from '@site/static/img/logo.png'

import styles from './styles.module.css'

export default function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext()

  return (
    <header className={clsx('hero', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className={`hero__title ${styles.titleWithLogo}`}>
          <img
            className={styles.logo}
            src={LogoUrl}
            alt="Flowg logo"
          />
          <span>{siteConfig.title}</span>
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link
            className="button button--primary button--lg"
            to="/docs"
          >
            Get Started
          </Link>
          <Link
            className="button button--secondary button--lg"
            to="https://github.com/link-society/flowg"
            target="_blank"
          >
            Bug Tracker
          </Link>
          <Link
            className="button button--success button--lg"
            to="https://demo.flowg.cloud"
            target="_blank"
          >
            Try the Demo
          </Link>
        </div>
        <div className={styles.news}>
          <Link to="/blog/osmc-2025-video">
            🚩
            Check out the video of the FlowG talk at the
            <b>OpenSource Monitoring Conference 2025</b>.
            🚩
          </Link>
        </div>
      </div>
    </header>
  )
}

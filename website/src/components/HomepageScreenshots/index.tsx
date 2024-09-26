import clsx from 'clsx'

import styles from './styles.module.css'

type ScreenshotItem = {
  title: string
  image: string
}

const ScreenshotList: ScreenshotItem[] = [
  {
    title: 'Pipeline Editor',
    image: '/img/screenshots/pipelines.png',
  },
  {
    title: 'Stream View',
    image: '/img/screenshots/streams.png',
  },
]

function Screenshot({ title, image }: ScreenshotItem): JSX.Element {
  return (
    <div className="row margin-vert--lg">
      <div className={clsx('col col--12')}>
        <img className={styles.screenshot} src={image} alt={title} />
      </div>
    </div>
  )
}

export default function HomepageScreenshots(): JSX.Element {
  return (
    <section>
      <div className="container">
        {ScreenshotList.map((item,  idx) => (
          <Screenshot key={idx} {...item} />
        ))}
      </div>
    </section>
  )
}

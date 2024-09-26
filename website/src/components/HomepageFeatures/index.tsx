import clsx from 'clsx'

import Heading from '@theme/Heading'

import styles from './styles.module.css'

type FeatureItem = {
  title: string;
  Svg: React.ComponentType<React.ComponentProps<'svg'>>;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Pipelines',
    Svg: require('@site/static/img/features/pipelines.svg').default,
    description: (
      <>
        Build and manage log processing pipelines visually, with support for
        transformations with VRL scripts and conditional routing to dedicated
        streams.
      </>
    ),
  },
  {
    title: 'Dynamic Logs',
    Svg: require('@site/static/img/features/dynamic-logs.svg').default,
    description: (
      <>
        Manage logs with varying structures without needing predefined schemas.
      </>
    ),
  },
  {
    title: 'Interoperability',
    Svg: require('@site/static/img/features/interoperability.svg').default,
    description: (
      <>
        Integrate with other softwares for log ingestion (via API or Syslog) and
        alerting (via Webhooks).
      </>
    ),
  },
  {
    title: 'Observability',
    Svg: require('@site/static/img/features/observability.svg').default,
    description: (
      <>
        Query and visualize logs in real-time in the integrated Web interface.
      </>
    ),
  },
]

function Feature({title, Svg, description}: FeatureItem) {
  return (
    <div className={clsx('col col--3')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  )
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  )
}

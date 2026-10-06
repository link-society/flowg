import type { CSSProperties } from 'react'

import { flowHeight, flowStages } from './config'
import EventFlow from './EventFlow'
import { sources, destinations } from './integrations'
import type { LogFlowDiagramProps } from './types'

import styles from './styles.module.css'

const LogFlowDiagram = ({ labelledBy }: LogFlowDiagramProps) => {
  return (
    <figure className={styles.flowFigure} aria-labelledby={labelledBy}>
      <div className={styles.flowSources} role="group" aria-label="Log sources">
        {sources.map((source) => (
          <div
            className={styles.flowSource}
            role="img"
            aria-label={source.name}
            key={source.name}
          >
            {source.logo}
          </div>
        ))}
      </div>
      <div className={styles.flowScene}>
        <EventFlow />
        <EventFlow compact />
        <ol
          className={styles.flowStages}
          aria-label="How FlowG processes events"
        >
          {flowStages.map((stage, i) => (
            <li
              key={stage.name}
              style={
                {
                  '--stage-y': `${(stage.y / flowHeight) * 100}%`,
                } as CSSProperties
              }
            >
              <span className={styles.flowStep}>0{i + 1}</span>
              <div>
                <strong>{stage.name}</strong>
                <span>{stage.description}</span>
              </div>
            </li>
          ))}
        </ol>
      </div>
      <ul className={styles.flowDestinations} aria-label="Log destinations">
        {destinations.map((item) => (
          <li key={item.name}>
            <span role="img" aria-label={item.name}>
              {item.logo}
            </span>
          </li>
        ))}
      </ul>
    </figure>
  )
}

export default LogFlowDiagram

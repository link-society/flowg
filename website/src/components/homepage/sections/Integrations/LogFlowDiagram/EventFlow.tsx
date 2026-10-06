import clsx from 'clsx'

import {
  flowKinds,
  flowStages,
  flowTiming,
  flowDuration,
  flowHeight,
  sourceLanes,
  getFlowLayout,
} from './config'
import { eventFlowPath, eventPoints } from './geometry'
import type { EventFlowProps, EventShapeProps } from './types'

import styles from './styles.module.css'

const EventShape = ({ kind }: EventShapeProps) => {
  return <polygon points={eventPoints(kind)} />
}

const EventFlow = ({ compact = false }: EventFlowProps) => {
  const { width, center, scaleY } = getFlowLayout(compact)
  return (
    <svg
      className={clsx(
        styles.eventFlow,
        compact ? styles.eventFlowCompact : styles.eventFlowWide
      )}
      viewBox={`0 0 ${width} ${flowHeight * scaleY}`}
      fill="none"
      aria-hidden="true"
      focusable="false"
    >
      <g strokeLinecap="round" strokeWidth="1.5">
        {sourceLanes.map((source) => (
          <path
            key={source}
            data-lane="input"
            d={eventFlowPath(source, 0, compact).incoming}
            stroke="#d6dfe9"
          />
        ))}
        <path
          d={`M${center} ${140 * scaleY}V${580 * scaleY}`}
          stroke="#c4d4e9"
        />
        {flowKinds
          .filter((kind) => kind.lane >= 0)
          .map((kind) => (
            <path
              key={kind.name}
              data-lane="output"
              d={eventFlowPath(0, kind.lane, compact).outgoing}
              stroke={kind.color}
              strokeOpacity=".4"
            />
          ))}
      </g>
      {flowStages.map((stage, index) => (
        <g key={stage.name}>
          <path
            d={`M${compact ? 480 * 0.58 : index % 2 ? 1200 * 0.64 : 1200 * 0.36} ${stage.y * scaleY}H${center}`}
            stroke="#c9d9ef"
            strokeWidth="1.5"
          />
          <circle
            cx={center}
            cy={stage.y * scaleY}
            r="17"
            fill="#f4f8ff"
            stroke="#cfddf1"
          />
        </g>
      ))}
      <g className={styles.flowMoving}>
        {sourceLanes.flatMap((source) =>
          flowKinds.map((kind, index) => {
            const route = eventFlowPath(source, kind.lane, compact)
            const begin = `${-(((index - source * 2 + 10) % 5) * 4 + source)}s`
            const timing = {
              dur: `${flowDuration}s`,
              begin,
              repeatCount: 'indefinite',
              calcMode: 'linear' as const,
            }
            const dot = eventPoints(kind.name, true),
              shape = eventPoints(kind.name)
            return (
              <g
                key={`${source}-${kind.name}`}
                data-source={source}
                data-kind={kind.name}
                data-delay={-parseFloat(begin)}
              >
                <animateMotion
                  {...timing}
                  path={route.path}
                  keyPoints={route.keyPoints}
                  keyTimes={flowTiming.join(';')}
                />
                {kind.lane < 0 && (
                  <animate
                    {...timing}
                    attributeName="opacity"
                    values="1;1;0;0"
                    keyTimes="0;.39;.425;1"
                  />
                )}
                <polygon
                  data-particle="true"
                  points={dot}
                  fill="#aab6c5"
                  stroke={kind.outline}
                  strokeWidth="0"
                  strokeLinejoin="round"
                  paintOrder="stroke fill"
                >
                  <animate
                    {...timing}
                    attributeName="fill"
                    values={`#aab6c5;#aab6c5;${kind.color};${kind.color}`}
                    keyTimes="0;.22;.255;1"
                  />
                  {kind.lane >= 0 && (
                    <>
                      <animate
                        {...timing}
                        attributeName="points"
                        values={`${dot};${dot};${shape};${shape}`}
                        keyTimes="0;.55;.585;1"
                      />
                      <animate
                        {...timing}
                        attributeName="stroke-width"
                        values="0;0;4;4"
                        keyTimes="0;.71;.745;1"
                      />
                    </>
                  )}
                </polygon>
              </g>
            )
          })
        )}
      </g>
      <g className={styles.flowStill}>
        {sourceLanes.flatMap((source) =>
          [16, 30].map((y) => (
            <circle
              key={`${source}-${y}`}
              cx={width / 8 + (source * width) / 4}
              cy={y * scaleY}
              r="6.5"
              fill="#aab6c5"
            />
          ))
        )}
        {flowKinds.map((kind, i) => (
          <circle
            key={kind.name}
            cx={center}
            cy={(160 + i * 17) * scaleY}
            r="6.5"
            fill={kind.color}
          />
        ))}
        {flowKinds
          .filter((kind) => kind.lane >= 0)
          .map((kind) => (
            <g key={kind.name}>
              <circle
                cx={center}
                cy={(277 + kind.lane * 20) * scaleY}
                r="6.5"
                fill={kind.color}
              />
              <g
                transform={`translate(${center} ${(387 + kind.lane * 20) * scaleY})`}
                fill={kind.color}
              >
                <EventShape kind={kind.name} />
              </g>
              <g
                transform={`translate(${center} ${(497 + kind.lane * 20) * scaleY})`}
                fill={kind.color}
                stroke={kind.outline}
                strokeWidth="4"
                strokeLinejoin="round"
                paintOrder="stroke fill"
              >
                <EventShape kind={kind.name} />
              </g>
              <g
                transform={`translate(${width / 8 + (kind.lane * width) / 4} ${690 * scaleY})`}
                fill={kind.color}
                stroke={kind.outline}
                strokeWidth="4"
                strokeLinejoin="round"
                paintOrder="stroke fill"
              >
                <EventShape kind={kind.name} />
              </g>
            </g>
          ))}
      </g>
    </svg>
  )
}

export default EventFlow

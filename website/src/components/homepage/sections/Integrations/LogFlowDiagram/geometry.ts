import { flowHeight, flowStages, getFlowLayout, type FlowKind } from './config'

// Two quarter-circle bends join vertical and horizontal runs. The exact
// length keeps each event's changes synchronized with its pipeline node.
function roundedFlowLink(
  fromX: number,
  fromY: number,
  toX: number,
  toY: number,
  bendY: number
) {
  const distance = Math.abs(toX - fromX)
  if (!distance)
    return { path: `M${fromX} ${fromY}V${toY}`, length: toY - fromY }
  const direction = Math.sign(toX - fromX)
  const radius = Math.min(16, distance / 2, bendY - fromY, toY - bendY)
  const firstSweep = direction > 0 ? 0 : 1
  return {
    path: `M${fromX} ${fromY}V${bendY - radius}A${radius} ${radius} 0 0 ${firstSweep} ${fromX + direction * radius} ${bendY}H${toX - direction * radius}A${radius} ${radius} 0 0 ${1 - firstSweep} ${toX} ${bendY + radius}V${toY}`,
    length: toY - fromY + distance + (Math.PI - 4) * radius,
  }
}

export function eventFlowPath(source: number, lane: number, compact: boolean) {
  const { width, center, scaleY } = getFlowLayout(compact)
  const inputX = width / 8 + (source * width) / 4
  const outputX = width / 8 + (Math.max(lane, 0) * width) / 4
  const input = roundedFlowLink(inputX, 0, center, 140 * scaleY, 70 * scaleY)
  const output = roundedFlowLink(
    center,
    580 * scaleY,
    outputX,
    flowHeight * scaleY,
    660 * scaleY
  )
  const firstStage = flowStages[0].y
  const lastStage = flowStages[flowStages.length - 1].y
  const length =
    input.length + (lastStage - firstStage) * scaleY + output.length
  const distances = [
    0,
    ...flowStages.map(
      (stage) => input.length + (stage.y - firstStage) * scaleY
    ),
    length,
  ]
  // Red events stop at Filter and fade out there; they never reach a forwarder.
  if (lane < 0) distances.fill(distances[2], 3)
  return {
    incoming: input.path,
    outgoing: output.path,
    path: `${input.path}V${580 * scaleY}${output.path.replace(/^M[^V]+/, '')}`,
    keyPoints: distances.map((n) => n / length).join(';'),
  }
}

// Matching vertex counts let SVG interpolate the entire silhouette, rather
// than swapping a circle for a shape. Sixty preserves every polygon corner.
export function eventPoints(kind: FlowKind['name'], circle = false) {
  const vertices: Record<string, number[][]> = {
    green: [
      [-7, -7],
      [7, -7],
      [7, 7],
      [-7, 7],
    ],
    blue: [
      [0, -9],
      [8, 6],
      [-8, 6],
    ],
    yellow: Array.from({ length: 10 }, (_, i) => {
      const angle = -Math.PI / 2 + (i * Math.PI) / 5
      const radius = i % 2 ? 4.5 : 10
      return [Math.cos(angle) * radius, Math.sin(angle) * radius]
    }),
    purple: [
      [0, -8],
      [7, -4],
      [7, 4],
      [0, 8],
      [-7, 4],
      [-7, -4],
    ],
  }
  const corners = vertices[kind] || vertices.blue
  const start = Math.atan2(corners[0][1], corners[0][0])
  return Array.from({ length: 60 }, (_, i) => {
    if (circle) {
      const angle = start + (i * Math.PI) / 30
      return `${(Math.cos(angle) * 6.5).toFixed(3)},${(Math.sin(angle) * 6.5).toFixed(3)}`
    }
    const position = (i * corners.length) / 60
    const index = Math.floor(position),
      fraction = position - index
    const from = corners[index],
      to = corners[(index + 1) % corners.length]
    return from
      .map((value, axis) => (value + (to[axis] - value) * fraction).toFixed(3))
      .join(',')
  }).join(' ')
}

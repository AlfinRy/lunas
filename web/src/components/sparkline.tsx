/**
 * F1.4 — the DSO sparkline. Pure SVG, no chart library: a 14-point line with
 * the jade chart-cousin color, tabular by construction.
 */
export function Sparkline({
  points,
  width = 200,
  height = 36,
}: {
  points: number[];
  width?: number;
  height?: number;
}) {
  if (!points.length) return null;
  const max = Math.max(...points, 1);
  const min = Math.min(...points, 0);
  const span = Math.max(max - min, 1);
  const step = points.length > 1 ? width / (points.length - 1) : width;

  const coords = points.map((p, i) => {
    const x = i * step;
    const y = height - 3 - ((p - min) / span) * (height - 6);
    return [x, y] as const;
  });

  const d = coords.map(([x, y], i) => `${i === 0 ? "M" : "L"}${x.toFixed(1)},${y.toFixed(1)}`).join(" ");
  const last = coords[coords.length - 1];

  return (
    <svg
      viewBox={`0 0 ${width} ${height}`}
      width="100%"
      height={height}
      role="img"
      aria-label={`DSO trend: ${Math.round(points[points.length - 1])} days now`}
      preserveAspectRatio="none"
    >
      <path d={`${d} L${width},${height} L0,${height} Z`} fill="var(--primary)" opacity="0.08" />
      <path d={d} fill="none" stroke="var(--chart-1)" strokeWidth="1.5" strokeLinecap="round" />
      <circle cx={last[0]} cy={last[1]} r="2.5" fill="var(--chart-1)" />
    </svg>
  );
}

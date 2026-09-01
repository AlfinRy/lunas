import type { ReactNode } from "react";
import { cn } from "@/lib/cn";
import { SquircleSurface } from "./oa";

/** Cards are squircle surfaces: white plates on the stage (oa-design 01). */
export function Card({ className, children }: { className?: string; children?: ReactNode }) {
  return (
    <SquircleSurface
      className={cn("border border-border p-5 shadow-[0_1px_2px_rgba(0,0,0,0.06)]", className)}
    >
      {children}
    </SquircleSurface>
  );
}

/** Fixed-height mini KPI card: eyebrow, tabular number, sub-line (+sparkline). */
export function KpiCard({
  label,
  value,
  sub,
  tone = "neutral",
  children,
}: {
  label: string;
  value: string | null;
  sub?: string;
  tone?: "neutral" | "success" | "warning";
  children?: ReactNode;
}) {
  const toneClass =
    tone === "success" ? "text-success" : tone === "warning" ? "text-warning" : "text-foreground";
  return (
    <Card className="flex h-[112px] flex-col justify-between">
      <p className="text-label uppercase tracking-[0.06em] text-muted-foreground">{label}</p>
      <div>
        <p
          className={cn(
            "tnum text-[27px] font-medium leading-none tracking-[-0.01em]",
            value === null && "text-muted-foreground",
            toneClass,
          )}
        >
          {value ?? "—"}
        </p>
        {sub ? <p className="mt-1.5 text-xs text-muted-foreground">{sub}</p> : null}
        {children}
      </div>
    </Card>
  );
}

/** Two-layer plate: frame with a title strip, holding a recessed grey inset. */
export function Plate({
  title,
  meta,
  action,
  children,
  className,
  contentClassName,
}: {
  title: ReactNode;
  meta?: ReactNode;
  action?: ReactNode;
  children: ReactNode;
  className?: string;
  contentClassName?: string;
}) {
  return (
    <SquircleSurface
      className={cn("border border-border p-1 shadow-[0_1px_2px_rgba(0,0,0,0.06)]", className)}
    >
      <div className="flex min-h-9 flex-wrap items-center justify-between gap-2 px-3 pb-1 pt-1.5">
        <h2 className="ml-1 flex min-w-0 items-center gap-2 text-sm font-medium text-foreground/80">{title}</h2>
        <div className="flex items-center gap-2 pr-1">
          {meta ? <span className="text-xs text-muted-foreground">{meta}</span> : null}
          {action}
        </div>
      </div>
      <SquircleSurface
        className={cn(
          "overflow-hidden rounded-[22px] border border-border bg-background [--card-clip-radius:12px] sm:rounded-[44px] sm:[--card-clip-radius:17px]",
          contentClassName,
        )}
      >
        {children}
      </SquircleSurface>
    </SquircleSurface>
  );
}

/** Skeleton block — pixel-matched waits; data arrives by blur, not by pop. */
export function Skeleton({ className }: { className?: string }) {
  return (
    <div
      aria-hidden="true"
      className={cn("animate-pulse rounded-md bg-accent-wash", className)}
      style={{ animationDuration: "1.6s" }}
    />
  );
}

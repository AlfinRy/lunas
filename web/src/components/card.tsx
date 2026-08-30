import type { ReactNode } from "react";
import { cn } from "@/lib/cn";

/** OA 01-squircle-card: white plate on the grey stage, no shadow. */
export function Card({ className, children }: { className?: string; children?: ReactNode }) {
  return (
    <section className={cn("rounded-2xl border border-border bg-card", className)}>{children}</section>
  );
}

/** Fixed-height mini card for KPIs: eyebrow, tabular number, sub-line. */
export function KpiCard({
  label,
  value,
  sub,
  tone = "neutral",
}: {
  label: string;
  value: string | null;
  sub?: string;
  tone?: "neutral" | "success" | "warning";
}) {
  const toneClass =
    tone === "success" ? "text-success" : tone === "warning" ? "text-warning" : "text-foreground";
  return (
    <Card className="flex h-[104px] flex-col justify-between p-5">
      <p className="text-label uppercase tracking-[0.06em] text-muted-foreground">{label}</p>
      <div>
        <p
          className={cn(
            "tnum text-[26px] font-medium leading-none tracking-tight",
            value === null && "text-muted-foreground",
            toneClass,
          )}
        >
          {value ?? "—"}
        </p>
        {sub ? <p className="mt-1.5 text-xs text-muted-foreground">{sub}</p> : null}
      </div>
    </Card>
  );
}

/** Skeleton block — pixel-matched waits (OA 07): same box, shimmering quietly. */
export function Skeleton({ className }: { className?: string }) {
  return (
    <div
      aria-hidden="true"
      className={cn("animate-pulse rounded-md bg-accent-wash", className)}
      style={{ animationDuration: "1.6s" }}
    />
  );
}

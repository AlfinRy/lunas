import type { InvoiceStatus } from "@/api/client";
import { cn } from "@/lib/cn";

/**
 * OA rule 7/8: states are neutral-surface pills — a small dot + semantic text
 * tint, never a colored fill. The jade accent stays reserved for actions.
 */
const meta: Record<string, { label: string; dot: string; text: string }> = {
  paid: { label: "Paid", dot: "bg-success", text: "text-success" },
  chasing: { label: "Chasing", dot: "bg-destructive", text: "text-foreground" },
  scheduled: { label: "Due soon", dot: "bg-warning", text: "text-warning" },
  paused: { label: "Paused", dot: "bg-muted-foreground", text: "text-muted-foreground" },
  written_off: { label: "Written off", dot: "bg-muted-foreground", text: "text-muted-foreground" },
  draft: { label: "Draft", dot: "bg-muted-foreground", text: "text-muted-foreground" },
};

export function StatusPill({
  status,
  daysOverdue,
}: {
  status: InvoiceStatus;
  daysOverdue: number;
}) {
  const m = meta[status] ?? { label: status, dot: "bg-muted-foreground", text: "text-muted-foreground" };
  const overdueOpen = status === "chasing" && daysOverdue > 0;
  return (
    <span
      className={cn(
        "inline-flex w-fit items-center gap-1.5 rounded-pill border border-border bg-card px-2.5 py-1 text-xs whitespace-nowrap",
        overdueOpen ? "text-destructive" : m.text,
      )}
    >
      <span
        aria-hidden="true"
        className={cn("size-1.5 shrink-0 rounded-full", overdueOpen ? "bg-destructive" : m.dot)}
      />
      {overdueOpen ? `Overdue ${daysOverdue}d` : m.label}
    </span>
  );
}

import type { Invoice } from "@/api/client";
import { Link } from "@tanstack/react-router";
import { Skeleton } from "./card";
import { StatusPill } from "./status-pill";
import { money, fmtDate } from "@/lib/format";
import { cn } from "@/lib/cn";

const agentLabel: Record<string, string> = {
  idle: "Idle",
  planning: "Planning",
  awaiting_approval: "Awaiting approval",
  chasing: "Chasing",
  holding: "Holding",
  stopped: "Stopped",
};

export function InvoiceTable({
  invoices,
  loading,
  skeletonRows = 4,
}: {
  invoices?: Invoice[];
  loading?: boolean;
  skeletonRows?: number;
}) {
  return (
    <div role="table" aria-label="Invoices" aria-busy={loading || undefined}>
      {/* header row */}
      <div
        role="row"
        className="hidden grid-cols-[minmax(140px,1.6fr)_92px_110px_110px_130px_130px] gap-3 border-t border-border px-5 py-2 text-xs text-muted-foreground md:grid"
      >
        <span role="columnheader">Payer</span>
        <span role="columnheader" className="text-end">
          Number
        </span>
        <span role="columnheader" className="text-end">
          Amount
        </span>
        <span role="columnheader" className="text-end">
          Due
        </span>
        <span role="columnheader">Status</span>
        <span role="columnheader">Agent</span>
      </div>

      {loading
        ? Array.from({ length: skeletonRows }).map((_, i) => (
            <div
              key={i}
              role="row"
              className="grid grid-cols-2 items-center gap-3 border-t border-border px-5 py-3.5 md:grid-cols-[minmax(140px,1.6fr)_92px_110px_110px_130px_130px]"
            >
              <Skeleton className="h-4 w-28" />
              <Skeleton className="h-4 w-16 justify-self-end" />
              <Skeleton className="h-4 w-20 justify-self-end" />
              <Skeleton className="hidden h-4 w-16 justify-self-end md:block" />
              <Skeleton className="hidden h-5 w-24 md:block" />
              <Skeleton className="hidden h-4 w-28 md:block" />
            </div>
          ))
        : (invoices ?? []).map((inv) => (
            <Link
              key={inv.id}
              to="/invoices/$id"
              params={{ id: String(inv.id) }}
              className="contents"
              aria-label={`${inv.client_name}, invoice ${inv.number}`}
            >
              <Row invoice={inv} />
            </Link>
          ))}
    </div>
  );
}

function Row({ invoice }: { invoice: Invoice }) {
  const overdue = (invoice.days_overdue ?? 0) > 0;
  return (
    <div
      role="row"
      className="grid grid-cols-2 items-center gap-x-3 gap-y-1.5 border-t border-border px-5 py-3.5 transition-colors duration-150 ease-out hover:bg-accent-wash md:grid-cols-[minmax(140px,1.6fr)_92px_110px_110px_130px_130px] md:gap-y-0"
    >
      {/* payer + number on mobile */}
      <div role="cell" className="min-w-0">
        <p className="truncate text-sm">{invoice.client_name}</p>
        <p className="truncate text-xs text-muted-foreground md:hidden">
          <span className="font-mono">{invoice.number}</span>
        </p>
      </div>

      <p role="cell" className="hidden truncate text-end text-xs text-muted-foreground md:block">
        <span className="font-mono">{invoice.number}</span>
      </p>

      <p role="cell" className="tnum text-end text-sm font-medium">
        {money(invoice.amount_cents, invoice.currency)}
      </p>

      <p
        role="cell"
        className={cn(
          "tnum text-end text-xs",
          overdue ? "text-destructive" : "text-muted-foreground",
        )}
      >
        {fmtDate(invoice.due_on)}
        {overdue ? ` · ${invoice.days_overdue}d late` : ""}
      </p>

      <div role="cell" className="col-span-2 md:col-span-1">
        <StatusPill status={invoice.status} daysOverdue={invoice.days_overdue ?? 0} />
      </div>

      <p role="cell" className="col-span-2 text-xs text-muted-foreground md:col-span-1">
        {agentLabel[invoice.agent_state] ?? invoice.agent_state}
      </p>
    </div>
  );
}

import { createRoute, Link } from "@tanstack/react-router";
import { ArrowRight } from "lucide-react";
import { Route as rootRoute } from "./__root";
import { Card } from "@/components/card";
import { useClients, useInvoices } from "@/api/queries";
import { money } from "@/lib/format";

const reliabilityMeta: Record<string, { label: string; cls: string }> = {
  pays_on_time: { label: "Pays on time", cls: "text-success" },
  usually_late: { label: "Usually late", cls: "text-warning" },
  chronically_late: { label: "Chronically late", cls: "text-destructive" },
};

function ClientsPage() {
  const clients = useClients();
  const invoices = useInvoices();

  return (
    <main className="mx-auto flex max-w-[1200px] flex-col gap-6">
      <header>
        <h1 className="text-[26px] font-medium leading-tight tracking-tight">Clients</h1>
        <p className="mt-0.5 text-sm text-muted-foreground">
          Payment behavior feeds the agent — it chases late payers earlier and gentle-clients warmer.
        </p>
      </header>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
        {(clients.data ?? []).map((c) => {
          const open = (invoices.data ?? []).filter(
            (i) => i.client_id === c.id && i.status !== "paid" && i.status !== "written_off",
          );
          const owed = open.reduce((s, i) => s + i.amount_cents - (i.amount_paid_cents ?? 0), 0);
          const rel = c.payment_score ? reliabilityMeta[c.payment_score.reliability] : null;
          return (
            <Card key={c.id} className="flex flex-col gap-3 p-5">
              <div className="flex items-start justify-between gap-2">
                <div className="min-w-0">
                  <p className="truncate font-medium">{c.name}</p>
                  <p className="truncate text-xs text-muted-foreground">{c.email}</p>
                </div>
                {rel ? (
                  <span className={`flex w-fit shrink-0 items-center gap-1.5 text-xs ${rel.cls}`}>
                    <span aria-hidden="true" className="size-1.5 rounded-full bg-current" />
                    {rel.label}
                  </span>
                ) : null}
              </div>

              {c.payment_score ? (
                <p className="tnum text-xs text-muted-foreground">
                  {c.payment_score.settled_count} settled · avg{" "}
                  {c.payment_score.avg_days_late > 0
                    ? `${c.payment_score.avg_days_late} days late`
                    : "on time or early"}
                </p>
              ) : (
                <p className="text-xs text-muted-foreground">
                  {open.length ? "Learning this payer's rhythm…" : "No settled invoices yet."}
                </p>
              )}

              {open.length ? (
                <div className="mt-auto flex items-end justify-between border-t border-border pt-3">
                  <div>
                    <p className="text-label uppercase tracking-[0.06em] text-muted-foreground">Owed</p>
                    <p className="tnum text-lg font-medium">{money(owed, open[0]?.currency)}</p>
                  </div>
                  <Link
                    to="/invoices"
                    className="flex items-center gap-1 text-xs text-info hover:underline hover:underline-offset-4"
                  >
                    {open.length} open <ArrowRight size={12} strokeWidth={1.5} aria-hidden="true" />
                  </Link>
                </div>
              ) : (
                <p className="mt-auto border-t border-border pt-3 text-xs text-muted-foreground">Nothing outstanding.</p>
              )}
            </Card>
          );
        })}
      </div>

      {clients.data && clients.data.length === 0 ? (
        <Card className="px-5 py-10">
          <p className="font-medium">No clients yet</p>
          <p className="mt-1 text-sm text-muted-foreground">Clients appear when you add your first invoice.</p>
        </Card>
      ) : null}
    </main>
  );
}

export const Route = createRoute({
  getParentRoute: () => rootRoute,
  path: "/clients",
  component: ClientsPage,
});

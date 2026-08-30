import { useState } from "react";
import { createRoute } from "@tanstack/react-router";
import { ClipboardPaste } from "lucide-react";
import { Route as rootRoute } from "./__root";
import { Card, KpiCard } from "@/components/card";
import { Button } from "@/components/button";
import { InvoiceTable } from "@/components/invoice-table";
import { AddInvoiceModal } from "@/components/add-invoice-modal";
import { RecordPaymentModal } from "@/components/record-payment-modal";
import { Sparkline } from "@/components/sparkline";
import { DemoTimePill } from "@/components/demo-time-pill";
import { useDashboard, useInvoices } from "@/api/queries";
import { money } from "@/lib/format";

function DashboardPage() {
  const dash = useDashboard();
  const invoices = useInvoices();
  const [adding, setAdding] = useState(false);
  const [recording, setRecording] = useState(false);

  return (
    <main className="mx-auto flex max-w-[1200px] flex-col gap-6">
      <header className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="text-[26px] font-medium leading-tight tracking-tight">Dashboard</h1>
          <p className="mt-0.5 text-sm text-muted-foreground">
            {dash.data
              ? `As of ${fmtDate(dash.data.sim_now)} — Lunas is watching ${invoices.data?.length ?? "…"} invoices.`
              : "Loading your numbers…"}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="secondary" onClick={() => setRecording(true)}>
            <ClipboardPaste size={15} strokeWidth={1.5} aria-hidden="true" />
            Paste payment
          </Button>
          <DemoTimePill />
        </div>
      </header>

      <section aria-label="Key numbers" className="grid grid-cols-2 gap-4 xl:grid-cols-4">
        <KpiCard
          label="Outstanding"
          value={dash.isPending ? null : money(dash.data!.outstanding_cents)}
          sub={`${dash.data?.counts.due_soon ?? 0} due soon`}
        />
        <KpiCard
          label="Overdue"
          value={dash.isPending ? null : money(dash.data!.overdue_cents)}
          sub={`${dash.data?.counts.overdue ?? 0} invoices`}
          tone={dash.data && dash.data.overdue_cents > 0 ? "warning" : "neutral"}
        />
        <KpiCard
          label="Recovered"
          value={dash.isPending ? null : money(dash.data!.recovered_cents)}
          sub="lifetime, after chasing"
          tone="success"
        />
        <KpiCard
          label="Days sales outstanding"
          value={dash.isPending ? null : `${Math.round(dash.data!.dso_days)} days`}
          sub="30-day rolling"
        >
          {dash.data?.dso_trend?.length ? (
            <div className="mt-2">
              <Sparkline points={dash.data.dso_trend.map((p) => p.dso_days)} />
            </div>
          ) : null}
        </KpiCard>
      </section>

      <Card className="p-0">
        <div className="flex items-center justify-between px-5 py-4">
          <h2 className="text-[15px] font-medium tracking-tight">Invoices</h2>
          <div className="flex items-center gap-3">
            <p className="text-xs text-muted-foreground">Overdue first</p>
            <Button variant="primary" onClick={() => setAdding(true)}>
              Add invoice
            </Button>
          </div>
        </div>
        <InvoiceTable invoices={invoices.data} loading={invoices.isPending} skeletonRows={5} />
      </Card>

      <AddInvoiceModal open={adding} onClose={() => setAdding(false)} />
      <RecordPaymentModal open={recording} onClose={() => setRecording(false)} />
    </main>
  );
}

function fmtDate(iso: string): string {
  return new Date(`${iso.slice(0, 10)}T00:00:00Z`).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    timeZone: "UTC",
  });
}

export const Route = createRoute({
  getParentRoute: () => rootRoute,
  path: "/",
  component: DashboardPage,
});

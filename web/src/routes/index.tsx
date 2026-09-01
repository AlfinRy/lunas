import { useState } from "react";
import { createRoute } from "@tanstack/react-router";
import { motion } from "motion/react";
import { ClipboardPaste, Plus } from "lucide-react";
import { Route as rootRoute } from "./__root";
import { KpiCard, Plate } from "@/components/card";
import { Button } from "@/components/button";
import { InvoiceTable } from "@/components/invoice-table";
import { AddInvoiceModal } from "@/components/add-invoice-modal";
import { RecordPaymentModal } from "@/components/record-payment-modal";
import { Sparkline } from "@/components/sparkline";
import { DemoTimePill } from "@/components/demo-time-pill";
import { useDashboard, useInvoices } from "@/api/queries";
import { money } from "@/lib/format";

/** Enter once by ~70ms beats (oa reveal); data swaps in by blur, not by pop. */
const enter = (i: number) => ({
  initial: { opacity: 0, y: 6, filter: "blur(4px)" },
  animate: { opacity: 1, y: 0, filter: "blur(0px)" },
  transition: { duration: 0.22, ease: [0.2, 0, 0, 1] as const, delay: i * 0.07 },
});

function DashboardPage() {
  const dash = useDashboard();
  const invoices = useInvoices();
  const [adding, setAdding] = useState(false);
  const [recording, setRecording] = useState(false);

  return (
    <main className="flex flex-col gap-6">
      <motion.header {...enter(0)} className="flex flex-wrap items-end justify-between gap-3">
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
      </motion.header>

      <section aria-label="Key numbers" className="grid grid-cols-2 gap-4 xl:grid-cols-4">
        <motion.div {...enter(1)}>
          <KpiCard
            label="Outstanding"
            value={dash.isPending ? null : money(dash.data!.outstanding_cents)}
            sub={`${dash.data?.counts.due_soon ?? 0} due soon`}
          />
        </motion.div>
        <motion.div {...enter(2)}>
          <KpiCard
            label="Overdue"
            value={dash.isPending ? null : money(dash.data!.overdue_cents)}
            sub={`${dash.data?.counts.overdue ?? 0} invoices`}
            tone={dash.data && dash.data.overdue_cents > 0 ? "warning" : "neutral"}
          />
        </motion.div>
        <motion.div {...enter(3)}>
          <KpiCard
            label="Recovered"
            value={dash.isPending ? null : money(dash.data!.recovered_cents)}
            sub="lifetime, after chasing"
            tone="success"
          />
        </motion.div>
        <motion.div {...enter(4)}>
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
        </motion.div>
      </section>

      <motion.div {...enter(5)}>
        <Plate
          title={
            <>
              <svg width="15" height="15" viewBox="0 0 16 16" fill="none" aria-hidden="true" className="text-muted-foreground">
                <rect x="2" y="3.5" width="12" height="10" rx="2" stroke="currentColor" strokeWidth="1.4" />
                <path d="M4.5 7h4M4.5 9.5h7" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" />
              </svg>
              Invoices
            </>
          }
          meta="Overdue first"
          action={
            <Button variant="primary" onClick={() => setAdding(true)}>
              <Plus size={14} strokeWidth={1.8} aria-hidden="true" />
              Add invoice
            </Button>
          }
          contentClassName="py-0"
        >
          <InvoiceTable invoices={invoices.data} loading={invoices.isPending} skeletonRows={5} />
        </Plate>
      </motion.div>

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

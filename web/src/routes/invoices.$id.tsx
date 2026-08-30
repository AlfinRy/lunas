import { useState } from "react";
import { createRoute, Link, useParams } from "@tanstack/react-router";
import { ArrowLeft, Pause, Play, BadgeCheck, ClipboardPaste } from "lucide-react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Route as rootRoute } from "./__root";
import { Card } from "@/components/card";
import { Button } from "@/components/button";
import { StatusPill } from "@/components/status-pill";
import { useToast } from "@/components/toast";
import { RecordPaymentModal } from "@/components/record-payment-modal";
import { ApiError, api } from "@/api/client";
import { money, fmtDate } from "@/lib/format";

function InvoiceDetailPage() {
  const { id } = useParams({ from: "/invoices/$id" });
  const invoiceId = Number(id);
  const qc = useQueryClient();
  const { toast } = useToast();

  const [recording, setRecording] = useState(false);
  const invoice = useQuery({ queryKey: ["invoice", invoiceId], queryFn: () => api.invoice(invoiceId) });
  const activity = useQuery({ queryKey: ["activity", invoiceId], queryFn: () => api.activity(invoiceId) });

  async function setStatus(status: "paused" | "chasing" | "paid") {
    const label =
      status === "paused" ? "Chasing paused" : status === "paid" ? "Invoice settled" : "Chasing resumed";
    try {
      await api.updateInvoice(invoiceId, { status });
      toast("success", label === "Invoice settled" ? `${inv!.number} settled — chasing stopped.` : `${label} for ${inv!.number}.`);
      qc.invalidateQueries();
    } catch (err) {
      toast("error", err instanceof ApiError ? err.message : "Unable to save. Try again.");
    }
  }

  const inv = invoice.data;
  const open = inv && inv.status !== "paid" && inv.status !== "written_off";

  return (
    <main className="mx-auto flex max-w-[760px] flex-col gap-6">
      <Link
        to="/invoices"
        className="inline-flex w-fit items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft size={14} strokeWidth={1.5} aria-hidden="true" /> Invoices
      </Link>

      {invoice.isPending || !inv ? (
        <Card className="h-40 animate-pulse" />
      ) : (
        <>
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div>
              <div className="flex items-center gap-3">
                <h1 className="text-[26px] font-medium tracking-tight">{inv.client_name}</h1>
                <StatusPill status={inv.status} daysOverdue={inv.days_overdue ?? 0} />
              </div>
              <p className="mt-1 text-sm text-muted-foreground">
                <span className="font-mono">{inv.number}</span> · issued {fmtDate(inv.issued_on)} · due{" "}
                {fmtDate(inv.due_on)}
              </p>
            </div>
            <p className="tnum text-[28px] font-medium tracking-tight">{money(inv.amount_cents, inv.currency)}</p>
          </div>

          {open ? (
            <div className="flex flex-wrap gap-2">
              {inv.status === "paused" ? (
                <Button variant="secondary" onClick={() => setStatus("chasing")}>
                  <Play size={14} strokeWidth={1.5} aria-hidden="true" /> Resume chasing
                </Button>
              ) : (
                <Button variant="secondary" onClick={() => setStatus("paused")}>
                  <Pause size={14} strokeWidth={1.5} aria-hidden="true" /> Pause chasing
                </Button>
              )}
              <Button variant="primary" onClick={() => setRecording(true)}>
                <ClipboardPaste size={14} strokeWidth={1.5} aria-hidden="true" /> Record payment
              </Button>
              <Button variant="ghost" onClick={() => setStatus("paid")}>
                <BadgeCheck size={14} strokeWidth={1.5} aria-hidden="true" /> Mark as settled
              </Button>
            </div>
          ) : (
            <Card className="flex items-center gap-2.5 px-5 py-4 text-sm text-success">
              <BadgeCheck size={16} strokeWidth={1.5} aria-hidden="true" />
              Settled{inv.amount_paid_cents ? ` — ${money(inv.amount_paid_cents, inv.currency)} recovered.` : "."} Chasing stopped.
            </Card>
          )}

          {inv.notes ? (
            <Card className="px-5 py-4">
              <p className="text-label uppercase tracking-[0.06em] text-muted-foreground">Notes</p>
              <p className="mt-1.5 text-sm">{inv.notes}</p>
            </Card>
          ) : null}

          <Card className="p-5">
            <h2 className="text-[15px] font-medium tracking-tight">Activity</h2>
            <ol className="mt-4 flex flex-col">
              {(activity.data ?? []).map((a, i) => (
                <li key={a.id} className="relative flex gap-3 pb-4 last:pb-0">
                  {i < (activity.data?.length ?? 0) - 1 ? (
                    <span aria-hidden="true" className="absolute inset-y-0 start-[5px] w-px bg-border" />
                  ) : null}
                  <span
                    aria-hidden="true"
                    className="relative mt-1.5 size-[11px] shrink-0 rounded-full border-2 border-border bg-card"
                  />
                  <div className="min-w-0">
                    <p className="text-sm leading-snug">{a.message}</p>
                    <p className="mt-0.5 text-xs text-muted-foreground">{fmtDate(a.created_at)}</p>
                  </div>
                </li>
              ))}
              {activity.data && activity.data.length === 0 ? (
                <p className="text-sm text-muted-foreground">Nothing yet — the agent will log its first plan here.</p>
              ) : null}
            </ol>
          </Card>
        </>
      )}
      <RecordPaymentModal open={recording} onClose={() => setRecording(false)} presetInvoiceId={invoiceId} />
    </main>
  );
}

export const Route = createRoute({
  getParentRoute: () => rootRoute,
  path: "/invoices/$id",
  component: InvoiceDetailPage,
});

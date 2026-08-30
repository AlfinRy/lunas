import { useState, type FormEvent } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { ClipboardPaste, BadgeCheck } from "lucide-react";
import { Modal } from "./modal";
import { Button } from "./button";
import { Field } from "./add-invoice-modal";
import { useToast } from "./toast";
import { api, ApiError } from "@/api/client";
import { money } from "@/lib/format";
import { cn } from "@/lib/cn";

type Parsed = Awaited<ReturnType<typeof api.parsePayment>>;

/**
 * F5 — paste a payment notification; the matcher proposes invoices; one tap
 * settles it. This is the flow that stops the agent.
 */
export function RecordPaymentModal({
  open,
  onClose,
  presetInvoiceId,
}: {
  open: boolean;
  onClose: () => void;
  presetInvoiceId?: number;
}) {
  const qc = useQueryClient();
  const { toast } = useToast();
  const [text, setText] = useState("");
  const [parsed, setParsed] = useState<Parsed | null>(null);
  const [parsing, setParsing] = useState(false);
  const [settled, setSettled] = useState<{ number?: string; cents: number } | null>(null);

  const parse = async () => {
    setParsing(true);
    setSettled(null);
    try {
      const res = await api.parsePayment(text);
      setParsed(res);
    } catch (err) {
      toast("error", err instanceof ApiError ? err.message : "Unable to read that text. Try again.");
    } finally {
      setParsing(false);
    }
  };

  const reconcile = useMutation({
    mutationFn: (m: { invoice_id: number; confidence?: string }) =>
      api.reconcilePayment({
        invoice_id: m.invoice_id,
        amount_cents: parsed!.parsed.amount_cents,
        paid_on: parsed?.parsed.paid_on ?? undefined,
        source: "paste",
        confidence: m.confidence,
        raw_text: text,
      }),
    onSuccess: (pay) => {
      setSettled({ cents: pay.amount_cents });
      toast(
        "success",
        pay.invoice_status_after === "paid"
          ? `Invoice settled — ${money(pay.amount_cents)} recovered. Chasing stopped.`
          : `Partial payment recorded — Lunas keeps chasing the remainder.`,
      );
      qc.invalidateQueries();
      setTimeout(onClose, 1400);
    },
    onError: (err) => toast("error", err instanceof ApiError ? err.message : "Unable to save. Try again."),
  });

  function onSubmit(e: FormEvent) {
    e.preventDefault();
    void parse();
  }

  function reset() {
    setText("");
    setParsed(null);
    setSettled(null);
  }

  return (
    <Modal
      open={open}
      onClose={() => {
        reset();
        onClose();
      }}
      title={presetInvoiceId ? "Record payment" : "Paste a payment notification"}
    >
      <form onSubmit={onSubmit} className="flex flex-col gap-4">
        {presetInvoiceId === undefined ? (
          <Field label="Paste the bank or payment email">
            {(id) => (
              <textarea
                id={id}
                value={text}
                onChange={(e) => setText(e.target.value)}
                rows={4}
                placeholder="Payment received: $2,400.00 from Meridian Coffee Co. on Sep 3, 2026"
                className="w-full resize-none rounded-lg border border-input bg-card px-3 py-2 text-[16px] leading-relaxed placeholder:text-muted-foreground focus-visible:outline-2 focus-visible:outline-ring sm:text-sm"
              />
            )}
          </Field>
        ) : null}

        {!parsed ? (
          <div className="flex justify-end gap-2">
            <Button type="button" variant="ghost" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" variant="primary" loading={parsing} disabled={!text.trim()}>
              <ClipboardPaste size={15} strokeWidth={1.5} aria-hidden="true" />
              Read payment
            </Button>
          </div>
        ) : settled ? (
          <div className="flex items-center gap-2.5 rounded-2xl border border-border bg-primary-soft px-4 py-3 text-sm text-foreground">
            <BadgeCheck size={16} strokeWidth={1.5} className="shrink-0 text-success" aria-hidden="true" />
            {settled.cents ? `${money(settled.cents)} linked.` : null} The agent stopped chasing this invoice.
          </div>
        ) : (
          <div className="flex flex-col gap-3" role="status">
            <p className="text-sm">
              Found <span className="tnum font-medium">{money(parsed.parsed.amount_cents)}</span>
              {parsed.parsed.paid_on ? ` paid on ${parsed.parsed.paid_on}` : ""}.
            </p>

            {parsed.matches.length === 0 ? (
              <p className="text-sm text-muted-foreground">
                No invoice matches this payment. Check the amount, or record it manually from the invoice page.
              </p>
            ) : (
              <ul className="flex flex-col gap-2">
                {parsed.matches.map((m) => (
                  <li
                    key={m.invoice_id}
                    className="flex items-center justify-between gap-3 rounded-2xl border border-border bg-background px-4 py-3"
                  >
                    <div className="min-w-0">
                      <p className="truncate text-sm">
                        <span className="font-medium">{m.client_name}</span>{" "}
                        <span className="font-mono text-xs text-muted-foreground">{m.invoice_number}</span>
                      </p>
                      <p className="tnum mt-0.5 text-xs text-muted-foreground">
                        {money(m.amount_cents)}
                        {m.days_overdue ? ` · ${m.days_overdue}d overdue` : ""}
                      </p>
                    </div>
                    <div className="flex shrink-0 items-center gap-2">
                      <span
                        className={cn(
                          "rounded-pill border border-border px-2 py-0.5 text-[11px]",
                          m.confidence === "high" && "text-success",
                          m.confidence === "medium" && "text-warning",
                          m.confidence === "low" && "text-muted-foreground",
                        )}
                      >
                        {m.confidence} match
                      </span>
                      <Button
                        variant={m.confidence === "high" ? "primary" : "secondary"}
                        loading={reconcile.isPending && reconcile.variables?.invoice_id === m.invoice_id}
                        onClick={() =>
                          reconcile.mutate({ invoice_id: m.invoice_id, confidence: m.confidence })
                        }
                      >
                        Link payment
                      </Button>
                    </div>
                  </li>
                ))}
              </ul>
            )}

            <div className="flex justify-end">
              <Button type="button" variant="ghost" onClick={reset}>
                Paste different text
              </Button>
            </div>
          </div>
        )}
      </form>
    </Modal>
  );
}

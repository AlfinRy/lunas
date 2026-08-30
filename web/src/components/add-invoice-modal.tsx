import { useState, type FormEvent, type ReactNode } from "react";
import { Modal } from "./modal";
import { Button } from "./button";
import { useClients, useRefresh } from "@/api/queries";
import { api, ApiError } from "@/api/client";
import { useToast } from "./toast";

/**
 * F2.1 add invoice. Five fields, one modal; validation errors arrive per-field
 * from the API and render next to where they broke.
 */
export function AddInvoiceModal({ open, onClose }: { open: boolean; onClose: () => void }) {
  const clients = useClients();
  const refresh = useRefresh();
  const { toast } = useToast();
  const [saving, setSaving] = useState(false);
  const [errors, setErrors] = useState<Record<string, string[]>>({});

  async function onSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const fd = new FormData(e.currentTarget);
    const amount = Math.round(parseFloat(String(fd.get("amount") || "0")) * 100);
    setSaving(true);
    setErrors({});
    try {
      const inv = await api.createInvoice({
        client_id: Number(fd.get("client_id")),
        number: String(fd.get("number") || "").trim(),
        amount_cents: amount,
        currency: String(fd.get("currency") || "USD").toUpperCase(),
        issued_on: String(fd.get("issued_on")),
        due_on: String(fd.get("due_on")),
        notes: String(fd.get("notes") || "") || undefined,
      });
      toast("success", `Invoice ${inv.number} added — the agent starts planning its chase.`);
      await refresh();
      onClose();
    } catch (err) {
      if (err instanceof ApiError && err.fields) setErrors(err.fields);
      else toast("error", err instanceof Error ? err.message : "Unable to add the invoice. Try again.");
    } finally {
      setSaving(false);
    }
  }

  const today = new Date().toISOString().slice(0, 10);

  return (
    <Modal open={open} onClose={onClose} title="Add invoice">
      <form onSubmit={onSubmit} className="flex flex-col gap-4" noValidate>
        <Field label="Payer" error={errors["client_id"]}>
          {(id, invalid) => (
            <select
              id={id}
              name="client_id"
              aria-invalid={invalid}
              required
              className={inputCls}
              defaultValue=""
            >
              <option value="" disabled>
                Pick the payer
              </option>
              {(clients.data ?? []).map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name}
                </option>
              ))}
            </select>
          )}
        </Field>

        <div className="grid grid-cols-2 gap-3">
          <Field label="Invoice number" error={errors["number"]}>
            {(id, invalid) => (
              <input id={id} name="number" aria-invalid={invalid} placeholder="INV-0046" className={inputCls} required />
            )}
          </Field>
          <Field label="Amount" error={errors["amount_cents"]}>
            {(id, invalid) => (
              <input
                id={id}
                name="amount"
                inputMode="decimal"
                aria-invalid={invalid}
                placeholder="2400.00"
                className={`${inputCls} tnum`}
                required
              />
            )}
          </Field>
        </div>

        <div className="grid grid-cols-2 gap-3">
          <Field label="Issued on" error={undefined}>
            {(id) => <input id={id} name="issued_on" type="date" defaultValue={today} className={inputCls} required />}
          </Field>
          <Field label="Due on" error={errors["due_on"]}>
            {(id, invalid) => (
              <input
                id={id}
                name="due_on"
                type="date"
                aria-invalid={invalid}
                defaultValue={new Date(Date.now() + 14 * 86400_000).toISOString().slice(0, 10)}
                className={inputCls}
                required
              />
            )}
          </Field>
        </div>

        <Field label="Currency" error={errors["currency"]}>
          {(id, invalid) => (
            <input
              id={id}
              name="currency"
              aria-invalid={invalid}
              defaultValue="USD"
              maxLength={3}
              className={`${inputCls} uppercase`}
              required
            />
          )}
        </Field>

        <Field label="Notes (optional)">
          {(id) => <input id={id} name="notes" placeholder="e.g. Phase 2 — brand system" className={inputCls} />}
        </Field>

        <div className="mt-2 flex justify-end gap-2">
          <Button type="button" variant="ghost" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" variant="primary" loading={saving}>
            Add invoice
          </Button>
        </div>
      </form>
    </Modal>
  );
}

const inputCls =
  "w-full rounded-lg border border-input bg-card px-3 py-2 text-[16px] text-foreground placeholder:text-muted-foreground focus-visible:outline-2 focus-visible:outline-ring sm:text-sm";

/** Label always visible; placeholder is only an example (ux-writing §7). */
export function Field({
  label,
  error,
  children,
}: {
  label: string;
  error?: string[];
  children: (id: string, invalid: boolean) => ReactNode;
}) {
  const id = `f-${label.toLowerCase().replace(/[^a-z0-9]+/g, "-")}`;
  const invalid = !!error?.length;
  return (
    <div className="flex flex-col gap-1.5">
      <label htmlFor={id} className="text-sm">
        {label}
      </label>
      {children(id, invalid)}
      {invalid ? (
        <p className="text-xs text-destructive" role="alert">
          {error![0]}
        </p>
      ) : null}
    </div>
  );
}

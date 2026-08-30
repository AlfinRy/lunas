import { useState } from "react";
import { createRoute } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Route as rootRoute } from "./__root";
import { Card } from "@/components/card";
import { Button } from "@/components/button";
import { Modal } from "@/components/modal";
import { Field } from "@/components/add-invoice-modal";
import { ModeSwitch } from "@/components/mode-switch";
import { useToast } from "@/components/toast";
import { api, ApiError } from "@/api/client";

const inputCls =
  "w-full rounded-lg border border-input bg-card px-3 py-2 text-[16px] text-foreground placeholder:text-muted-foreground focus-visible:outline-2 focus-visible:outline-ring sm:text-sm";

function SettingsPage() {
  const qc = useQueryClient();
  const { toast } = useToast();
  const settings = useQuery({ queryKey: ["settings"], queryFn: api.settings });
  const [confirmOpen, setConfirmOpen] = useState(false);

  const save = useMutation({
    mutationFn: api.updateSettings,
    onSuccess: (s) => {
      toast("success", "Settings saved.");
      qc.setQueryData(["settings"], s);
      qc.invalidateQueries();
    },
    onError: (err) => toast("error", err instanceof ApiError ? err.message : "Unable to save. Try again."),
  });

  const reset = useMutation({
    mutationFn: api.resetDemo,
    onSuccess: () => {
      toast("success", "Demo data reset.");
      qc.invalidateQueries();
      setConfirmOpen(false);
    },
    onError: (err) => toast("error", err instanceof ApiError ? err.message : "Unable to reset. Try again."),
  });

  const s = settings.data;

  return (
    <main className="flex w-full max-w-[640px] flex-col gap-6">
      <header>
        <h1 className="text-[26px] font-medium leading-tight tracking-tight">Settings</h1>
        <p className="mt-0.5 text-sm text-muted-foreground">How Lunas signs, waits, and escalates.</p>
      </header>

      {settings.isPending || !s ? (
        <Card className="h-40 animate-pulse" />
      ) : (
        <>
          <Card className="flex flex-col gap-4 p-5">
            <h2 className="text-[15px] font-medium tracking-tight">Sender identity</h2>
            <form
              className="flex flex-col gap-4"
              onSubmit={(e) => {
                e.preventDefault();
                const fd = new FormData(e.currentTarget);
                save.mutate({
                  sender_name: String(fd.get("sender_name") || ""),
                  sender_email: String(fd.get("sender_email") || ""),
                  default_terms_days: Number(fd.get("default_terms_days") || 14),
                });
              }}
            >
              <Field label="Your name">
                {(id) => <input id={id} name="sender_name" defaultValue={s.sender_name} className={inputCls} />}
              </Field>
              <Field label="Billing email">
                {(id) => (
                  <input id={id} name="sender_email" type="email" defaultValue={s.sender_email} className={inputCls} />
                )}
              </Field>
              <Field label="Default payment terms (days)">
                {(id) => (
                  <input
                    id={id}
                    name="default_terms_days"
                    type="number"
                    min={0}
                    defaultValue={s.default_terms_days}
                    className={inputCls}
                  />
                )}
              </Field>
              <div className="flex justify-end">
                <Button type="submit" variant="primary" loading={save.isPending}>
                  Save changes
                </Button>
              </div>
            </form>
          </Card>

          <Card className="flex flex-wrap items-center justify-between gap-4 p-5">
            <div className="min-w-0">
              <h2 className="text-[15px] font-medium tracking-tight">Agent mode</h2>
              <p className="mt-0.5 text-sm text-muted-foreground">
                Full auto sends gentle through firm chases; final notices always wait for you.
              </p>
            </div>
            <ModeSwitch />
          </Card>

          <Card className="flex flex-col gap-3 p-5">
            <h2 className="text-[15px] font-medium tracking-tight text-destructive">Danger zone</h2>
            <p className="text-sm text-muted-foreground">
              Reset demo data. This deletes all invoices, drafts, sent emails, and payments, then reseeds the demo
              story.
            </p>
            <div className="flex justify-end">
              <Button variant="secondary" className="text-destructive" onClick={() => setConfirmOpen(true)}>
                Reset demo data
              </Button>
            </div>
          </Card>
        </>
      )}

      <Modal open={confirmOpen} onClose={() => setConfirmOpen(false)} title="Reset demo data?">
        <p className="text-sm text-muted-foreground">
          This deletes every invoice, draft, sent email, and payment, then reseeds the demo story. There is no undo.
        </p>
        <div className="mt-5 flex justify-end gap-2">
          <Button variant="ghost" onClick={() => setConfirmOpen(false)}>
            Cancel
          </Button>
          <Button
            variant="primary"
            className="bg-destructive hover:bg-destructive/85"
            loading={reset.isPending}
            onClick={() => reset.mutate()}
          >
            Reset demo data
          </Button>
        </div>
      </Modal>
    </main>
  );
}

export const Route = createRoute({
  getParentRoute: () => rootRoute,
  path: "/settings",
  component: SettingsPage,
});

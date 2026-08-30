import { createRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { BotMessageSquare, CalendarClock, Info } from "lucide-react";
import { Route as rootRoute } from "./__root";
import { Card } from "@/components/card";
import { Button } from "@/components/button";
import { DemoTimePill } from "@/components/demo-time-pill";
import { ModeSwitch } from "@/components/mode-switch";
import { useDraftActions, stageInfo } from "@/api/agent";
import { api } from "@/api/client";
import { money } from "@/lib/format";
import { cn } from "@/lib/cn";

function AgentPage() {
  const inbox = useQuery({ queryKey: ["agent-inbox"], queryFn: api.agentInbox });
  const settings = useQuery({ queryKey: ["settings"], queryFn: api.settings, staleTime: 10_000 });
  const { approve, skip } = useDraftActions();

  const drafts = inbox.data?.drafts ?? [];
  const plans = inbox.data?.plans ?? [];

  return (
    <main className="mx-auto flex max-w-[900px] flex-col gap-6">
      <header className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="text-[26px] font-medium leading-tight tracking-tight">Agent inbox</h1>
          <p className="mt-0.5 text-sm text-muted-foreground">
            What Lunas is planning, and what waits for your approval.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <ModeSwitch />
          <DemoTimePill />
        </div>
      </header>

      {settings.data?.template_mode ? (
        <Card className="flex items-start gap-2.5 px-5 py-3.5 text-sm text-muted-foreground">
          <Info size={15} strokeWidth={1.5} className="mt-0.5 shrink-0 text-info" aria-hidden="true" />
          <p>
            Template mode — no AI provider configured, so chases use the standard ladder copy. Set{" "}
            <code className="font-mono text-xs">LLM_API_KEY</code> to enable model-drafted emails. Decisions run either way.
          </p>
        </Card>
      ) : null}

      <section aria-label="Drafts awaiting approval" className="flex flex-col gap-4">
        <h2 className="text-[15px] font-medium tracking-tight">
          Awaiting approval {drafts.length ? <span className="text-muted-foreground">({drafts.length})</span> : null}
        </h2>

        {inbox.isPending ? (
          <Card className="h-32 animate-pulse" />
        ) : drafts.length === 0 ? (
          <Card className="flex flex-col items-start gap-1 px-5 py-8">
            <p className="flex items-center gap-2 font-medium">
              <BotMessageSquare size={16} strokeWidth={1.5} aria-hidden="true" /> Nothing needs your approval
            </p>
            <p className="text-sm text-muted-foreground">
              {plans.length
                ? `Lunas is watching the plans below. Drafts appear here the moment a chase is due.`
                : "Lunas is watching your invoices. You'll see drafts here the moment one needs chasing."}
            </p>
          </Card>
        ) : (
          drafts.map((d) => {
            const meta = stageInfo(d.stage);
            return (
              <Card key={d.id} className="flex flex-col gap-4 p-5">
                <div className="flex flex-wrap items-start justify-between gap-2">
                  <div>
                    <p className="flex items-center gap-2 font-medium">
                      <span aria-hidden="true" className={cn("size-1.5 rounded-full", meta.dot)} />
                      {d.client_name}
                      <span className="font-mono text-xs text-muted-foreground">{d.invoice_number}</span>
                    </p>
                    <p className="mt-0.5 text-xs text-muted-foreground">
                      {meta.label} · {meta.tone} tone · to {d.client_email}
                    </p>
                  </div>
                  {d.stage === "stage3_final" ? (
                    <span className="rounded-pill border border-border px-2.5 py-1 text-xs text-destructive">
                      Needs your send — final notices never auto-send
                    </span>
                  ) : null}
                </div>

                <div className="rounded-2xl border border-border bg-background px-4 py-3">
                  <p className="text-sm">{d.subject}</p>
                  <p className="mt-2 line-clamp-4 whitespace-pre-line text-sm text-muted-foreground">{d.body}</p>
                </div>

                <div className="flex justify-end gap-2">
                  <Button variant="ghost" loading={skip.isPending && skip.variables === d.id} onClick={() => skip.mutate(d.id)}>
                    Skip this draft
                  </Button>
                  <Button
                    variant="primary"
                    loading={approve.isPending && approve.variables === d.id}
                    onClick={() => approve.mutate(d.id)}
                  >
                    {d.stage === "stage4_escalation" ? "Open dossier" : "Approve & send"}
                  </Button>
                </div>
              </Card>
            );
          })
        )}
      </section>

      <section aria-label="Planned chases" className="flex flex-col gap-4">
        <h2 className="text-[15px] font-medium tracking-tight">
          Planned {plans.length ? <span className="text-muted-foreground">({plans.length})</span> : null}
        </h2>
        {plans.length === 0 ? (
          <Card className="px-5 py-6 text-sm text-muted-foreground">Nothing scheduled — every open invoice has its next step queued or is awaiting approval above.</Card>
        ) : (
          plans.map((p) => {
            const meta = stageInfo(p.stage);
            return (
              <Card key={p.invoice_id} className="flex items-start gap-3.5 p-5">
                <span aria-hidden="true" className={cn("mt-1.5 size-1.5 shrink-0 rounded-full", meta.dot)} />
                <div className="min-w-0 flex-1">
                  <p className="text-sm">
                    <span className="font-medium">{p.client_name}</span>{" "}
                    <span className="font-mono text-xs text-muted-foreground">{p.invoice_number}</span> ·{" "}
                    <span className="tnum">{money(p.amount_cents)}</span>
                  </p>
                  <p className="mt-1 text-sm text-muted-foreground">{p.reasoning}</p>
                </div>
                <p className="flex shrink-0 items-center gap-1.5 text-xs text-muted-foreground">
                  <CalendarClock size={13} strokeWidth={1.5} aria-hidden="true" />
                  {fmtDate(p.planned_on)}
                </p>
              </Card>
            );
          })
        )}
      </section>
    </main>
  );
}

function fmtDate(iso: string): string {
  return new Date(`${iso}T00:00:00Z`).toLocaleDateString("en-US", { month: "short", day: "numeric", timeZone: "UTC" });
}

export const Route = createRoute({
  getParentRoute: () => rootRoute,
  path: "/agent",
  component: AgentPage,
});

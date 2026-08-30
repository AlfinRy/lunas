import { useState } from "react";
import { createRoute } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { Mail } from "lucide-react";
import { Route as rootRoute } from "./__root";
import { Card } from "@/components/card";
import { api } from "@/api/client";
import { cn } from "@/lib/cn";

/** The outbox reads like a mail client: list on the left, reading pane right. */
function OutboxPage() {
  const outbox = useQuery({ queryKey: ["outbox"], queryFn: api.outbox });
  const emails = outbox.data ?? [];
  const [selected, setSelected] = useState<number | null>(null);
  const active = emails.find((e) => e.id === selected) ?? emails[0];

  return (
    <main className="flex flex-col gap-6">
      <header>
        <h1 className="text-[26px] font-medium leading-tight tracking-tight">Outbox</h1>
        <p className="mt-0.5 text-sm text-muted-foreground">Every chase Lunas has sent, newest first.</p>
      </header>

      {outbox.isPending ? (
        <Card className="h-40 animate-pulse" />
      ) : emails.length === 0 ? (
        <Card className="flex flex-col items-start gap-1 px-5 py-8">
          <p className="flex items-center gap-2 font-medium">
            <Mail size={16} strokeWidth={1.5} aria-hidden="true" /> No emails sent yet
          </p>
          <p className="text-sm text-muted-foreground">Approved chases land here with their delivery status.</p>
        </Card>
      ) : (
        <div className="flex flex-col gap-4 lg:flex-row">
          <Card className="min-w-0 flex-1 p-0">
            <ul>
              {emails.map((e) => (
                <li key={e.id}>
                  <button
                    onClick={() => setSelected(e.id)}
                    aria-current={active?.id === e.id}
                    className={cn(
                      "flex w-full flex-col items-start gap-0.5 border-t border-border px-4 py-3 text-start transition-colors first:border-t-0 hover:bg-accent-wash",
                      active?.id === e.id && "bg-primary-soft/60",
                    )}
                  >
                    <span className="flex w-full items-baseline justify-between gap-2">
                      <span className="truncate text-sm font-medium">{e.to_name}</span>
                      <span className="shrink-0 text-xs text-muted-foreground">{fmt(e.sent_at)}</span>
                    </span>
                    <span className="line-clamp-1 text-xs text-muted-foreground">{e.subject}</span>
                  </button>
                </li>
              ))}
            </ul>
          </Card>

          {active ? (
            <Card className="min-w-0 flex-1 p-5">
              <p className="text-xs text-muted-foreground">
                To {active.to_name} &lt;{active.to_email}&gt; · {active.invoice_number}
              </p>
              <h2 className="mt-2 text-[15px] font-medium tracking-tight">{active.subject}</h2>
              <p className="mt-4 whitespace-pre-line text-sm leading-relaxed">{active.body}</p>
            </Card>
          ) : null}
        </div>
      )}
    </main>
  );
}

function fmt(ts: string): string {
  return new Date(ts).toLocaleDateString("en-US", { month: "short", day: "numeric" });
}

export const Route = createRoute({
  getParentRoute: () => rootRoute,
  path: "/outbox",
  component: OutboxPage,
});

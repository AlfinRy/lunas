import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { CalendarClock, ChevronDown } from "lucide-react";
import { motion, AnimatePresence } from "motion/react";
import { api } from "@/api/client";
import { useSimulate } from "@/api/agent";
import { BANNER } from "@/lib/springs";
import { cn } from "@/lib/cn";

/**
 * OA 09-floating-pill: the demo clock is the app chrome's one floating word.
 * Its controls expand in place with BANNER.
 */
export function DemoTimePill() {
  const settings = useQuery({ queryKey: ["settings"], queryFn: api.settings, staleTime: 10_000 });
  const simulate = useSimulate();
  const [open, setOpen] = useState(false);

  const simNow = settings.data?.sim_now;
  const label = simNow ? `Demo time: ${fmt(simNow)}` : "Real time";

  return (
    <div className="relative">
      <button
        onClick={() => setOpen((v) => !v)}
        aria-expanded={open}
        className={cn(
          "inline-flex items-center gap-1.5 rounded-pill border border-border bg-card px-3 py-1.5 text-xs transition-colors hover:bg-accent-wash",
          simNow && "text-warning",
        )}
      >
        <CalendarClock size={13} strokeWidth={1.5} aria-hidden="true" />
        {label}
        <ChevronDown
          size={12}
          strokeWidth={1.5}
          aria-hidden="true"
          className={cn("transition-transform duration-150", open && "rotate-180")}
        />
      </button>

      <AnimatePresence>
        {open ? (
          <motion.div
            role="menu"
            aria-label="Time controls"
            initial={{ opacity: 0, y: -6, scale: 0.98 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: -4, scale: 0.98 }}
            transition={BANNER}
            className="absolute end-0 top-[calc(100%+6px)] z-30 flex w-44 flex-col gap-1 rounded-2xl border border-border bg-popover p-2 shadow-lg"
          >
            <p className="px-2 py-1 text-xs text-muted-foreground">Advance the demo clock</p>
            {[
              { days: 1, label: "+1 day" },
              { days: 3, label: "+3 days" },
              { days: 7, label: "+7 days" },
              { days: 14, label: "+14 days" },
            ].map((o) => (
              <button
                key={o.days}
                role="menuitem"
                disabled={simulate.isPending}
                onClick={() => {
                  setOpen(false);
                  simulate.mutate({ days: o.days });
                }}
                className="rounded-pill px-3 py-1.5 text-left text-sm transition-colors hover:bg-accent-wash disabled:opacity-50"
              >
                {o.label}
              </button>
            ))}
            <p className="px-2 pb-1 pt-2 text-[11px] text-muted-foreground">
              The agent acts as of the new time — drafts and sends appear here.
            </p>
          </motion.div>
        ) : null}
      </AnimatePresence>
    </div>
  );
}

function fmt(iso: string): string {
  return new Date(`${iso}T00:00:00Z`).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    timeZone: "UTC",
  });
}

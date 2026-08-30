import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { motion } from "motion/react";
import { api, ApiError } from "@/api/client";
import { useToast } from "./toast";
import { LAYOUT } from "@/lib/springs";

/**
 * Mode switcher with the OA sliding highlight: one shared pill travels
 * between Approve each / Full auto (LAYOUT spring), never a repaint.
 */
const modes = [
  { key: "approve_each", label: "Approve each" },
  { key: "full_auto", label: "Full auto" },
] as const;

export function ModeSwitch() {
  const qc = useQueryClient();
  const { toast } = useToast();
  const settings = useQuery({ queryKey: ["settings"], queryFn: api.settings, staleTime: 10_000 });

  const update = useMutation({
    mutationFn: (global_mode: (typeof modes)[number]["key"]) => api.updateSettings({ global_mode }),
    onSuccess: (s) => {
      toast(
        "success",
        s.global_mode === "full_auto"
          ? "Full auto on — Lunas sends gentle through firm chases on its own. Final notices still wait for you."
          : "Approve-each on — every chase waits for your approval.",
      );
      qc.setQueryData(["settings"], s);
      qc.invalidateQueries();
    },
    onError: (err) => toast("error", err instanceof ApiError ? err.message : "Unable to save. Try again."),
  });

  const current = settings.data?.global_mode ?? "approve_each";

  return (
    <div role="radiogroup" aria-label="Agent mode" className="relative flex rounded-pill border border-border bg-card p-0.5">
      {modes.map((m) => {
        const active = current === m.key;
        return (
          <button
            key={m.key}
            role="radio"
            aria-checked={active}
            disabled={update.isPending}
            onClick={() => update.mutate(m.key)}
            className="relative rounded-pill px-3 py-1 text-xs transition-colors disabled:opacity-60"
          >
            {active ? (
              <motion.span
                layoutId="mode-highlight"
                transition={LAYOUT}
                className="absolute inset-0 rounded-pill bg-primary-soft"
              />
            ) : null}
            <span className={active ? "relative text-foreground" : "relative text-muted-foreground"}>{m.label}</span>
          </button>
        );
      })}
    </div>
  );
}

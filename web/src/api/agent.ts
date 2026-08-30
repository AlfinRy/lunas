import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api, ApiError } from "./client";
import { useToast } from "@/components/toast";
import type { ChaseStage } from "./client";


export const stageMeta: Record<string, { label: string; tone: string; dot: string }> = {
  stage0_reminder: { label: "Reminder", tone: "gentle", dot: "bg-primary" },
  stage1_polite: { label: "Polite follow-up", tone: "polite", dot: "bg-info" },
  stage2_firm: { label: "Firm follow-up", tone: "firm", dot: "bg-warning" },
  stage3_final: { label: "Final notice", tone: "final", dot: "bg-destructive" },
  stage4_escalation: { label: "Collections dossier", tone: "escalation", dot: "bg-destructive" },
};

export function stageInfo(s: ChaseStage | string) {
  return stageMeta[s] ?? { label: s, tone: "", dot: "bg-muted-foreground" };
}

/** Shared approve/skip mutations for drafts. */
export function useDraftActions() {
  const qc = useQueryClient();
  const { toast } = useToast();
  const invalidate = () => {
    qc.invalidateQueries({ queryKey: ["agent-inbox"] });
    qc.invalidateQueries({ queryKey: ["invoices"] });
    qc.invalidateQueries({ queryKey: ["dashboard"] });
    qc.invalidateQueries({ queryKey: ["outbox"] });
    qc.invalidateQueries({ queryKey: ["activity"] });
    qc.invalidateQueries({ queryKey: ["invoice"] });
  };

  const approve = useMutation({
    mutationFn: api.approveDraft,
    onSuccess: (sent) => {
      toast("success", `Chase sent to ${sent.to_name}.`);
      invalidate();
    },
    onError: (err) => toast("error", err instanceof ApiError ? err.message : "Unable to send. Try again."),
  });

  const skip = useMutation({
    mutationFn: api.skipDraft,
    onSuccess: (d) => {
      toast("success", `Skipped the ${stageInfo(d.stage).label.toLowerCase()} — Lunas will replan.`);
      invalidate();
    },
    onError: (err) => toast("error", err instanceof ApiError ? err.message : "Unable to skip. Try again."),
  });

  return { approve, skip };
}

export function useSimulate() {
  const qc = useQueryClient();
  const { toast } = useToast();
  return useMutation({
    mutationFn: api.simulateAdvance,
    onSuccess: (s) => {
      toast("success", `Demo clock advanced to ${s.sim_now}. The agent ran as of the new time.`);
      qc.invalidateQueries();
    },
    onError: (err) => toast("error", err instanceof ApiError ? err.message : "Unable to advance time. Try again."),
  });
}

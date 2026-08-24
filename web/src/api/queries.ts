import { useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "./client";

export const qk = {
  dashboard: ["dashboard"] as const,
  invoices: ["invoices"] as const,
  clients: ["clients"] as const,
  settings: ["settings"] as const,
};

export function useDashboard() {
  return useQuery({ queryKey: qk.dashboard, queryFn: api.dashboard });
}
export function useInvoices() {
  return useQuery({ queryKey: qk.invoices, queryFn: api.invoices });
}
export function useClients() {
  return useQuery({ queryKey: qk.clients, queryFn: api.clients });
}
export function useSettings() {
  return useQuery({ queryKey: qk.settings, queryFn: api.settings });
}

/** Invalidate everything the demo touches — called after any mutation. */
export function useRefresh() {
  const qc = useQueryClient();
  return () => qc.invalidateQueries();
}

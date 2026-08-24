import type { components } from "./schema";

export type Dashboard = components["schemas"]["Dashboard"];
export type Invoice = components["schemas"]["Invoice"];
export type Client = components["schemas"]["Client"];
export type Activity = components["schemas"]["Activity"];
export type Settings = components["schemas"]["Settings"];
export type InvoiceStatus = components["schemas"]["InvoiceStatus"];

export class ApiError extends Error {
  fields?: Record<string, string[]>;
  status: number;
  constructor(status: number, message: string, fields?: Record<string, string[]>) {
    super(message);
    this.status = status;
    this.fields = fields;
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    headers: { "Content-Type": "application/json" },
    ...init,
  });
  if (!res.ok) {
    let message = `Request failed (${res.status}).`;
    let fields: Record<string, string[]> | undefined;
    try {
      const body = await res.json();
      if (body?.message) message = body.message;
      if (body?.fields) fields = body.fields;
    } catch {
      /* non-JSON error body */
    }
    throw new ApiError(res.status, message, fields);
  }
  if (res.status === 204) return undefined as T;
  return res.json();
}

export const api = {
  dashboard: () => request<Dashboard>("/api/dashboard"),
  invoices: () => request<Invoice[]>("/api/invoices"),
  invoice: (id: number) => request<Invoice>(`/api/invoices/${id}`),
  createInvoice: (body: components["schemas"]["InvoiceInput"]) =>
    request<Invoice>("/api/invoices", { method: "POST", body: JSON.stringify(body) }),
  activity: (id: number) => request<Activity[]>(`/api/invoices/${id}/activity`),
  clients: () => request<Client[]>("/api/clients"),
  settings: () => request<Settings>("/api/settings"),
};

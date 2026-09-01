import { useState } from "react";
import { Plus } from "lucide-react";
import { createRoute } from "@tanstack/react-router";
import { Route as rootRoute } from "./__root";
import { Plate } from "@/components/card";
import { Button } from "@/components/button";
import { InvoiceTable } from "@/components/invoice-table";
import { AddInvoiceModal } from "@/components/add-invoice-modal";
import { useInvoices } from "@/api/queries";

function InvoicesPage() {
  const invoices = useInvoices();
  const [adding, setAdding] = useState(false);

  return (
    <main className="flex flex-col gap-6">
      <header className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="text-[26px] font-medium leading-tight tracking-tight">Invoices</h1>
          <p className="mt-0.5 text-sm text-muted-foreground">
            Every outstanding and settled invoice, oldest due first.
          </p>
        </div>
        <Button variant="primary" onClick={() => setAdding(true)}>
          <Plus size={15} strokeWidth={1.5} aria-hidden="true" />
          Add invoice
        </Button>
      </header>

      <Plate title="All invoices" meta="Overdue first" contentClassName="py-0">
        <InvoiceTable invoices={invoices.data} loading={invoices.isPending} skeletonRows={6} />
        {invoices.data && invoices.data.length === 0 ? (
          <div className="flex flex-col items-start gap-1 px-5 py-10">
            <p className="font-medium">No invoices yet</p>
            <p className="text-sm text-muted-foreground">
              Add an invoice and Lunas will start planning the first chase.
            </p>
            <Button variant="primary" className="mt-3" onClick={() => setAdding(true)}>
              Add invoice
            </Button>
          </div>
        ) : null}
      </Plate>

      <AddInvoiceModal open={adding} onClose={() => setAdding(false)} />
    </main>
  );
}

export const Route = createRoute({
  getParentRoute: () => rootRoute,
  path: "/invoices",
  component: InvoicesPage,
});

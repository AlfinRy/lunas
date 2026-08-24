import { createRoute } from "@tanstack/react-router";
import { Route as rootRoute } from "./__root";

function Page() {
  return (
    <main className="mx-auto flex max-w-[1200px] flex-col gap-6">
      <h1 className="text-[26px] font-medium tracking-tight">Invoices</h1>
      <p className="text-sm text-muted-foreground">Ships in week 2 of the build plan.</p>
    </main>
  );
}

export const Route = createRoute({ getParentRoute: () => rootRoute, path: "/invoices", component: Page });

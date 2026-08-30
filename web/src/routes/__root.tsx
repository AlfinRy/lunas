import { createRootRouteWithContext, Link, Outlet } from "@tanstack/react-router";
import { QueryClient } from "@tanstack/react-query";
import {
  LayoutDashboard,
  FileText,
  Users,
  BotMessageSquare,
  Send,
  Settings as SettingsIcon,
} from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/api/client";

type RouterContext = { queryClient: QueryClient };

const nav = [
  { to: "/", label: "Dashboard", icon: LayoutDashboard },
  { to: "/invoices", label: "Invoices", icon: FileText },
  { to: "/clients", label: "Clients", icon: Users },
  { to: "/agent", label: "Agent inbox", icon: BotMessageSquare, badge: "agent" },
  { to: "/outbox", label: "Outbox", icon: Send },
  { to: "/settings", label: "Settings", icon: SettingsIcon },
] as const;


function Shell() {
  return (
    <div className="min-h-dvh bg-background">
      <a
        href="#content"
        className="sr-only focus:not-sr-only focus:absolute focus:z-50 focus:m-2 focus:rounded-pill focus:bg-card focus:px-4 focus:py-2 focus:text-sm"
      >
        Skip to content
      </a>
      <div className="mx-auto flex min-h-dvh w-full max-w-[1440px]">
        <nav
          aria-label="Main"
          className="sticky top-0 hidden h-dvh w-[200px] shrink-0 flex-col gap-1 border-e border-border px-3 py-4 md:flex"
        >
          <Link to="/" className="mb-4 flex items-center gap-2 px-2 py-1">
            <span className="flex size-7 items-center justify-center rounded-[8px] bg-primary text-primary-foreground">
              <svg width="15" height="15" viewBox="0 0 16 16" fill="none" aria-hidden="true">
                <path
                  d="M3 8.5 6.2 11.5 13 4.5"
                  stroke="currentColor"
                  strokeWidth="2.2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            </span>
            <span className="text-[15px] font-medium tracking-tight">Lunas</span>
          </Link>
          {nav.map((item) => {
            const Icon = item.icon;
            const badge = "badge" in item ? item.badge : undefined;
            return (
            <Link
              key={item.to}
              to={item.to}
              className="group flex items-center gap-2.5 rounded-pill px-3 py-1.5 text-[13.5px] text-muted-foreground transition-colors duration-150 ease-out hover:bg-accent-wash hover:text-foreground aria-[current=page]:bg-primary-soft aria-[current=page]:text-foreground"
            >
              <Icon size={16} strokeWidth={1.5} className="shrink-0" aria-hidden="true" />
              {item.label}
              {badge === "agent" ? <AgentBadge /> : null}
            </Link>
            );
          })}
          <p className="mt-auto px-3 text-xs text-muted-foreground">
            The AI collections agent
            <br />
            that gets you paid.
          </p>
        </nav>

        <div id="content" className="min-w-0 flex-1 px-4 pb-16 pt-6 md:px-8">
          <Outlet />
        </div>
      </div>
    </div>
  );
}

function AgentBadge() {
  const dash = useQuery({ queryKey: ["dashboard"], queryFn: api.dashboard, staleTime: 10_000 });
  const n = dash.data?.counts.awaiting_approval ?? 0;
  if (!n) return null;
  return (
    <span
      aria-label={`${n} awaiting approval`}
      className="tnum ms-auto rounded-pill bg-primary px-1.5 py-0.5 text-[10px] leading-none text-primary-foreground"
    >
      {n}
    </span>
  );
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: Shell,
  notFoundComponent: () => (
    <div className="py-24 text-center">
      <p className="font-medium">Page not found</p>
      <Link to="/" className="mt-2 inline-block text-info underline underline-offset-4">
        Back to the dashboard
      </Link>
    </div>
  ),
});

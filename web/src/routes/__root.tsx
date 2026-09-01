import { createRootRouteWithContext, Link, Outlet } from "@tanstack/react-router";
import { QueryClient, useQuery } from "@tanstack/react-query";
import { AnimatePresence, motion } from "motion/react";
import {
  LayoutDashboard,
  FileText,
  Users,
  BotMessageSquare,
  Send,
  Settings as SettingsIcon,
  Moon,
  Sun,
} from "lucide-react";
import { api } from "@/api/client";
import { cn } from "@/lib/cn";
import { useTheme } from "@/lib/theme";
import { FLICK, LAYOUT } from "@/lib/springs";

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
      <MobileTopNav />
      <div className="mx-auto flex min-h-dvh w-full max-w-[1440px]">
        <SideNav />

        {/* One consistent content column: same width, same start edge on every page. */}
        <div id="content" className="mx-auto min-w-0 w-full max-w-[1200px] flex-1 px-4 pb-16 pt-6 md:px-8">
          <Outlet />
        </div>
      </div>
    </div>
  );
}

/** Desktop sidebar with the OA traveling highlight: one shared pill, LAYOUT spring. */
function SideNav() {
  return (
    <nav
      aria-label="Main"
      className="sticky top-0 hidden h-dvh w-[200px] shrink-0 flex-col gap-1 border-e border-border px-3 py-4 md:flex"
    >
      <Link to="/" className="mb-4 flex items-center gap-2 px-2 py-1">
        <span className="flex size-7 items-center justify-center rounded-[9px] bg-primary text-primary-foreground shadow-[inset_0_1px_0_rgba(255,255,255,0.25)]">
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

      <NavLinks />

      <div className="mt-auto flex items-center justify-between px-1">
        <p className="text-xs leading-relaxed text-muted-foreground">
          The AI collections agent
          <br />
          that gets you paid.
        </p>
        <ThemeToggle />
      </div>
    </nav>
  );
}

function NavLinks() {
  return (
    <>
      {nav.map((item) => {
        const Icon = item.icon;
        const badge = "badge" in item ? item.badge : undefined;
        return (
          <Link
            key={item.to}
            to={item.to}
            className="group relative flex items-center gap-2.5 rounded-pill px-3 py-1.5 text-[13.5px] text-muted-foreground transition-colors duration-150 ease-out hover:text-foreground aria-[current=page]:text-foreground"
          >
            {({ isActive }) => (
              <>
                {/* the traveling highlight: one shared object the eye follows */}
                {isActive ? (
                  <motion.span
                    layoutId="sidebar-highlight"
                    transition={LAYOUT}
                    aria-hidden="true"
                    className="absolute inset-0 rounded-pill bg-primary-soft"
                  />
                ) : (
                  <span
                    aria-hidden="true"
                    className="absolute inset-0 rounded-pill opacity-0 transition-colors duration-150 ease-out group-hover:bg-accent-wash group-hover:opacity-100"
                  />
                )}
                <Icon size={16} strokeWidth={1.5} className="relative shrink-0" aria-hidden="true" />
                <span className="relative">{item.label}</span>
                {badge === "agent" ? <AgentBadge /> : null}
              </>
            )}
          </Link>
        );
      })}
    </>
  );
}

/** Below md the sidebar is hidden; this sticky top bar keeps every page reachable. */
function MobileTopNav() {
  return (
    <header className="sticky top-0 z-30 border-b border-border bg-background/90 backdrop-blur md:hidden">
      <div className="flex items-center gap-3 px-4 py-2.5">
        <Link to="/" className="flex shrink-0 items-center gap-2">
          <span className="flex size-6 items-center justify-center rounded-[8px] bg-primary text-primary-foreground">
            <svg width="13" height="13" viewBox="0 0 16 16" fill="none" aria-hidden="true">
              <path
                d="M3 8.5 6.2 11.5 13 4.5"
                stroke="currentColor"
                strokeWidth="2.2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </span>
          <span className="text-sm font-medium tracking-tight">Lunas</span>
        </Link>
        <nav
          aria-label="Main"
          className="scrollbar-none -me-1 flex min-w-0 flex-1 items-center gap-1 overflow-x-auto px-1"
        >
          {nav.map((item) => {
            const Icon = item.icon;
            const badge = "badge" in item ? item.badge : undefined;
            return (
              <Link
                key={item.to}
                to={item.to}
                className="flex shrink-0 items-center gap-1.5 rounded-pill px-2.5 py-1 text-xs text-muted-foreground transition-colors duration-150 ease-out hover:bg-accent-wash hover:text-foreground aria-[current=page]:bg-primary-soft aria-[current=page]:text-foreground"
              >
                <Icon size={13} strokeWidth={1.5} className="shrink-0" aria-hidden="true" />
                {item.label}
                {badge === "agent" ? <AgentBadge compact /> : null}
              </Link>
            );
          })}
        </nav>
        <ThemeToggle compact />
      </div>
    </header>
  );
}

/** Sun/moon swap on FLICK (oa 07 contextual icon animation). */
function ThemeToggle({ compact = false }: { compact?: boolean }) {
  const { theme, toggle } = useTheme();
  return (
    <button
      onClick={toggle}
      aria-label={theme === "dark" ? "Switch to light appearance" : "Switch to dark appearance"}
      className={cn(
        "relative flex shrink-0 items-center justify-center rounded-pill text-muted-foreground transition-colors hover:bg-accent-wash hover:text-foreground",
        compact ? "size-7" : "size-8",
      )}
    >
      <AnimatePresence initial={false} mode="popLayout">
        <motion.span
          key={theme}
          initial={{ opacity: 0, scale: 0.25, filter: "blur(4px)" }}
          animate={{ opacity: 1, scale: 1, filter: "blur(0px)" }}
          exit={{ opacity: 0, scale: 0.25, filter: "blur(4px)" }}
          transition={FLICK}
          className="flex"
        >
          {theme === "dark" ? (
            <Sun size={compact ? 14 : 15} strokeWidth={1.5} aria-hidden="true" />
          ) : (
            <Moon size={compact ? 14 : 15} strokeWidth={1.5} aria-hidden="true" />
          )}
        </motion.span>
      </AnimatePresence>
    </button>
  );
}

function AgentBadge({ compact = false }: { compact?: boolean }) {
  const dash = useQuery({ queryKey: ["dashboard"], queryFn: api.dashboard, staleTime: 10_000 });
  const n = dash.data?.counts.awaiting_approval ?? 0;
  if (!n) return null;
  return (
    <span
      aria-label={`${n} awaiting approval`}
      className={cn(
        "tnum rounded-pill bg-primary px-1.5 leading-none text-primary-foreground shadow-[inset_0_1px_0_rgba(255,255,255,0.2)]",
        compact ? "py-px text-[9px]" : "py-0.5 text-[10px]",
        !compact && "ms-auto",
      )}
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

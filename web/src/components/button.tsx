import type { ButtonHTMLAttributes, ReactNode } from "react";
import { Loader2 } from "lucide-react";
import { cn } from "@/lib/cn";

type Variant = "primary" | "secondary" | "ghost";

/**
 * OA 02-button: pills that press. One primary per view — the jade accent is
 * spent on the single most important action; secondaries stay flat grey.
 */
export function Button({
  variant = "secondary",
  loading = false,
  className,
  children,
  ...props
}: ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: Variant;
  loading?: boolean;
  children: ReactNode;
}) {
  return (
    <button
      {...props}
      disabled={props.disabled || loading}
      className={cn(
        "inline-flex items-center justify-center gap-2 rounded-pill px-4 py-1.5 text-sm transition-colors duration-150 ease-out",
        "active:scale-[0.96] transition-transform",
        "focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring",
        "disabled:pointer-events-none disabled:opacity-50",
        variant === "primary" &&
          "bg-primary text-primary-foreground hover:bg-primary-strong",
        variant === "secondary" &&
          "bg-secondary text-secondary-foreground hover:bg-accent-wash",
        variant === "ghost" &&
          "text-muted-foreground hover:bg-accent-wash hover:text-foreground",
        className,
      )}
    >
      {loading ? <Loader2 size={15} strokeWidth={1.5} className="animate-spin" aria-hidden="true" /> : null}
      {children}
    </button>
  );
}

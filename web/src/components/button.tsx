import type { ButtonHTMLAttributes, ReactNode } from "react";
import { cn } from "@/lib/cn";

/**
 * The OA button (02-button), as a pill (our house style — "everything
 * clickable that is not a card is a pill"):
 * - press = translate-y-px + scale-0.98 — the only pointer geometry change
 * - primary carries the bevel: fill mixed toward deep indigo at rest, pure
 *   jade on hover, white inset highlight on top, dark seat below
 * - loading KEEPS the label and adds a ring spinner; width never jumps
 */
type Variant = "primary" | "secondary" | "ghost" | "destructive";

const BASE =
  "inline-flex cursor-pointer items-center justify-center gap-1.5 whitespace-nowrap font-medium outline-none transition-all " +
  "focus-visible:ring-[3px] focus-visible:ring-ring/50 " +
  "disabled:pointer-events-none disabled:opacity-50 " +
  "active:translate-y-px active:scale-[0.98] " +
  "rounded-pill px-4 py-1.5 text-sm";

const VARIANTS: Record<Variant, string> = {
  primary:
    "border border-[color-mix(in_srgb,var(--primary)_80%,#3a3480)] " +
    "bg-[color-mix(in_srgb,var(--primary)_90%,#3a3480)] text-primary-foreground " +
    "shadow-[inset_0_1px_0_rgba(255,255,255,0.22),inset_0_-1px_0_rgba(24,38,30,0.35)] " +
    "hover:bg-primary hover:border-[color-mix(in_srgb,var(--primary)_70%,#3a3480)]",
  secondary:
    "border border-transparent bg-secondary text-secondary-foreground " +
    "hover:bg-[color-mix(in_srgb,var(--secondary)_95%,var(--ink))]",
  ghost: "text-muted-foreground hover:bg-accent-wash hover:text-foreground",
  destructive:
    "border border-transparent bg-destructive text-white " +
    "shadow-[inset_0_1px_0_rgba(255,255,255,0.18),inset_0_-1px_0_rgba(0,0,0,0.25)] hover:bg-destructive/90",
};

export function Spinner({ className }: { className?: string }) {
  return (
    <span
      aria-label="Loading"
      role="status"
      className={cn(
        "inline-block size-3.5 shrink-0 animate-spin rounded-full border-2 border-current border-t-transparent",
        className,
      )}
    />
  );
}

export function Button({
  variant = "secondary",
  loading = false,
  className,
  children,
  type,
  ...props
}: ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: Variant;
  loading?: boolean;
  children: ReactNode;
}) {
  return (
    <button
      {...props}
      type={type ?? "button"}
      disabled={props.disabled || loading}
      aria-disabled={loading || undefined}
      className={cn(BASE, VARIANTS[variant], className)}
    >
      {loading ? <Spinner /> : null}
      {children}
    </button>
  );
}

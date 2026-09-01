import { createContext, useCallback, useContext, useRef, useState, type ReactNode } from "react";
import { AnimatePresence, motion } from "motion/react";
import { CheckCircle2, AlertCircle } from "lucide-react";
import { BANNER } from "@/lib/springs";

/**
 * OA 10-toast: one-off outcomes retire alone. Success pulses once; errors
 * shake. `role="status"` for successes, `role="alert"` for errors.
 */
type Toast = { id: number; kind: "success" | "error"; message: string };

const ToastCtx = createContext<{ toast: (kind: Toast["kind"], message: string) => void }>({
  toast: () => {},
});

export function useToast() {
  return useContext(ToastCtx);
}

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);
  const nextId = useRef(1);

  const toast = useCallback((kind: Toast["kind"], message: string) => {
    const id = nextId.current++;
    setToasts((t) => [...t, { id, kind, message }]);
    setTimeout(() => setToasts((t) => t.filter((x) => x.id !== id)), 4200);
  }, []);

  return (
    <ToastCtx.Provider value={{ toast }}>
      {children}
      <div
        aria-live="polite"
        className="pointer-events-none fixed bottom-4 end-4 z-50 flex w-[min(92vw,360px)] flex-col gap-2"
      >
        <AnimatePresence>
          {toasts.map((t) => (
            <motion.div
              key={t.id}
              role={t.kind === "error" ? "alert" : "status"}
              initial={{ opacity: 0, y: 12, scale: 0.97 }}
              animate={{
                opacity: 1,
                y: 0,
                scale: t.kind === "success" ? [0.97, 1.03, 1] : 1,
                x: t.kind === "error" ? [0, -6, 6, -3, 0] : 0,
              }}
              transition={{ ...BANNER, scale: { duration: 0.28 }, x: { duration: 0.32 } }}
              exit={{ opacity: 0, y: 6 }}
              className="pointer-events-auto flex items-center gap-2.5 rounded-2xl border border-border bg-popover px-4 py-3 text-sm shadow-lg"
            >
              {t.kind === "success" ? (
                <CheckCircle2 size={16} strokeWidth={1.5} className="shrink-0 text-success" aria-hidden="true" />
              ) : (
                <AlertCircle size={16} strokeWidth={1.5} className="shrink-0 text-destructive" aria-hidden="true" />
              )}
              {t.message}
            </motion.div>
          ))}
        </AnimatePresence>
      </div>
    </ToastCtx.Provider>
  );
}

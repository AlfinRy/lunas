import { useCallback, useEffect, useState } from "react";

type Theme = "light" | "dark";

function initial(): Theme {
  if (typeof window === "undefined") return "light";
  const saved = window.localStorage.getItem("lunas-theme");
  if (saved === "light" || saved === "dark") return saved;
  return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
}

/** Light/dark with localStorage persistence; tokens live in tokens.css. */
export function useTheme() {
  const [theme, setTheme] = useState<Theme>(initial);

  useEffect(() => {
    document.documentElement.classList.toggle("dark", theme === "dark");
    window.localStorage.setItem("lunas-theme", theme);
  }, [theme]);

  const toggle = useCallback(() => setTheme((t) => (t === "dark" ? "light" : "dark")), []);
  return { theme, toggle };
}

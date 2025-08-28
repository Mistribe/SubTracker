// Runtime + Build-time env helper
// Priority: runtime window.__ENV__ (Docker), fallback: import.meta.env (Vite build)

declare global {
  interface Window {
    __ENV__?: Record<string, string | undefined>;
  }
}

function readRuntime(name: string): string | undefined {
  if (typeof window !== 'undefined' && window.__ENV__) {
    return window.__ENV__[name];
  }
  return undefined;
}

function readBuild(name: string): string | undefined {
  try {
    const meta: ImportMeta | undefined = import.meta as unknown as ImportMeta | undefined;
    const env: Record<string, unknown> | undefined = (meta as unknown as { env?: Record<string, unknown> } | undefined)?.env;
    const val = env?.[name];
    return typeof val === 'string' ? val : (val !== undefined ? String(val) : undefined);
  } catch {
    return undefined;
  }
}

export function envVar(name: string): string | undefined {
  return readRuntime(name) ?? readBuild(name);
}

export function requireEnv(name: string): string {
  const v = envVar(name);
  if (v === undefined || v === '') {
    // Keep console warning to avoid throwing in runtime; components can decide defaults if needed
    console.warn(`[env] Missing required env var: ${name}`);
  }
  return (v ?? '') as string;
}

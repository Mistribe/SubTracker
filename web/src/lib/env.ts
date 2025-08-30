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
        const val = (import.meta as any).env?.[name];
        return typeof val === 'string' ? val : (val !== undefined ? String(val) : undefined);
    } catch (error) {
        console.log(`[env] Error reading build variable '${name}':`, error);
        return undefined;
    }
}

export function envVar(name: string): string | undefined {
    return readBuild(name) ?? readRuntime(name);
}

export function requireEnv(name: string): string {
    const v = envVar(name);
    if (v === undefined || v === '') {
        // Keep console warning to avoid throwing in runtime; components can decide defaults if needed
        console.debug('[env] Runtime variables:', window.__ENV__);
        console.debug('[env] Build variables:', import.meta.env);
        console.debug(`[env] Requested variable '${name}' value:`, v);
        console.warn(`[env] Missing required env var: ${name}`);
    }
    return (v ?? '') as string;
}

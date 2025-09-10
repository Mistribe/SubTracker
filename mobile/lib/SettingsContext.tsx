import React, { createContext, useCallback, useContext, useEffect, useMemo, useRef, useState } from "react";
import { Appearance, ColorSchemeName } from "react-native";
import { NativeWindStyleSheet } from "nativewind";
// We use AsyncStorage for persistence of settings
import AsyncStorage from "@react-native-async-storage/async-storage";

export type ThemeMode = "light" | "dark" | "system";

export type SettingsState = {
  themeMode: ThemeMode;
  resolvedScheme: Exclude<ColorSchemeName, null>;
  setThemeMode: (mode: ThemeMode) => void;
  currency: string;
  setCurrency: (code: string) => void;
  accountConnected: boolean;
  setAccountConnected: (connected: boolean) => void;
  initializing: boolean;
};

const SettingsContext = createContext<SettingsState | undefined>(undefined);

const STORAGE_KEYS = {
  themeMode: "settings.themeMode",
  currency: "settings.currency",
  accountConnected: "settings.accountConnected",
};

const DEFAULTS = {
  themeMode: "system" as ThemeMode,
  currency: "USD",
  accountConnected: false,
};

function resolveScheme(mode: ThemeMode, system: Exclude<ColorSchemeName, null>): Exclude<ColorSchemeName, null> {
  return mode === "system" ? system : mode;
}

export function SettingsProvider({ children }: { children: React.ReactNode }) {
  const systemScheme = Appearance.getColorScheme() ?? "light";
  const [themeMode, setThemeModeState] = useState<ThemeMode>(DEFAULTS.themeMode);
  const [resolvedScheme, setResolvedScheme] = useState<Exclude<ColorSchemeName, null>>(resolveScheme(DEFAULTS.themeMode, systemScheme));
  const [currency, setCurrencyState] = useState<string>(DEFAULTS.currency);
  const [accountConnected, setAccountConnectedState] = useState<boolean>(DEFAULTS.accountConnected);
  const [initializing, setInitializing] = useState<boolean>(true);
  const systemListenerCleanup = useRef<null | (() => void)>(null);

  // Apply NativeWind color scheme override
  const applySchemeToNativeWind = useCallback((scheme: Exclude<ColorSchemeName, null>) => {
    try {
      // NativeWind v4 API
      NativeWindStyleSheet.setColorScheme(scheme);
    } catch {}
  }, []);

  // Update resolved scheme when themeMode or system scheme changes
  const handleSchemeUpdate = useCallback(
    (mode: ThemeMode, currentSystem: Exclude<ColorSchemeName, null>) => {
      const finalScheme = resolveScheme(mode, currentSystem);
      setResolvedScheme(finalScheme);
      applySchemeToNativeWind(finalScheme);
    },
    [applySchemeToNativeWind]
  );

  // Load persisted settings on mount
  useEffect(() => {
    (async () => {
      try {
        const [savedMode, savedCurrency, savedAccount] = await Promise.all([
          AsyncStorage.getItem(STORAGE_KEYS.themeMode),
          AsyncStorage.getItem(STORAGE_KEYS.currency),
          AsyncStorage.getItem(STORAGE_KEYS.accountConnected),
        ]);
        const mode = (savedMode as ThemeMode) || DEFAULTS.themeMode;
        const curr = savedCurrency || DEFAULTS.currency;
        const acc = savedAccount === "true";
        setThemeModeState(mode);
        setCurrencyState(curr);
        setAccountConnectedState(acc);
        handleSchemeUpdate(mode, (Appearance.getColorScheme() ?? "light"));
      } catch {
        // ignore errors, continue with defaults
        handleSchemeUpdate(DEFAULTS.themeMode, (Appearance.getColorScheme() ?? "light"));
      } finally {
        setInitializing(false);
      }
    })();
  }, [handleSchemeUpdate]);

  // Listen to system theme changes only when mode is system
  useEffect(() => {
    if (systemListenerCleanup.current) {
      systemListenerCleanup.current();
      systemListenerCleanup.current = null;
    }
    if (themeMode === "system") {
      const sub = Appearance.addChangeListener(({ colorScheme }) => {
        const scheme = (colorScheme ?? "light") as Exclude<ColorSchemeName, null>;
        handleSchemeUpdate("system", scheme);
      });
      systemListenerCleanup.current = () => sub.remove();
    }
    return () => {
      if (systemListenerCleanup.current) {
        systemListenerCleanup.current();
        systemListenerCleanup.current = null;
      }
    };
  }, [themeMode, handleSchemeUpdate]);

  const setThemeMode = useCallback(async (mode: ThemeMode) => {
    setThemeModeState(mode);
    await AsyncStorage.setItem(STORAGE_KEYS.themeMode, mode);
    const currentSystem = (Appearance.getColorScheme() ?? "light") as Exclude<ColorSchemeName, null>;
    handleSchemeUpdate(mode, currentSystem);
  }, [handleSchemeUpdate]);

  const setCurrency = useCallback(async (code: string) => {
    setCurrencyState(code);
    await AsyncStorage.setItem(STORAGE_KEYS.currency, code);
  }, []);

  const setAccountConnected = useCallback(async (connected: boolean) => {
    setAccountConnectedState(connected);
    await AsyncStorage.setItem(STORAGE_KEYS.accountConnected, connected ? "true" : "false");
  }, []);

  const value = useMemo<SettingsState>(() => ({
    themeMode,
    resolvedScheme,
    setThemeMode,
    currency,
    setCurrency,
    accountConnected,
    setAccountConnected,
    initializing,
  }), [themeMode, resolvedScheme, currency, accountConnected, setThemeMode, setCurrency, setAccountConnected, initializing]);

  return (
    <SettingsContext.Provider value={value}>{children}</SettingsContext.Provider>
  );
}

export function useSettings() {
  const ctx = useContext(SettingsContext);
  if (!ctx) throw new Error("useSettings must be used within SettingsProvider");
  return ctx;
}

import "../global.css";
import { Stack } from "expo-router";
import { StatusBar } from "expo-status-bar";
import { useColorScheme } from "react-native";
import { SettingsProvider, useSettings } from "@/lib/SettingsContext";

function InnerLayout() {
  const { resolvedScheme } = useSettings();
  // StatusBar: use light content on dark scheme backgrounds
  return (
    <>
      <Stack screenOptions={{ headerShown: false }} />
      <StatusBar style={resolvedScheme === "dark" ? "light" : "dark"} />
    </>
  );
}

export default function RootLayout() {
  // keep useColorScheme import to avoid tree-shaking removing RN Appearance integration in some bundlers
  void useColorScheme();
  return (
    <SettingsProvider>
      <InnerLayout />
    </SettingsProvider>
  );
}

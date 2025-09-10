import "../global.css";
import { Stack } from "expo-router";
import { StatusBar } from "expo-status-bar";
import { useColorScheme } from "react-native";
import { ClerkProvider, useAuth } from "@clerk/clerk-expo";
import { SettingsProvider, useSettings } from "@/lib/SettingsContext";
import { getClerkPublishableKey, tokenCache } from "@/lib/auth/clerk";
import { useEffect } from "react";
import { SafeAreaProvider, SafeAreaView } from "react-native-safe-area-context";

function InnerLayout() {
  const { setAccountConnected } = useSettings();
  const { isSignedIn } = useAuth();

  useEffect(() => {
    setAccountConnected(Boolean(isSignedIn));
  }, [isSignedIn, setAccountConnected]);

  return <Stack screenOptions={{ headerShown: false }} />;
}

function StatusBarController() {
  const { resolvedScheme } = useSettings();
  return <StatusBar style={resolvedScheme === "dark" ? "light" : "dark"} />;
}

export default function RootLayout() {
  // keep useColorScheme import to avoid tree-shaking removing RN Appearance integration in some bundlers
  void useColorScheme();
  const publishableKey = getClerkPublishableKey();
  if (!publishableKey && __DEV__) {
    console.warn("EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY is not set. Clerk features will be disabled until configured.");
  }
  return (
    <ClerkProvider publishableKey={publishableKey} tokenCache={tokenCache}>
      <SettingsProvider>
        <SafeAreaProvider>
          <SafeAreaView className="flex-1 bg-white dark:bg-[#0b1220]" edges={["top"]}>
            <InnerLayout />
            <StatusBarController />
          </SafeAreaView>
        </SafeAreaProvider>
      </SettingsProvider>
    </ClerkProvider>
  );
}

import React, { useCallback, useEffect, useState } from "react";
import { View, Text, TouchableOpacity, ScrollView, ActivityIndicator } from "react-native";
import { useOAuth, isClerkAPIResponseError } from "@clerk/clerk-expo";
import * as WebBrowser from "expo-web-browser";
import { makeRedirectUri, useRouter } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { useSettings } from "@/lib/SettingsContext";

WebBrowser.maybeCompleteAuthSession();

const ONBOARDING_KEY = "onboarding.done";

export default function SignInScreen() {
  const router = useRouter();
  const { resolvedScheme, setAccountConnected } = useSettings();
  const isDark = resolvedScheme === "dark";

  const [loading, setLoading] = useState(false);
  const { startOAuthFlow } = useOAuth({ strategy: "oauth_google" });

  const onSignIn = useCallback(async () => {
    setLoading(true);
    try {
      const redirectUrl = makeRedirectUri({ scheme: "subtracker" });
      const { createdSessionId, setActive } = await startOAuthFlow({ redirectUrl });
      if (createdSessionId) {
        await setActive?.({ session: createdSessionId });
        await AsyncStorage.setItem(ONBOARDING_KEY, "true");
        await setAccountConnected(true);
        router.replace("/(tabs)");
      }
    } catch (err: any) {
      if (isClerkAPIResponseError(err)) {
        console.warn("Clerk error:", err.errors?.[0]?.message ?? err.message);
      } else {
        console.warn("Sign-in error:", err?.message ?? String(err));
      }
    } finally {
      setLoading(false);
    }
  }, [router, setAccountConnected, startOAuthFlow]);

  return (
    <ScrollView className="flex-1" contentContainerClassName={`flex-1 justify-center p-6 ${isDark ? "bg-[#0b1220]" : "bg-white"}`}>
      <View className="items-center mb-10">
        <Text className={`text-3xl font-bold ${isDark ? "text-white" : "text-gray-900"}`}>Sign in</Text>
        <Text className={`mt-2 text-base ${isDark ? "text-gray-300" : "text-gray-600"}`}>Connect your account to enable sync and sharing</Text>
      </View>
      <View className="gap-4">
        <TouchableOpacity disabled={loading} onPress={onSignIn} className={`rounded-xl py-4 px-5 ${loading ? "bg-blue-400" : "bg-blue-600"}`}>
          {loading ? (
            <View className="flex-row items-center justify-center gap-2">
              <ActivityIndicator color="#fff" />
              <Text className="text-white text-center text-base font-semibold">Signing in…</Text>
            </View>
          ) : (
            <Text className="text-white text-center text-base font-semibold">Continue with Google</Text>
          )}
        </TouchableOpacity>
        <TouchableOpacity onPress={() => router.back()} className={`rounded-xl py-4 px-5 border ${isDark ? "border-gray-700" : "border-gray-300"}`}>
          <Text className={`text-center text-base font-medium ${isDark ? "text-gray-200" : "text-gray-800"}`}>Back</Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
}

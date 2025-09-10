import { router } from "expo-router";
import { View, Text, TouchableOpacity, ScrollView } from "react-native";
import { useSettings } from "@/lib/SettingsContext";
import AsyncStorage from "@react-native-async-storage/async-storage";

const ONBOARDING_KEY = "onboarding.done";

export default function LandingScreen() {
  const { resolvedScheme, setAccountConnected } = useSettings();
  const isDark = resolvedScheme === "dark";

  async function useAnonymously() {
    await AsyncStorage.setItem(ONBOARDING_KEY, "true");
    await setAccountConnected(false);
    router.replace("/(tabs)");
  }

  function goToSignIn() {
    router.push("/(auth)/sign-in");
  }

  function goToSignUp() {
    router.push("/(auth)/sign-in"); // For now, same screen handles both
  }

  return (
    <ScrollView className="flex-1" contentContainerClassName={`flex-1 justify-center p-6 ${isDark ? "bg-[#0b1220]" : "bg-white"}`}>
      <View className="items-center mb-10">
        <Text className={`text-3xl font-bold ${isDark ? "text-white" : "text-gray-900"}`}>Welcome to SubTracker</Text>
        <Text className={`mt-2 text-base ${isDark ? "text-gray-300" : "text-gray-600"}`}>Choose how you want to get started</Text>
      </View>

      <View className="gap-4">
        <TouchableOpacity onPress={goToSignUp} className="rounded-xl py-4 px-5 bg-blue-600">
          <Text className="text-white text-center text-base font-semibold">Create an account</Text>
        </TouchableOpacity>
        <TouchableOpacity onPress={goToSignIn} className={`rounded-xl py-4 px-5 ${isDark ? "bg-gray-800" : "bg-gray-100"}`}>
          <Text className={`text-center text-base font-semibold ${isDark ? "text-white" : "text-gray-900"}`}>Log in</Text>
        </TouchableOpacity>
        <View className="mt-2" />
        <TouchableOpacity onPress={useAnonymously} className={`rounded-xl py-4 px-5 border ${isDark ? "border-gray-700" : "border-gray-300"}`}>
          <Text className={`text-center text-base font-medium ${isDark ? "text-gray-200" : "text-gray-800"}`}>Use anonymously (offline mode)</Text>
        </TouchableOpacity>
        <Text className={`mt-2 text-xs ${isDark ? "text-gray-400" : "text-gray-500"}`}>Anonymous mode stores data locally only. Family features and cloud sync are disabled until you connect an account.</Text>
      </View>
    </ScrollView>
  );
}

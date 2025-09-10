import React from "react";
import { Alert, ScrollView, Text, View, Pressable } from "react-native";
import { useSettings } from "@/lib/SettingsContext";

function SegmentButton({ label, selected, onPress }: { label: string; selected?: boolean; onPress: () => void }) {
  return (
    <Pressable
      onPress={onPress}
      className={`px-3 py-2 rounded-lg border ${selected ? "bg-blue-600 border-blue-700" : "bg-transparent border-gray-300 dark:border-gray-700"}`}
    >
      <Text className={`text-sm ${selected ? "text-white" : "text-gray-800 dark:text-gray-200"}`}>{label}</Text>
    </Pressable>
  );
}

export default function SettingsScreen() {
  const { themeMode, setThemeMode, currency, setCurrency, accountConnected, setAccountConnected } = useSettings();

  const onConnectPress = () => {
    Alert.alert(
      "Connect Account",
      "Clerk authentication flow is not configured in this build. This is a placeholder action.",
      [
        { text: "Cancel", style: "cancel" },
        {
          text: "Mark as Connected",
          onPress: () => setAccountConnected(true),
        },
      ]
    );
  };

  const onDisconnectPress = () => {
    setAccountConnected(false);
    Alert.alert("Disconnected", "Your account is now in anonymous mode.");
  };

  return (
    <ScrollView className="flex-1 bg-white dark:bg-[#0b1220]">
      <View className="p-4 space-y-6">
        <View>
          <Text className="text-2xl font-bold text-gray-900 dark:text-gray-100">Settings</Text>
          <Text className="text-gray-600 dark:text-gray-400 mt-1">Manage account, currency, and theme preferences.</Text>
        </View>

        {/* Account Section */}
        <View className="rounded-xl p-4 bg-gray-50 dark:bg-gray-900">
          <Text className="text-base font-semibold text-gray-900 dark:text-gray-100">Account</Text>
          <Text className="text-sm text-gray-600 dark:text-gray-400 mt-1">
            Status: {accountConnected ? "Connected" : "Anonymous / Not Connected"}
          </Text>
          <View className="mt-3 flex-row gap-3">
            {!accountConnected ? (
              <Pressable onPress={onConnectPress} className="px-4 py-2 rounded-lg bg-blue-600">
                <Text className="text-white font-medium">Connect account</Text>
              </Pressable>
            ) : (
              <Pressable onPress={onDisconnectPress} className="px-4 py-2 rounded-lg bg-gray-200 dark:bg-gray-800">
                <Text className="text-gray-900 dark:text-gray-100 font-medium">Disconnect</Text>
              </Pressable>
            )}
          </View>
        </View>

        {/* Preferred Currency Section */}
        <View className="rounded-xl p-4 bg-gray-50 dark:bg-gray-900">
          <Text className="text-base font-semibold text-gray-900 dark:text-gray-100">Preferred currency</Text>
          <Text className="text-sm text-gray-600 dark:text-gray-400 mt-1">Choose your display currency.</Text>
          <View className="mt-3 flex-row gap-2">
            {(["USD", "EUR", "GBP"] as const).map((code) => (
              <SegmentButton
                key={code}
                label={code}
                selected={currency === code}
                onPress={() => setCurrency(code)}
              />
            ))}
          </View>
        </View>

        {/* Theme Section */}
        <View className="rounded-xl p-4 bg-gray-50 dark:bg-gray-900">
          <Text className="text-base font-semibold text-gray-900 dark:text-gray-100">Theme</Text>
          <Text className="text-sm text-gray-600 dark:text-gray-400 mt-1">Select light, dark, or follow system.</Text>
          <View className="mt-3 flex-row gap-2">
            <SegmentButton label="Light" selected={themeMode === "light"} onPress={() => setThemeMode("light")} />
            <SegmentButton label="Dark" selected={themeMode === "dark"} onPress={() => setThemeMode("dark")} />
            <SegmentButton label="System" selected={themeMode === "system"} onPress={() => setThemeMode("system")} />
          </View>
        </View>
      </View>
    </ScrollView>
  );
}

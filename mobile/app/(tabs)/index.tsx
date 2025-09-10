import { View, Text, ScrollView } from "react-native";

export default function DashboardScreen() {
  return (
    <ScrollView className="flex-1 bg-white dark:bg-[#0b1220]">
      <View className="p-4">
        <Text className="text-2xl font-bold text-gray-900 dark:text-gray-100">Dashboard</Text>
        <Text className="text-gray-600 dark:text-gray-400 mt-1">
          Overview of your subscriptions and spending.
        </Text>

        <View className="mt-4 flex-row gap-4">
          <View className="flex-1 rounded-xl p-4 bg-blue-50 dark:bg-blue-950/30">
            <Text className="text-sm text-blue-700 dark:text-blue-300">Active</Text>
            <Text className="text-2xl font-semibold text-blue-900 dark:text-blue-200">12</Text>
          </View>
          <View className="flex-1 rounded-xl p-4 bg-emerald-50 dark:bg-emerald-950/30">
            <Text className="text-sm text-emerald-700 dark:text-emerald-300">Monthly</Text>
            <Text className="text-2xl font-semibold text-emerald-900 dark:text-emerald-200">$89.99</Text>
          </View>
        </View>

        <View className="mt-6 rounded-xl p-4 bg-gray-50 dark:bg-gray-900">
          <Text className="text-base font-semibold text-gray-900 dark:text-gray-100">Upcoming renewals</Text>
          <View className="mt-3 space-y-3">
            {["Spotify", "Netflix", "iCloud"].map((item) => (
              <View key={item} className="flex-row items-center justify-between rounded-lg p-3 bg-white dark:bg-gray-800">
                <Text className="text-gray-900 dark:text-gray-100">{item}</Text>
                <Text className="text-gray-600 dark:text-gray-400">in 5 days</Text>
              </View>
            ))}
          </View>
        </View>
      </View>
    </ScrollView>
  );
}

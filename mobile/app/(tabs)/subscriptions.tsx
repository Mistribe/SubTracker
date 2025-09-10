import { ScrollView, Text, View } from "react-native";

export default function SubscriptionsScreen() {
  return (
    <ScrollView className="flex-1 bg-white dark:bg-[#0b1220]">
      <View className="p-4">
        <Text className="text-2xl font-bold text-gray-900 dark:text-gray-100">Subscriptions</Text>
        <Text className="text-gray-600 dark:text-gray-400 mt-1">
          Manage all your subscriptions here.
        </Text>

        <View className="mt-4">
          {["Spotify", "Netflix", "Notion", "YouTube Premium"].map((name, idx) => (
            <View
              key={name}
              className="mb-3 rounded-xl p-4 bg-gray-50 dark:bg-gray-900"
            >
              <Text className="text-base font-medium text-gray-900 dark:text-gray-100">{name}</Text>
              <Text className="text-sm text-gray-600 dark:text-gray-400 mt-1">$9.{idx}9 / month</Text>
            </View>
          ))}
        </View>
      </View>
    </ScrollView>
  );
}

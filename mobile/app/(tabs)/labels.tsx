import { ScrollView, Text, View } from "react-native";

export default function LabelsScreen() {
  return (
    <ScrollView className="flex-1 bg-white dark:bg-[#0b1220]">
      <View className="p-4">
        <Text className="text-2xl font-bold text-gray-900 dark:text-gray-100">Labels</Text>
        <Text className="text-gray-600 dark:text-gray-400 mt-1">
          Organize subscriptions by labels.
        </Text>

        <View className="mt-4">
          {[
            { name: "Productivity", count: 4 },
            { name: "Entertainment", count: 5 },
            { name: "Cloud", count: 3 },
            { name: "Family", count: 2 },
          ].map((l) => (
            <View key={l.name} className="mb-3 rounded-xl p-4 bg-gray-50 dark:bg-gray-900">
              <View className="flex-row items-center justify-between">
                <Text className="text-base font-medium text-gray-900 dark:text-gray-100">{l.name}</Text>
                <Text className="text-sm text-gray-600 dark:text-gray-400">{l.count} items</Text>
              </View>
            </View>
          ))}
        </View>
      </View>
    </ScrollView>
  );
}

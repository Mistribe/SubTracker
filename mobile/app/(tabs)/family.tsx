import { ScrollView, Text, View } from "react-native";

export default function FamilyScreen() {
  return (
    <ScrollView className="flex-1 bg-white dark:bg-[#0b1220]">
      <View className="p-4">
        <Text className="text-2xl font-bold text-gray-900 dark:text-gray-100">Family</Text>
        <Text className="text-gray-600 dark:text-gray-400 mt-1">
          Manage shared subscriptions with your family.
        </Text>

        <View className="mt-4">
          {[{ name: "Max", shared: 3 }, { name: "Alex", shared: 2 }, { name: "Sam", shared: 1 }].map((m) => (
            <View key={m.name} className="mb-3 rounded-xl p-4 bg-gray-50 dark:bg-gray-900">
              <View className="flex-row items-center justify-between">
                <Text className="text-base font-medium text-gray-900 dark:text-gray-100">{m.name}</Text>
                <Text className="text-sm text-gray-600 dark:text-gray-400">{m.shared} shared</Text>
              </View>
            </View>
          ))}
        </View>
      </View>
    </ScrollView>
  );
}

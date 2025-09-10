import { ScrollView, Text, View } from "react-native";

export default function ProvidersScreen() {
  return (
    <ScrollView className="flex-1 bg-white dark:bg-[#0b1220]">
      <View className="p-4">
        <Text className="text-2xl font-bold text-gray-900 dark:text-gray-100">Providers</Text>
        <Text className="text-gray-600 dark:text-gray-400 mt-1">
          Browse services providing your subscriptions.
        </Text>

        <View className="mt-4">
          {[
            { name: "Apple", subs: 3 },
            { name: "Google", subs: 2 },
            { name: "Microsoft", subs: 1 },
            { name: "Spotify", subs: 1 },
          ].map((p) => (
            <View key={p.name} className="mb-3 rounded-xl p-4 bg-gray-50 dark:bg-gray-900">
              <View className="flex-row items-center justify-between">
                <Text className="text-base font-medium text-gray-900 dark:text-gray-100">{p.name}</Text>
                <Text className="text-sm text-gray-600 dark:text-gray-400">{p.subs} subscriptions</Text>
              </View>
            </View>
          ))}
        </View>
      </View>
    </ScrollView>
  );
}

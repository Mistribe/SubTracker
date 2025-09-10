import "../global.css";
import { Stack } from "expo-router";
import { StatusBar } from "expo-status-bar";
import { useColorScheme } from "react-native";

export default function RootLayout() {
  const scheme = useColorScheme();
  return (
    <>
      <Stack screenOptions={{ headerShown: false }} />
      <StatusBar style={scheme === "dark" ? "light" : "dark"} />
    </>
  );
}

import React, { useEffect, useState } from "react";
import { Redirect } from "expo-router";
import AsyncStorage from "@react-native-async-storage/async-storage";
import { useAuth } from "@clerk/clerk-expo";

const ONBOARDING_KEY = "onboarding.done";

export default function Index() {
  const [onboarded, setOnboarded] = useState<boolean | null>(null);
  const { isSignedIn } = useAuth();

  useEffect(() => {
    (async () => {
      try {
        const v = await AsyncStorage.getItem(ONBOARDING_KEY);
        const done = v === "true";
        if (isSignedIn && !done) {
          await AsyncStorage.setItem(ONBOARDING_KEY, "true");
          setOnboarded(true);
          return;
        }
        setOnboarded(done);
      } catch {
        setOnboarded(Boolean(isSignedIn));
      }
    })();
  }, [isSignedIn]);

  if (onboarded === null) return null;

  if (!onboarded) {
    return <Redirect href="/(onboarding)/landing" />;
  }

  return <Redirect href="/(tabs)" />;
}

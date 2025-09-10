import * as SecureStore from 'expo-secure-store';
import Constants from 'expo-constants';

export const tokenCache = {
  getToken: async (key: string) => {
    try {
      return await SecureStore.getItemAsync(key);
    } catch {
      return null;
    }
  },
  saveToken: async (key: string, value: string) => {
    try {
      await SecureStore.setItemAsync(key, value);
    } catch {}
  },
};

export function getClerkPublishableKey(): string {
  // Priority: env -> expo extra -> dev fallback
  const fromEnv = process.env.EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY;
  const extra = (Constants.expoConfig?.extra ?? {}) as Record<string, any>;
  const fromExtra = extra.clerkPublishableKey as string | undefined;
  if (fromEnv) return fromEnv;
  if (fromExtra) return fromExtra;
  if (__DEV__) return 'pk_test_default_development_key';
  return '';
}

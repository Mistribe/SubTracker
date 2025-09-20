import { useRef } from 'react';
import { useAuth } from '@clerk/clerk-react';
import { requireEnv } from '@/lib/env';
import { Configuration } from '@/api';
import { CurrenciesApi, FamilyApi, LabelsApi, ProvidersApi, SubscriptionsApi, UsersApi } from '@/api/apis';
import type { ApiClient } from '@/lib/api-client';

export function useApiClient() {
  const { getToken } = useAuth();
  const clientRef = useRef<ApiClient | null>(null);

  if (clientRef.current === null) {
    const basePath = requireEnv('VITE_BACKEND_URL');

    const authMiddleware = {
      pre: async ({ url, init }: { url: string; init: RequestInit }) => {
        const token = await getToken();
        const headers: Record<string, string> = {
          ...(init.headers as Record<string, string> | undefined),
        };
        if (token) {
          headers['Authorization'] = `Bearer ${token}`;
        }
        return { url, init: { ...init, headers } };
      },
    };

    const configuration = new Configuration({
      basePath,
      middleware: [authMiddleware],
    });

    clientRef.current = {
      currencies: new CurrenciesApi(configuration),
      families: new FamilyApi(configuration),
      labels: new LabelsApi(configuration),
      providers: new ProvidersApi(configuration),
      subscriptions: new SubscriptionsApi(configuration),
      users: new UsersApi(configuration),
    };
  }

  return { apiClient: clientRef.current };
}
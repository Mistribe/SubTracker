import { type ApiClient, createApiClient } from '../api/apiClient';
import { FetchRequestAdapter } from '@microsoft/kiota-http-fetchlibrary';
import { useRef } from 'react';
import { useAuth } from '@clerk/clerk-react';
import { AuthTokenProvider } from '@/lib/AuthTokenProvider';
import { requireEnv } from '@/lib/env';

export function useApiClient() {
  const { getToken } = useAuth();
  const clientRef = useRef<ApiClient | null>(null);

  if (clientRef.current === null) {
    const tokenProvider = new AuthTokenProvider(async () => await getToken());
    const adapter = new FetchRequestAdapter(tokenProvider);
    adapter.baseUrl = requireEnv('VITE_BACKEND_URL');
    clientRef.current = createApiClient(adapter);
  }

  return { apiClient: clientRef.current };
}
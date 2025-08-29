import {type ApiClient, createApiClient} from '../api/apiClient';
import {FetchRequestAdapter} from '@microsoft/kiota-http-fetchlibrary';
import {useState, useEffect} from 'react';
import { useAuth } from '@clerk/clerk-react';
import { AuthTokenProvider } from '@/lib/AuthTokenProvider';
import { requireEnv } from '@/lib/env';

export function useApiClient() {
    const [apiClient, setApiClient] = useState<ApiClient | null>(null)
    const [isLoading, setIsLoading] = useState(true);

    const { getToken } = useAuth();

    useEffect(() => {
        const initializeClient = async () => {
            try {
                const tokenProvider = new AuthTokenProvider(async () => await getToken());
                const adapter = new FetchRequestAdapter(tokenProvider);
                adapter.baseUrl = requireEnv('VITE_BACKEND_URL');
                const client = createApiClient(adapter);
                setApiClient(client);
            } catch (error) {
                console.error('Failed to initialize API client:', error);
            } finally {
                setIsLoading(false);
            }
        };

        initializeClient();
    }, [getToken]);


    return {apiClient, isLoading};
}
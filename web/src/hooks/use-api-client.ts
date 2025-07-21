import {type ApiClient, createApiClient} from '../api/apiClient';
import {FetchRequestAdapter} from '@microsoft/kiota-http-fetchlibrary';
import {useState, useEffect} from 'react';
import {useKindeAuth} from "@kinde-oss/kinde-auth-react";
import {KindeAuthProvider} from '@/lib/KindeAuthProvider';

export function useApiClient() {
    const [apiClient, setApiClient] = useState<ApiClient | null>(null)
    const [isLoading, setIsLoading] = useState(true);

    const {getToken} = useKindeAuth();

    useEffect(() => {
        const initializeClient = async () => {
            try {
                const authProvider = new KindeAuthProvider(getToken);
                const adapter = new FetchRequestAdapter(authProvider);
                adapter.baseUrl = import.meta.env.VITE_BACKEND_URL;
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
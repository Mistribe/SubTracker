import {type ApiClient, createApiClient} from '../api/apiClient';
import {FetchRequestAdapter} from '@microsoft/kiota-http-fetchlibrary';
import {useState, useEffect} from 'react';
import {useKindeAuth} from "@kinde-oss/kinde-auth-react";
import {KindeAuthProvider} from '@/lib/KindeAuthProvider';

export function useApiClient() {
    const [apiClient, setApiClient] = useState<ApiClient | null>(null);
    const {getToken} = useKindeAuth();

    useEffect(() => {
        const authProvider = new KindeAuthProvider(getToken);
        const adapter = new FetchRequestAdapter(authProvider);
        const client = createApiClient(adapter);
        setApiClient(client.withUrl(import.meta.env.VITE_BACKEND_URL));
    }, [getToken]);


    return apiClient;
}
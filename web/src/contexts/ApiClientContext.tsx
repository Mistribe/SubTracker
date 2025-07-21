// src/contexts/ApiClientContext.tsx
import {createContext, type ReactNode, useContext} from 'react';
import {useApiClient} from "@/hooks/use-api-client.ts";
import type {ApiClient} from "@/api/apiClient.ts";


const ApiClientContext = createContext<ApiClient | null>(null);

export function ApiClientProvider({children}: { children: ReactNode }) {
    const {apiClient} = useApiClient();

    return (
        <ApiClientContext.Provider value={apiClient}>
            {children}
        </ApiClientContext.Provider>
    );
}

export function useApiClientContext() {
    const context = useContext(ApiClientContext);
    if (!context) {
        throw new Error('useApiClientContext must be used within an ApiClientProvider');
    }
    return context;
}
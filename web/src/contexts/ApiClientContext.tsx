// src/contexts/ApiClientContext.tsx
import {createContext, type ReactNode} from 'react';
import {useApiClient} from "@/hooks/use-api-client.ts";
import type {ApiClient} from "@/lib/api-client";


const ApiClientContext = createContext<ApiClient | null>(null);

export function ApiClientProvider({children}: { children: ReactNode }) {
    const {apiClient} = useApiClient();

    return (
        <ApiClientContext.Provider value={apiClient}>
            {children}
        </ApiClientContext.Provider>
    );
}
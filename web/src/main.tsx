import {StrictMode} from 'react'
import {createRoot} from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import {KindeProvider} from "@kinde-oss/kinde-auth-react";
import { requireEnv } from './lib/env';

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <KindeProvider
            clientId={requireEnv('VITE_KINDE_CLIENT_ID')}
            domain={requireEnv('VITE_KINDE_DOMAIN')}
            logoutUri={requireEnv('VITE_KINDE_LOGOUT_URL')}
            redirectUri={requireEnv('VITE_KINDE_REDIRECT_URL')}
            audience={requireEnv('VITE_KINDE_AUDIENCE')}
            // When running local against a custom domain, include the line below
            // useInsecureForRefreshToken={true}
        >
            <App/>
        </KindeProvider>
    </StrictMode>,
)

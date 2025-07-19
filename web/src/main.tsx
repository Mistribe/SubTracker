import {StrictMode} from 'react'
import {createRoot} from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import {KindeProvider} from "@kinde-oss/kinde-auth-react";

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <KindeProvider
            clientId={import.meta.env.VITE_KINDE_CLIENT_ID}
            domain={import.meta.env.VITE_KINDE_DOMAIN}
            logoutUri={import.meta.env.VITE_KINDE_LOGOUT_URL}
            redirectUri={import.meta.env.VITE_KINDE_REDIRECT_URL}
            // When running local against a custom domain, include the line below
            // useInsecureForRefreshToken={true}
        >
            <App/>
        </KindeProvider>
    </StrictMode>,
)

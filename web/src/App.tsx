import './App.css'
import { KindeProvider } from "@kinde-oss/kinde-auth-react";
import { ThemeProvider } from './contexts/ThemeContext';
import Routes from './Routes';

function App() {
  // Use window.location.origin to dynamically set the redirect and logout URIs
  const baseUrl = window.location.origin;
  
  return (
    <ThemeProvider>
      <KindeProvider
          clientId="8021135b4546409488061c2ed6ac3a51"
          domain="https://mistribe.kinde.com"
          redirectUri={`${baseUrl}/`}
          logoutUri={`${baseUrl}/`}
      >
          <Routes />
      </KindeProvider>
    </ThemeProvider>
  )
}

export default App

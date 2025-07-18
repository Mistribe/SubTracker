import './App.css'
import { KindeProvider } from "@kinde-oss/kinde-auth-react";
import Routes from './Routes';

function App() {
  return (
      <KindeProvider
          clientId="8021135b4546409488061c2ed6ac3a51"
          domain="https://mistribe.kinde.com"
          redirectUri="com.mistribe.app://"
          logoutUri="com.mistribe.app://"
      >
          <Routes />
      </KindeProvider>
  )
}

export default App

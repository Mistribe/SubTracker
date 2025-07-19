import {ThemeProvider} from "@/components/theme-provider"
import {BrowserRouter, Routes, Route} from "react-router-dom"
import {AuthProvider} from "@/contexts/AuthProvider"
import {ProtectedRoute} from "@/contexts/ProtectedRoute"
import HomePage from "@/pages/HomePage"
import LoginPage from "@/pages/LoginPage"
import DashboardPage from "@/pages/DashboardPage"
import ProfilePage from "@/pages/ProfilePage"
import {KindeProvider} from "@kinde-oss/kinde-auth-react";

function App() {
    return (
        <BrowserRouter>
            <KindeProvider
                clientId="8021135b4546409488061c2ed6ac3a51"
                domain="https://mistribe.kinde.com"
                redirectUri="com.mistribe.subtracker://"
                logoutUri="com.mistribe.subtracker://">
                <AuthProvider>
                    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
                        <Routes>
                            {/* Public Routes */}
                            <Route path="/" element={<HomePage/>}/>
                            <Route path="/login" element={<LoginPage/>}/>

                            {/* Protected Routes */}
                            <Route
                                path="/dashboard"
                                element={
                                    <ProtectedRoute>
                                        <DashboardPage/>
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/profile"
                                element={
                                    <ProtectedRoute>
                                        <ProfilePage/>
                                    </ProtectedRoute>
                                }
                            />

                            {/* Fallback Route */}
                            <Route path="*" element={<HomePage/>}/>
                        </Routes>
                    </ThemeProvider>
                </AuthProvider>
            </KindeProvider>
        </BrowserRouter>
    )
}

export default App
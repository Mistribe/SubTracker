import {ThemeProvider} from "@/components/theme-provider"
import {BrowserRouter, Routes, Route} from "react-router-dom"
import {ProtectedRoute} from "@/contexts/ProtectedRoute"
import {AppLayout} from "@/layouts/AppLayout"
import HomePage from "@/pages/HomePage"
import DashboardPage from "@/pages/DashboardPage"
import ProfilePage from "@/pages/ProfilePage"
import FamiliesPage from "@/pages/FamiliesPage"
import LabelsPage from "@/pages/LabelsPage"
import {ApiClientProvider} from "@/contexts/ApiClientContext.tsx";
import {QueryClientProvider, QueryClient} from "@tanstack/react-query";

const queryClient = new QueryClient()

function App() {
    return (
        <BrowserRouter>
            <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
                <ApiClientProvider>
                    <QueryClientProvider client={queryClient}>
                        <Routes>
                            {/* Public Routes */}
                            <Route path="/" element={<HomePage/>}/>

                            {/* Protected Routes */}
                            <Route
                                path="/dashboard"
                                element={
                                    <ProtectedRoute>
                                        <AppLayout>
                                            <DashboardPage/>
                                        </AppLayout>
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/families"
                                element={
                                    <ProtectedRoute>
                                        <AppLayout>
                                            <FamiliesPage/>
                                        </AppLayout>
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/labels"
                                element={
                                    <ProtectedRoute>
                                        <AppLayout>
                                            <LabelsPage/>
                                        </AppLayout>
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/profile"
                                element={
                                    <ProtectedRoute>
                                        <AppLayout>
                                            <ProfilePage/>
                                        </AppLayout>
                                    </ProtectedRoute>
                                }
                            />

                            {/* Fallback Route */}
                            <Route path="*" element={<HomePage/>}/>
                        </Routes>
                    </QueryClientProvider>
                </ApiClientProvider>
            </ThemeProvider>
        </BrowserRouter>
    )
}

export default App
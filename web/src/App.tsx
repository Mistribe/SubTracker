import { ThemeProvider } from "@/components/theme-provider"
import { BrowserRouter, Routes, Route } from "react-router-dom"
import { ProtectedRoute } from "@/contexts/ProtectedRoute"
import { AppLayout } from "@/layouts/AppLayout"
import HomePage from "@/pages/HomePage"
import SignInPage from "@/pages/SignInPage"
import SignUpPage from "@/pages/SignUpPage"
import DashboardPage from "@/pages/DashboardPage"
import PreferencePage from "@/pages/PreferencePage.tsx"
import FamilyPage from "@/pages/FamilyPage"
import LabelsPage from "@/pages/LabelsPage"
import ImportLabelsPage from "@/pages/ImportLabelsPage"
import ProvidersPage from "@/pages/ProvidersPage"
import ImportProvidersPage from "@/pages/ImportProvidersPage"
import SubscriptionsPage from "@/pages/SubscriptionsPage"
import ImportSubscriptionsPage from "@/pages/ImportSubscriptionsPage"
import SubscriptionFormPage from "@/pages/SubscriptionFormPage"
import ProviderDetailPage from "@/pages/ProviderDetailPage";
import AcceptInvitationPage from "@/pages/AcceptInvitationPage";
import { ApiClientProvider } from "@/contexts/ApiClientContext.tsx";
import { QueryClientProvider, QueryClient } from "@tanstack/react-query";
import { Toaster } from "@/components/ui/sonner";

const queryClient = new QueryClient()

function App() {
    return (
        <BrowserRouter>
            <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
                <ApiClientProvider>
                    <QueryClientProvider client={queryClient}>
                        <Toaster />
                        <Routes>
                            {/* Public Routes */}
                            <Route path="/" element={<HomePage />} />
                            <Route path="/sign-in" element={<SignInPage />} />
                            <Route path="/sign-up" element={<SignUpPage />} />

                            {/* Protected Routes */}
                            <Route
                                path="/dashboard"
                                element={
                                    <ProtectedRoute>
                                        <AppLayout>
                                            <DashboardPage />
                                        </AppLayout>
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/family"
                                element={
                                    <ProtectedRoute>
                                        <AppLayout>
                                            <FamilyPage />
                                        </AppLayout>
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/labels"
                                element={
                                    <ProtectedRoute>
                                        <AppLayout>
                                            <LabelsPage />
                                        </AppLayout>
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/labels/import"
                                element={
                                    <ProtectedRoute>
                                        <AppLayout>
                                            <ImportLabelsPage />
                                        </AppLayout>
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/profile"
                                element={
                                    <ProtectedRoute>
                                        <AppLayout>
                                            <PreferencePage />
                                        </AppLayout>
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/providers"
                                element={
                                    <ProtectedRoute>
                                        <AppLayout>
                                            <ProvidersPage />
                                        </AppLayout>
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/providers/import"
                                element={
                                    <ProtectedRoute>
                                        <AppLayout>
                                            <ImportProvidersPage />
                                        </AppLayout>
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/providers/:providerId"
                                element={
                                    <ProtectedRoute>
                                        <AppLayout>
                                            <ProviderDetailPage />
                                        </AppLayout>
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/subscriptions"
                                element={
                                    <ProtectedRoute>
                                        <AppLayout>
                                            <SubscriptionsPage />
                                        </AppLayout>
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/subscriptions/import"
                                element={
                                    <ProtectedRoute>
                                        <AppLayout>
                                            <ImportSubscriptionsPage />
                                        </AppLayout>
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/subscriptions/create"
                                element={
                                    <ProtectedRoute>
                                        <AppLayout>
                                            <SubscriptionFormPage />
                                        </AppLayout>
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/subscriptions/edit/:subscriptionId"
                                element={
                                    <ProtectedRoute>
                                        <AppLayout>
                                            <SubscriptionFormPage />
                                        </AppLayout>
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/invite/accept"
                                element={
                                    <ProtectedRoute>
                                        <AppLayout>
                                            <AcceptInvitationPage />
                                        </AppLayout>
                                    </ProtectedRoute>
                                }
                            />

                            {/* Fallback Route */}
                            <Route path="*" element={<HomePage />} />
                        </Routes>
                    </QueryClientProvider>
                </ApiClientProvider>
            </ThemeProvider>
        </BrowserRouter>
    )
}

export default App
import {useEffect, useState} from "react";
import {useKindeAuth} from "@kinde-oss/kinde-auth-react";
import {useTheme} from "@/components/theme-provider";
import {PageHeader} from "@/components/ui/page-header";
import {UserProfileSection} from "@/components/profile/UserProfileSection";
import {ProfileDetailsContainer} from "@/components/profile/ProfileDetailsSection";
import {PreferencesSection} from "@/components/profile/PreferencesSection";
import {AccountSettingsSection} from "@/components/profile/AccountSettingsSection";
import {useProfileManagement} from "@/hooks/profile/useProfileManagement";

/**
 * ProfilePage component
 * Displays and manages user profile information, including preferences like currency and theme
 */
const ProfilePage = () => {
    // Get authentication data from Kinde
    const {user: kindeUser} = useKindeAuth();
    // Get theme information from theme context
    const {theme} = useTheme();
    // Use profile management hook for currency preferences and profile updates
    const {
        preferredCurrency,
        isLoadingPreferredCurrency,
        isErrorPreferredCurrency,
        availableCurrencies,
        isLoadingAvailableCurrencies,
        isErrorAvailableCurrencies,
        updateProfile,
        updateProfileName,
        isUpdating
    } = useProfileManagement();

    // User data from Kinde Auth with default values for missing fields
    const [user, setUser] = useState({
        givenName: kindeUser?.givenName || "John",
        familyName: kindeUser?.familyName || "Doe",
        email: kindeUser?.email || "john.doe@example.com",
        joinDate: new Date().toISOString().split('T')[0], // Current date as join date
        preferredCurrency: "USD", // Default currency
        preferredTheme: "system"
    });

    // Update user data when Kinde Auth data changes
    useEffect(() => {
        if (kindeUser) {
            setUser(prev => ({
                ...prev,
                givenName: kindeUser.givenName || prev.givenName,
                familyName: kindeUser.familyName || prev.familyName,
                email: kindeUser.email || prev.email
            }));
        }
    }, [kindeUser]);

    // Update preferred currency when it's loaded from the backend
    useEffect(() => {
        if (preferredCurrency) {
            setUser(prev => ({
                ...prev,
                preferredCurrency: preferredCurrency.currency || prev.preferredCurrency
            }));
        }
    }, [preferredCurrency]);

    // Update theme preference when theme changes in the theme context
    useEffect(() => {
        setUser(prev => ({...prev, preferredTheme: theme}));
    }, [theme]);

    // Format available currencies from API for the dropdown
    const currencies = availableCurrencies?.map(currency => ({
        value: currency.code,
        label: `${currency.name} (${currency.symbol})`
    })) || [
        // Fallback currencies if API data is not available
        {value: "USD", label: "US Dollar ($)"},
        {value: "EUR", label: "Euro (€)"}
    ];

    return (
        <div className="container mx-auto py-6">
            <PageHeader
                title="Profile"
                description="Manage your profile and preferences"
            />

            {(isLoadingPreferredCurrency || isLoadingAvailableCurrencies) && (
                <div className="flex justify-center my-12">
                    <div className="text-center">
                        <div
                            className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
                        <p className="text-muted-foreground">Loading profile data...</p>
                    </div>
                </div>
            )}

            {(isErrorPreferredCurrency || isErrorAvailableCurrencies) && (
                <div className="bg-destructive/10 text-destructive p-4 rounded-md my-6">
                    <h3 className="font-medium">Error loading profile</h3>
                    <p>There was a problem loading your profile data. Please try refreshing the page.</p>
                </div>
            )}

            <div className="max-w-3xl mx-auto space-y-12">
                {/* User Profile Section */}
                <div className="space-y-8">
                    {/* User Avatar and Basic Info */}
                    <UserProfileSection
                        name={`${user.givenName} ${user.familyName}`}
                        email={user.email}
                        picture={kindeUser?.picture}
                    />

                    {/* Profile Details */}
                    <ProfileDetailsContainer
                        givenName={user.givenName}
                        familyName={user.familyName}
                        email={user.email}
                        joinDate={user.joinDate}
                        isLoading={isUpdating}
                        onSave={(data) => {
                            setUser(prev => ({...prev, ...data}));
                            // Call the API to update the profile
                            if (updateProfileName) {
                                updateProfileName(data.givenName, data.familyName);
                            }
                        }}
                    />
                </div>

                {/* Preferences Section */}
                <PreferencesSection
                    preferredCurrency={user.preferredCurrency}
                    isUpdating={isUpdating}
                    onCurrencyChange={(value) => {
                        // Update local state for immediate UI feedback
                        setUser(prev => ({...prev, preferredCurrency: value}));

                        // Persist to backend using the profile management hook
                        updateProfile(value);
                    }}
                    currencies={currencies}
                />

                {/* Account Settings Section */}
                <AccountSettingsSection
                    onEmailNotificationsToggle={() => {
                        // Handle email notifications toggle
                        console.log("Toggle email notifications");
                    }}
                    onChangePassword={() => {
                        // Handle password change
                        console.log("Change password");
                    }}
                    onDeleteAccount={() => {
                        // Handle account deletion
                        console.log("Delete account");
                    }}
                />
            </div>
        </div>
    );
};

export default ProfilePage;
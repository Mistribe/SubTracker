import {PageHeader} from "@/components/ui/page-header";
import {PreferencesSection} from "@/components/profile/PreferencesSection";
import {useProfileManagement} from "@/hooks/profile/useProfileManagement";

/**
 * Preferences page
 * Only includes Preferred Currency and Theme chooser
 */
const PreferencePage = () => {
    const {
        preferredCurrency,
        isLoadingPreferredCurrency,
        isErrorPreferredCurrency,
        availableCurrencies,
        isLoadingAvailableCurrencies,
        isErrorAvailableCurrencies,
        updateProfile,
        isUpdating
    } = useProfileManagement();

    // Format available currencies from API for the dropdown
    const currencies = availableCurrencies?.map(currency => ({
        value: currency.code,
        label: `${currency.name} (${currency.symbol})`
    })) || [
        // Fallback currencies if API data is not available
        {value: "USD", label: "US Dollar ($)"},
        {value: "EUR", label: "Euro (€)"}
    ];

    const selectedCurrency = preferredCurrency?.currency ?? "USD";

    return (
        <div className="container mx-auto py-6">
            <PageHeader
                title="Preferences"
                description="Manage your app preferences"
            />

            {(isLoadingPreferredCurrency || isLoadingAvailableCurrencies) && (
                <div className="flex justify-center my-12">
                    <div className="text-center">
                        <div
                            className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
                        <p className="text-muted-foreground">Loading preferences...</p>
                    </div>
                </div>
            )}

            {(isErrorPreferredCurrency || isErrorAvailableCurrencies) && (
                <div className="bg-destructive/10 text-destructive p-4 rounded-md my-6">
                    <h3 className="font-medium">Error loading preferences</h3>
                    <p>There was a problem loading your preferences. Please try refreshing the page.</p>
                </div>
            )}

            <div className="max-w-3xl mx-auto">
                <PreferencesSection
                    preferredCurrency={selectedCurrency}
                    isUpdating={isUpdating}
                    onCurrencyChange={(value) => {
                        // Persist to backend using the profile management hook
                        updateProfile(value);
                    }}
                    currencies={currencies}
                />
            </div>
        </div>
    );
};

export default PreferencePage;
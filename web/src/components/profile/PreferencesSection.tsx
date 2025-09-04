import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import { ModeToggle } from "@/components/mode-toggle";

interface Currency {
  value: string;
  label: string;
}

interface PreferencesSectionProps {
  preferredCurrency: string;
  isUpdating: boolean;
  onCurrencyChange: (currency: string) => void;
  currencies: Currency[];
}

/**
 * PreferencesSection component
 * Displays and manages user preferences like currency and theme
 */
export const PreferencesSection = ({
  preferredCurrency,
  isUpdating,
  onCurrencyChange,
  currencies,
}: PreferencesSectionProps) => {
  return (
    <div className="space-y-6">
      <div className="space-y-6">
        {/* Currency Preference */}
        <div className="flex justify-between items-center">
          <div>
            <h4 className="font-medium">Preferred Currency</h4>
            <p className="text-sm text-muted-foreground">
              Select your preferred currency for the application
            </p>
          </div>
          <div className="flex items-center">
            <Select
              value={preferredCurrency}
              onValueChange={onCurrencyChange}
              disabled={isUpdating}
            >
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Select currency" />
              </SelectTrigger>
              <SelectContent>
                {currencies.map((currency) => (
                  <SelectItem key={currency.value} value={currency.value}>
                    {currency.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {isUpdating && <span className="ml-2 text-sm text-muted-foreground">Saving...</span>}
          </div>
        </div>

        <Separator />

        {/* Theme Preference */}
        <div className="flex justify-between items-center">
          <div>
            <h4 className="font-medium">Theme</h4>
            <p className="text-sm text-muted-foreground">Choose between light and dark mode</p>
          </div>
          <ModeToggle />
        </div>
      </div>
    </div>
  );
};
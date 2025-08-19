import React from "react";
import {Sheet, SheetContent, SheetFooter, SheetHeader, SheetTitle} from "@/components/ui/sheet";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {Label} from "@/components/ui/label";
import {Switch} from "@/components/ui/switch";
import Provider from "@/models/provider";
import {SubscriptionRecurrency} from "@/models/subscriptionRecurrency";
import type {MultiSelectOption} from "@/components/ui/multi-select";
import {MultiSelect} from "@/components/ui/multi-select";

export interface SubscriptionsFiltersProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    fromDateStr: string;
    onFromDateChange: (value: string) => void;
    toDateStr: string;
    onToDateChange: (value: string) => void;
    providerMap?: Map<string, Provider>;
    providersFilter: string[];
    onToggleProvider: (id: string) => void;
    recurrenciesFilter: SubscriptionRecurrency[];
    onToggleRecurrency: (r: SubscriptionRecurrency) => void;
    usersCsv: string;
    onUsersCsvChange: (value: string) => void;
    usersOptions?: MultiSelectOption[];
    withInactive: boolean;
    onWithInactiveChange: (value: boolean) => void;
    onClear: () => void;
    onApply: () => void;
}

// todo use an object to avoid this crap
export const SubscriptionsFilters: React.FC<SubscriptionsFiltersProps> = ({
                                                                              open,
                                                                              onOpenChange,
                                                                              fromDateStr,
                                                                              onFromDateChange,
                                                                              toDateStr,
                                                                              onToDateChange,
                                                                              providerMap,
                                                                              providersFilter,
                                                                              onToggleProvider,
                                                                              recurrenciesFilter,
                                                                              onToggleRecurrency,
                                                                              usersCsv,
                                                                              onUsersCsvChange,
                                                                              usersOptions,
                                                                              withInactive,
                                                                              onWithInactiveChange,
                                                                              onClear,
                                                                              onApply,
                                                                          }) => {
    return (
        <Sheet open={open} onOpenChange={onOpenChange}>
            <SheetContent side="right" className="sm:max-w-md">
                <SheetHeader>
                    <SheetTitle>Filters</SheetTitle>
                </SheetHeader>
                <div className="p-4 space-y-6 overflow-y-auto">
                    {/* Date range */}
                    <div className="grid grid-cols-2 gap-4">
                        <div className="space-y-2">
                            <Label htmlFor="fromDate">From date</Label>
                            <Input
                                id="fromDate"
                                type="date"
                                value={fromDateStr}
                                onChange={(e) => onFromDateChange(e.target.value)}
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="toDate">To date</Label>
                            <Input
                                id="toDate"
                                type="date"
                                value={toDateStr}
                                onChange={(e) => onToDateChange(e.target.value)}
                            />
                        </div>
                    </div>

                    {/* Providers */}
                    {providerMap && providerMap.size > 0 && (
                        <div className="space-y-2">
                            <div className="font-medium">Providers</div>
                            <MultiSelect
                                options={Array.from(providerMap.values()).map((p) => ({value: p.id, label: p.name}))}
                                selectedValues={providersFilter}
                                onToggle={onToggleProvider}
                                placeholder="Select providers"
                                searchPlaceholder="Search providers..."
                                emptyMessage="No provider found."
                            />
                        </div>
                    )}

                    {/* Recurrencies */}
                    <div className="space-y-2">
                        <div className="font-medium">Recurrencies</div>
                        <MultiSelect
                            options={(Object.values(SubscriptionRecurrency) as SubscriptionRecurrency[]).map((r) => ({
                                value: r,
                                label: r
                            }))}
                            selectedValues={recurrenciesFilter as string[]}
                            onToggle={(v) => onToggleRecurrency(v as SubscriptionRecurrency)}
                            placeholder="Select recurrencies"
                            searchPlaceholder="Search recurrencies..."
                            emptyMessage="No recurrencies found."
                        />
                    </div>

                    {/* Family Members */}
                    <div className="space-y-2">
                        <div className="font-medium">Family members</div>
                        {usersOptions && usersOptions.length > 0 ? (
                            <MultiSelect
                                options={usersOptions}
                                selectedValues={usersCsv ? usersCsv.split(',').map(s => s.trim()).filter(Boolean) : []}
                                onChange={(values) => onUsersCsvChange(values.join(','))}
                                placeholder="Select family members"
                                searchPlaceholder="Search members..."
                                emptyMessage="No family members found."
                            />
                        ) : (
                            <>
                                <Label htmlFor="users">Users (CSV)</Label>
                                <Input
                                    id="users"
                                    placeholder="user1,user2"
                                    value={usersCsv}
                                    onChange={(e) => onUsersCsvChange(e.target.value)}
                                />
                            </>
                        )}
                    </div>

                    {/* Include inactive */}
                    <div className="flex items-center justify-between">
                        <Label htmlFor="withInactive">Include inactive</Label>
                        <Switch
                            id="withInactive"
                            checked={withInactive}
                            onCheckedChange={(v) => onWithInactiveChange(!!v)}
                        />
                    </div>
                </div>
                <SheetFooter>
                    <div className="flex w-full gap-2">
                        <Button variant="secondary" className="flex-1" onClick={onClear}>
                            Clear
                        </Button>
                        <Button className="flex-1" onClick={onApply}>
                            Apply
                        </Button>
                    </div>
                </SheetFooter>
            </SheetContent>
        </Sheet>
    );
};

export default SubscriptionsFilters;

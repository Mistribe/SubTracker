import React from "react";
import {Sheet, SheetContent, SheetFooter, SheetHeader, SheetTitle} from "@/components/ui/sheet";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {Checkbox} from "@/components/ui/checkbox";
import {Label} from "@/components/ui/label";
import {Switch} from "@/components/ui/switch";
import Provider from "@/models/provider";
import {SubscriptionRecurrency} from "@/models/subscriptionRecurrency";

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
    withInactive: boolean;
    onWithInactiveChange: (value: boolean) => void;
    onClear: () => void;
    onApply: () => void;
}

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
                            <div className="grid grid-cols-1 gap-2 max-h-48 overflow-y-auto pr-2">
                                {Array.from(providerMap.values()).map((p) => (
                                    <label key={p.id} className="flex items-center gap-2">
                                        <Checkbox
                                            checked={providersFilter.includes(p.id)}
                                            onCheckedChange={() => onToggleProvider(p.id)}
                                        />
                                        <span className="text-sm">{p.name}</span>
                                    </label>
                                ))}
                            </div>
                        </div>
                    )}

                    {/* Recurrencies */}
                    <div className="space-y-2">
                        <div className="font-medium">Recurrencies</div>
                        <div className="grid grid-cols-1 gap-2">
                            {(Object.values(SubscriptionRecurrency) as SubscriptionRecurrency[]).map((r) => (
                                <label key={r} className="flex items-center gap-2">
                                    <Checkbox
                                        checked={recurrenciesFilter.includes(r)}
                                        onCheckedChange={() => onToggleRecurrency(r)}
                                    />
                                    <span className="text-sm">{r}</span>
                                </label>
                            ))}
                        </div>
                    </div>

                    {/* Users CSV */}
                    <div className="space-y-2">
                        <Label htmlFor="users">Users (CSV)</Label>
                        <Input
                            id="users"
                            placeholder="user1,user2"
                            value={usersCsv}
                            onChange={(e) => onUsersCsvChange(e.target.value)}
                        />
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

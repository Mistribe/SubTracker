import React, {useEffect, useState} from "react";
import {Sheet, SheetContent, SheetFooter, SheetHeader, SheetTitle} from "@/components/ui/sheet";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {Label} from "@/components/ui/label";
import {Switch} from "@/components/ui/switch";
import Provider from "@/models/provider";
import {SubscriptionRecurrency} from "@/models/subscriptionRecurrency";
import type {MultiSelectOption} from "@/components/ui/multi-select";
import {MultiSelect} from "@/components/ui/multi-select";

export interface SubscriptionsFiltersValues {
    fromDate?: Date;
    toDate?: Date;
    providers?: string[];
    recurrenciesFilter?: SubscriptionRecurrency[];
    users?: string[];
    withInactive?: boolean;
}

export interface SubscriptionsFiltersProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    fromDate?: Date;
    toDate?: Date;
    providerMap?: Map<string, Provider>;
    providersFilter: string[];
    recurrenciesFilter: SubscriptionRecurrency[];
    users: string[];
    usersOptions?: MultiSelectOption[];
    withInactive: boolean;
    onClear: () => void;
    onApply: (values: SubscriptionsFiltersValues) => void;
}

// Format a Date into yyyy-MM-dd for <input type="date">
function formatDateForInput(d?: Date): string {
    if (!d) return "";
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, "0");
    const day = String(d.getDate()).padStart(2, "0");
    return `${year}-${month}-${day}`;
}

// Parse yyyy-MM-dd into a Date at start of day
function parseDateFromInput(s?: string): Date | undefined {
    if (!s) return undefined;
    // Create date at local midnight to avoid timezone shifts
    const [y, m, d] = s.split("-").map(Number);
    if (!y || !m || !d) return undefined;
    return new Date(y, m - 1, d, 0, 0, 0, 0);
}

export const SubscriptionsFilters: React.FC<SubscriptionsFiltersProps> = ({
                                                                              open,
                                                                              onOpenChange,
                                                                              fromDate,
                                                                              toDate,
                                                                              providerMap,
                                                                              providersFilter,
                                                                              recurrenciesFilter,
                                                                              users,
                                                                              usersOptions,
                                                                              withInactive,
                                                                              onClear,
                                                                              onApply,
                                                                          }) => {
    // Local state mirrors props but is editable within the sheet
    const [fromDateStr, setFromDateStr] = useState<string>(formatDateForInput(fromDate));
    const [toDateStr, setToDateStr] = useState<string>(formatDateForInput(toDate));
    const [providersSel, setProvidersSel] = useState<string[]>(providersFilter ?? []);
    const [recurrenciesSel, setRecurrenciesSel] = useState<SubscriptionRecurrency[]>(recurrenciesFilter ?? []);
    const [usersCsv, setUsersCsv] = useState<string>(users?.join(",") ?? "");
    const [withInactiveState, setWithInactiveState] = useState<boolean>(!!withInactive);

    // Keep local state in sync when parent props change
    useEffect(() => {
        setFromDateStr(formatDateForInput(fromDate));
    }, [fromDate]);

    useEffect(() => {
        setToDateStr(formatDateForInput(toDate));
    }, [toDate]);

    useEffect(() => {
        setProvidersSel(providersFilter ?? []);
    }, [providersFilter]);

    useEffect(() => {
        setRecurrenciesSel(recurrenciesFilter ?? []);
    }, [recurrenciesFilter]);

    useEffect(() => {
        setUsersCsv(users?.join(",") ?? "");
    }, [users]);

    useEffect(() => {
        setWithInactiveState(!!withInactive);
    }, [withInactive]);

    // Handlers
    const onFromDateChange = (value: string) => setFromDateStr(value);
    const onToDateChange = (value: string) => setToDateStr(value);

    const onToggleProvider = (value: string) =>
        setProvidersSel((prev) => (prev.includes(value) ? prev.filter((v) => v !== value) : [...prev, value]));

    const onToggleRecurrency = (value: SubscriptionRecurrency) =>
        setRecurrenciesSel((prev) =>
            prev.includes(value) ? prev.filter((v) => v !== value) : [...prev, value]
        );

    const onUsersCsvChange = (value: string) => setUsersCsv(value);
    const onWithInactiveChange = (v: boolean) => setWithInactiveState(!!v);

    const handleClearClick = () => {
        setFromDateStr("");
        setToDateStr("");
        setProvidersSel([]);
        setRecurrenciesSel([]);
        setUsersCsv("");
        setWithInactiveState(false);
        onClear();
    };

    const handleApplyClick = () => {
        const values: SubscriptionsFiltersValues = {
            fromDate: parseDateFromInput(fromDateStr),
            toDate: parseDateFromInput(toDateStr),
            providers: providersSel,
            recurrenciesFilter: recurrenciesSel,
            users: usersCsv ? usersCsv.split(",").map((s) => s.trim()).filter(Boolean) : [],
            withInactive: withInactiveState,
        };
        onApply(values);
    };

    return (
        <Sheet open={open} onOpenChange={onOpenChange}>
            <SheetContent side="right" className="sm:max-w-md">
                <SheetHeader>
                    <SheetTitle>Filters</SheetTitle>
                </SheetHeader>
                <div className="p-4 space-y-6 overflow-y-auto">
                    <div className="flex items-center justify-between">
                        <Label htmlFor="withInactive">Include inactive</Label>
                        <Switch
                            id="withInactive"
                            checked={withInactiveState}
                            onCheckedChange={(v) => onWithInactiveChange(!!v)}
                        />
                    </div>
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

                    {providerMap && providerMap.size > 0 && (
                        <div className="space-y-2">
                            <div className="font-medium">Providers</div>
                            <MultiSelect
                                options={Array.from(providerMap.values()).map((p) => ({value: p.id, label: p.name}))}
                                selectedValues={providersSel}
                                onToggle={onToggleProvider}
                                placeholder="Select providers"
                                searchPlaceholder="Search providers..."
                                emptyMessage="No provider found."
                            />
                        </div>
                    )}

                    <div className="space-y-2">
                        <div className="font-medium">Recurrencies</div>
                        <MultiSelect
                            options={(Object.values(SubscriptionRecurrency) as SubscriptionRecurrency[]).map((r) => ({
                                value: r,
                                label: r
                            }))}
                            selectedValues={recurrenciesSel as string[]}
                            onToggle={(v) => onToggleRecurrency(v as SubscriptionRecurrency)}
                            placeholder="Select recurrencies"
                            searchPlaceholder="Search recurrencies..."
                            emptyMessage="No recurrencies found."
                        />
                    </div>

                    <div className="space-y-2">
                        <div className="font-medium">Family members</div>
                        {usersOptions && usersOptions.length > 0 ? (
                            <MultiSelect
                                options={usersOptions}
                                selectedValues={usersCsv ? usersCsv.split(",").map((s) => s.trim()).filter(Boolean) : []}
                                onChange={(values) => onUsersCsvChange(values.join(","))}
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


                </div>
                <SheetFooter>
                    <div className="flex w-full gap-2">
                        <Button variant="secondary" className="flex-1 cursor-pointer" onClick={handleClearClick}>
                            Clear
                        </Button>
                        <Button className="flex-1 cursor-pointer" onClick={handleApplyClick}>
                            Apply
                        </Button>
                    </div>
                </SheetFooter>
            </SheetContent>
        </Sheet>
    );
};

export default SubscriptionsFilters;

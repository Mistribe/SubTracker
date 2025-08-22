import type TopProvider from "@/models/topProvider.ts";
import type {Amount} from "@/models/amount.ts";
import type UpcomingRenewal from "@/models/upcomingRenewals.ts";
import type TopLabel from "@/models/topLabel.ts";

export default interface Summary {
    activeSubscriptions: number;
    topProviders: TopProvider[];
    topLabels: TopLabel[];
    totalMonthly: Amount;
    totalYearly: Amount;
    totalLastMonth?: Amount;
    totalLastYear?: Amount;
    upcomingRenewals: UpcomingRenewal[];
}


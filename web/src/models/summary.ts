import type TopProvider from "@/models/topProvider.ts";
import type {Amount} from "@/models/amount.ts";
import type UpcomingRenewal from "@/models/upcomingRenewals.ts";
import type TopLabel from "@/models/topLabel.ts";

export default interface Summary {
    activeSubscriptions: number;
    activePersonal: number;
    activeFamily: number;
    topProviders: TopProvider[];
    topLabels: TopLabel[];
    totalMonthly: Amount;
    totalYearly: Amount;
    totalLastMonth?: Amount;
    totalLastYear?: Amount;
    personalMonthly?: Amount;
    personalLastMonth?: Amount;
    personalYearly?: Amount;
    personalLastYear?: Amount;
    familyMonthly?: Amount;
    familyLastMonth?: Amount;
    familyYearly?: Amount;
    familyLastYear?: Amount;
    upcomingRenewals: UpcomingRenewal[];
}


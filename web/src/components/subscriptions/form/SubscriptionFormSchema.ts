import * as z from "zod";
import {OwnerType} from "@/models/ownerType";
import {SubscriptionRecurrency} from "@/models/subscriptionRecurrency";

// Define the form schema with Zod
export const formSchema = z.object({
    friendlyName: z.string().optional(),
    providerId: z.string().min(1, "Provider is required"),
    planId: z.string().optional(),
    priceId: z.string().optional(),
    recurrency: z.enum(SubscriptionRecurrency),
    customRecurrencyValue: z.number().positive("Value must be positive").optional(),
    customRecurrencyUnit: z.enum(["days", "months", "years"]).optional(),
    startDate: z.date(),
    endDate: z.date().optional(),
    ownerType: z.enum(OwnerType),
    familyId: z.string().optional(),
    serviceUsers: z.array(z.string()).optional(),
    payerType: z.enum(["family", "family_member"]).optional(),
    payerId: z.string().optional(),
    customPrice: z.object({
        amount: z.number().positive("Amount must be positive"),
        currency: z.string().min(1, "Currency is required")
    }),
    hasFreeTrialPeriod: z.boolean(),
    freeTrialStartDate: z.date().optional(),
    freeTrialEndDate: z.date().optional(),
});

export type FormValues = z.infer<typeof formSchema>;

// Helper function to convert custom recurrency to days
export const convertToDays = (value: number, unit: string): number => {
    switch (unit) {
        case "days":
            return value;
        case "months":
            return value * 30; // Approximate days in a month
        case "years":
            return value * 365; // Approximate days in a year
        default:
            return value;
    }
};
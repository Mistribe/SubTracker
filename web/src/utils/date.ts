import {addMonths, differenceInMonths} from "date-fns";

export interface MonthsBetweenOptions {
    /**
     * If true (default), any remaining partial month counts as +1.
     * If false, only full months are counted.
     */
    includePartial?: boolean;
}

/**
 * Returns the number of months between two dates.
 * - Counts full calendar months between `start` and `end`.
 * - If `includePartial` is true and there's any remainder beyond full months, it rounds up by one.
 * - Returns 0 if `end` <= `start`.
 */
export function monthsBetween(start: Date, end: Date, options: MonthsBetweenOptions = {}): number {
    const {includePartial = true} = options;

    if (!(start instanceof Date) || !(end instanceof Date)) {
        throw new TypeError("monthsBetween expects Date instances");
    }
    if (isNaN(start.getTime()) || isNaN(end.getTime())) {
        throw new RangeError("Invalid Date provided");
    }
    if (end <= start) {
        return 0;
    }

    const fullMonths = differenceInMonths(end, start);
    const boundary = addMonths(start, fullMonths);
    const hasRemainder = end > boundary;

    return includePartial && hasRemainder ? fullMonths + 1 : fullMonths;
}
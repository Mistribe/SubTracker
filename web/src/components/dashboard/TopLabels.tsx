import {Skeleton} from "@/components/ui/skeleton";
import {Tag} from "lucide-react";
import Money from "@/components/ui/money.tsx";
import type TopLabel from "@/models/topLabel.ts";
import {useLabelQuery} from "@/hooks/labels/useLabelQuery";

interface TopLabelsProps {
    labels: TopLabel[];
    isLoading: boolean;
}

const LabelName = ({id}: { id: string }) => {
    const {data} = useLabelQuery(id);
    return <h4 className="font-medium">{data?.name ?? id}</h4>;
};

const TopLabels = ({labels, isLoading}: TopLabelsProps) => {
    return (
        <div>
            <div className="flex flex-row items-center justify-between pb-5 space-y-0">
                <p className="text-xl font-medium">Top Spent by Label</p>
                <Tag className="h-5 w-5 text-blue-500 animate-pulse"/>
            </div>
            <div>
                {isLoading ? (
                    <div className="space-y-3">
                        {Array.from({length: 3}).map((_, i) => (
                            <div key={`toplabel-skeleton-${i}`} className="p-3 border rounded-lg bg-muted/30">
                                <Skeleton className="h-6 w-full mb-2"/>
                                <Skeleton className="h-4 w-1/2"/>
                            </div>
                        ))}
                    </div>
                ) : labels.length > 0 ? (
                    <div className="space-y-3">
                        {labels.map((label) => (
                            <div key={label.labelId}
                                 className="p-3 border rounded-lg bg-card transition-all duration-300 hover:shadow-lg ">
                                <div className="flex justify-between items-start">
                                    <div className="flex items-center">
                                        <LabelName id={label.labelId}/>
                                    </div>
                                    <span
                                        className="font-semibold bg-gradient-to-r from-blue-500 to-blue-700 bg-clip-text text-transparent">
                                            <Money amount={label.total}/>
                                        </span>
                                </div>
                            </div>
                        ))}
                    </div>
                ) : (
                    <p className="text-muted-foreground">No label spending data available.</p>
                )}
            </div>
        </div>
    );
};

export default TopLabels;

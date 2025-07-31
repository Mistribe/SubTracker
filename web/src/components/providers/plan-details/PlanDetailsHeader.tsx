import {
    DialogDescription,
    DialogHeader,
    DialogTitle
} from "@/components/ui/dialog";
import Plan from "@/models/plan";

interface PlanDetailsHeaderProps {
    plan: Plan;
}

export function PlanDetailsHeader({ plan }: PlanDetailsHeaderProps) {
    return (
        <DialogHeader>
            <DialogTitle>{plan.name}</DialogTitle>
            {plan.description && (
                <DialogDescription>
                    {plan.description}
                </DialogDescription>
            )}
        </DialogHeader>
    );
}
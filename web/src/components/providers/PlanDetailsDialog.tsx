// Re-export the PlanDetailsDialog component from the plan-details directory
export { PlanDetailsDialog } from './plan-details';

// Keep the original interface for backward compatibility
import Plan from "@/models/plan";

interface PlanDetailsDialogProps {
    isOpen: boolean;
    onClose: () => void;
    plan: Plan;
    providerId: string;
    onPlanUpdated?: () => void;
}

// The implementation has been moved to ./plan-details/PlanDetailsDialog.tsx
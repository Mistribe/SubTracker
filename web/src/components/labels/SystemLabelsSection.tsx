import Label from "@/models/label";
import { LabelItem } from "./LabelItem";

interface SystemLabelsSectionProps {
  labels: Label[];
}

export const SystemLabelsSection = ({ labels }: SystemLabelsSectionProps) => {
  // Filter to only include system labels
  const systemLabels = labels.filter(label => label.owner.isSystem);

  return (
    <div className="mb-8">
      <h2 className="text-xl font-semibold mb-4">System Labels</h2>
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3">
        {systemLabels.map((label) => (
          <LabelItem
            key={label.id}
            label={label}
            isReadOnly={true}
          />
        ))}
      </div>
    </div>
  );
};
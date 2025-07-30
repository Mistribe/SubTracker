import { CreateFamilyDialog } from "./CreateFamilyDialog";

export const EmptyFamiliesState = () => {
  return (
    <div className="text-center p-8 bg-gray-50 rounded-lg">
      <p className="text-gray-500 mb-4">No families found. Create your first family to get started.</p>
      <div className="flex justify-center">
        <CreateFamilyDialog />
      </div>
    </div>
  );
};
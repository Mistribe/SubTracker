import { OwnerType } from "@/models/ownerType";

// Get badge color based on owner type
export const getBadgeVariant = (ownerType: OwnerType) => {
  switch (ownerType) {
    case OwnerType.System:
      return "secondary";
    case OwnerType.Family:
      return "outline";
    case OwnerType.Personal:
      return "default";
    default:
      return "default";
  }
};

// Get badge text based on owner type
export const getBadgeText = (ownerType: OwnerType) => {
  switch (ownerType) {
    case OwnerType.System:
      return "System";
    case OwnerType.Family:
      return "Family";
    case OwnerType.Personal:
      return "Personal";
    default:
      return "Unknown";
  }
};
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import Family from "@/models/family.ts";
import { FamilyHeader } from "./FamilyHeader";
import { FamilyMembersTable } from "./FamilyMembersTable";

interface FamilyCardProps {
  family: Family;
}

export const FamilyCard = ({ family }: FamilyCardProps) => {
  return (
    <Card>
      <CardHeader>
        <FamilyHeader family={family} />
      </CardHeader>
      <CardContent>
        <FamilyMembersTable family={family} />
      </CardContent>
    </Card>
  );
};
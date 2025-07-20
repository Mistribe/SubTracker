import { useState } from "react";
import { 
  Card, 
  CardContent, 
  CardDescription, 
  CardHeader, 
  CardTitle 
} from "@/components/ui/card";
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow 
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { PlusIcon } from "lucide-react";

interface FamilyMember {
  id: number;
  name: string;
  role: string;
  email: string;
}

interface Family {
  id: number;
  name: string;
  members: FamilyMember[];
}

const FamiliesPage = () => {
  // Sample data - in a real app, this would come from an API
  const [families] = useState<Family[]>([
    {
      id: 1,
      name: "Smith Family",
      members: [
        {
          id: 1,
          name: "John Smith",
          role: "Admin",
          email: "john.smith@example.com"
        },
        {
          id: 2,
          name: "Jane Smith",
          role: "Member",
          email: "jane.smith@example.com"
        },
        {
          id: 3,
          name: "Jimmy Smith",
          role: "Member",
          email: "jimmy.smith@example.com"
        }
      ]
    },
    {
      id: 2,
      name: "Johnson Family",
      members: [
        {
          id: 4,
          name: "Mike Johnson",
          role: "Admin",
          email: "mike.johnson@example.com"
        },
        {
          id: 5,
          name: "Sarah Johnson",
          role: "Member",
          email: "sarah.johnson@example.com"
        }
      ]
    }
  ]);

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Families</h1>
        <Button>
          <PlusIcon className="h-4 w-4 mr-2" />
          Add Family
        </Button>
      </div>

      <div className="grid gap-6">
        {families.map((family) => (
          <Card key={family.id}>
            <CardHeader>
              <div className="flex justify-between items-center">
                <div>
                  <CardTitle>{family.name}</CardTitle>
                  <CardDescription>{family.members.length} members</CardDescription>
                </div>
                <div className="flex gap-2">
                  <Button variant="outline" size="sm">Edit</Button>
                  <Button variant="outline" size="sm">
                    <PlusIcon className="h-4 w-4 mr-2" />
                    Add Member
                  </Button>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Name</TableHead>
                    <TableHead>Role</TableHead>
                    <TableHead>Email</TableHead>
                    <TableHead className="w-[100px]">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {family.members.map((member) => (
                    <TableRow key={member.id}>
                      <TableCell className="font-medium">{member.name}</TableCell>
                      <TableCell>{member.role}</TableCell>
                      <TableCell>{member.email}</TableCell>
                      <TableCell>
                        <Button variant="ghost" size="sm">Edit</Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
};

export default FamiliesPage;
import {Table, TableBody, TableHead, TableHeader, TableRow} from "@/components/ui/table";
import Family from "@/models/family.ts";
import {FamilyMemberRow} from "./FamilyMemberRow";

interface FamilyMembersTableProps {
    family: Family;
}

export const FamilyMembersTable = ({family}: FamilyMembersTableProps) => {
    return (
        <div className="overflow-hidden rounded-lg border">

            <Table>
                <TableHeader className="bg-muted sticky">
                    <TableRow>
                        <TableHead>Name</TableHead>
                        <TableHead>Role</TableHead>
                        <TableHead className="w-[100px]">Actions</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {family.members.map((member) => (
                        <FamilyMemberRow
                            key={member.id}
                            member={member}
                            familyId={family.id}
                            isOwner={family.isOwner}
                            className="bg-white dark:bg-black"
                        />
                    ))}
                </TableBody>
            </Table>
        </div>
    );
};
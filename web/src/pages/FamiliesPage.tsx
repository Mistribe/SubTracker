import {useState} from "react";
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from "@/components/ui/card";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table";
import {Button} from "@/components/ui/button";
import {CheckIcon, Loader2, XIcon} from "lucide-react";
import Family from "@/models/family.ts";
import FamilyMember from "@/models/familyMember.ts";
import {useApiClient} from "@/hooks/use-api-client.ts";
import {useQuery, useQueryClient} from "@tanstack/react-query";
import type {FamilyModel, PatchFamilyModel, UpdateFamilyMemberModel} from "@/api/models";
import {CreateFamilyDialog} from "@/components/families/CreateFamilyDialog.tsx";
import {Input} from "@/components/ui/input";
import {AddFamilyMemberDialog} from "@/components/families/AddFamilyMemberDialog.tsx";
import {Checkbox} from "@/components/ui/checkbox";

const FamiliesPage = () => {
    const {apiClient} = useApiClient();
    const queryClient = useQueryClient();

    const [offset] = useState(0);
    const [limit] = useState(10);
    const [editingFamilyId, setEditingFamilyId] = useState<string | null>(null);
    const [editedName, setEditedName] = useState<string>("");
    const [isUpdating, setIsUpdating] = useState<boolean>(false);

    // Family member editing state
    const [editingMemberId, setEditingMemberId] = useState<string | null>(null);
    const [editedMemberName, setEditedMemberName] = useState<string>("");
    const [editedMemberIsKid, setEditedMemberIsKid] = useState<boolean>(false);
    const [isUpdatingMember, setIsUpdatingMember] = useState<boolean>(false);
    const [isDeletingMember, setIsDeletingMember] = useState<boolean>(false);

    // Function to start editing a family
    const startEditing = (family: Family) => {
        setEditingFamilyId(family.id);
        setEditedName(family.name);
    };

    // Function to cancel editing
    const cancelEditing = () => {
        setEditingFamilyId(null);
        setEditedName("");
    };

    // Function to start editing a family member
    const startEditingMember = (member: FamilyMember) => {
        setEditingMemberId(member.id);
        setEditedMemberName(member.name);
        setEditedMemberIsKid(member.isKid);
    };

    // Function to cancel editing a family member
    const cancelEditingMember = () => {
        setEditingMemberId(null);
        setEditedMemberName("");
        setEditedMemberIsKid(false);
    };

    // Function to save family member changes
    const saveMemberChanges = async (familyId: string, member: FamilyMember) => {
        if (!apiClient) return;

        try {
            setIsUpdatingMember(true);

            // Only update if values have changed
            if (editedMemberName === member.name && editedMemberIsKid === member.isKid) {
                cancelEditingMember();
                return;
            }

            const patchModel: Partial<UpdateFamilyMemberModel> = {
                name: editedMemberName,
                isKid: editedMemberIsKid,
            };

            await apiClient.families.byFamilyId(familyId).members.byId(member.id).put(patchModel);

            // Invalidate and refetch the families query
            await queryClient.invalidateQueries({queryKey: ['families']});

            // Reset editing state
            cancelEditingMember();
        } catch (error) {
            console.error("Failed to update family member:", error);
        } finally {
            setIsUpdatingMember(false);
        }
    };

    // Function to remove a family member
    const removeMember = async (familyId: string, member: FamilyMember, isOwner: boolean) => {
        if (!apiClient) return;

        // Prevent removing yourself if you are the owner of the family
        if (member.isYou && isOwner) {
            console.error("Cannot remove yourself as the owner of the family");
            return;
        }

        try {
            setIsDeletingMember(true);

            await apiClient.families.byFamilyId(familyId).members.byId(member.id).delete();

            // Invalidate and refetch the families query
            await queryClient.invalidateQueries({queryKey: ['families']});
        } catch (error) {
            console.error("Failed to remove family member:", error);
        } finally {
            setIsDeletingMember(false);
        }
    };

    // Function to save changes
    const saveChanges = async (family: Family) => {
        if (!apiClient) return;

        try {
            setIsUpdating(true);

            // Only update if values have changed
            if (editedName === family.name) {
                cancelEditing();
                return;
            }

            const patchModel: Partial<PatchFamilyModel> = {
                id: family.id,
                name: editedName,
                updatedAt: new Date()
            };

            await apiClient.families.patch(patchModel);

            // Invalidate and refetch the families query
            await queryClient.invalidateQueries({queryKey: ['families']});

            // Reset editing state
            cancelEditing();
        } catch (error) {
            console.error("Failed to update family:", error);
        } finally {
            setIsUpdating(false);
        }
    };

    const {
        data: queryResponse,
        isLoading,
        error
    } = useQuery({
        queryKey: ['families'],
        queryFn: async () => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }
            const result = await apiClient.families.get({
                queryParameters: {
                    offset: offset,
                    limit: limit
                }
            });
            if (result && result.data) {
                return {
                    families: result.data.map((model: FamilyModel) => {
                        return Family.fromModel(model);
                    }),
                    length: result.data.length,
                    total: result.total ?? 0
                }
            }
            return {families: [], length: 0, total: 0};
        },
        enabled: !!apiClient,
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchOnWindowFocus: true,
    });

    if (isLoading) {
        return (
            <div className="flex justify-center items-center h-64">
                <Loader2 className="h-8 w-8 animate-spin"/>
                <span className="ml-2">Loading families...</span>
            </div>
        );
    }

    if (error) {
        return (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded relative" role="alert">
                <strong className="font-bold">Error!</strong>
                <span
                    className="block sm:inline"> {error instanceof Error ? error.message : 'Failed to load families'}</span>
            </div>
        );
    }

    const families = queryResponse?.families || [];

    return (
        <div>
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-3xl font-bold">Families</h1>
                {families.length > 0 && families.filter(x => x.isOwner).length === 0 && <CreateFamilyDialog/>}
            </div>

            {families.length === 0 ? (
                <div className="text-center p-8 bg-gray-50 rounded-lg">
                    <p className="text-gray-500 mb-4">No families found. Create your first family to get started.</p>
                    <div className="flex justify-center">
                        <CreateFamilyDialog/>
                    </div>
                </div>
            ) : (
                <div className="grid gap-6">
                    {families.map((family) => (
                        <Card key={family.id}>
                            <CardHeader>
                                <div className="flex justify-between items-center">
                                    <div>
                                        {editingFamilyId === family.id ? (
                                            <div className="flex flex-col gap-2">
                                                <div className="flex items-center gap-2">
                                                    <Input
                                                        value={editedName}
                                                        onChange={(e) => setEditedName(e.target.value)}
                                                        className="w-full"
                                                        placeholder="Family name"
                                                    />
                                                </div>
                                            </div>
                                        ) : (
                                            <>
                                                <div className="flex items-center gap-2">
                                                    <CardTitle>{family.name}</CardTitle>
                                                </div>
                                                <CardDescription>{family.members.length} members</CardDescription>
                                            </>
                                        )}
                                    </div>
                                    <div className="flex gap-2">
                                        {editingFamilyId === family.id ? (
                                            <>
                                                <Button
                                                    variant="outline"
                                                    size="sm"
                                                    onClick={() => saveChanges(family)}
                                                    disabled={isUpdating}
                                                >
                                                    {isUpdating ? (
                                                        <Loader2 className="h-4 w-4 animate-spin mr-2"/>
                                                    ) : (
                                                        <CheckIcon className="h-4 w-4 mr-2"/>
                                                    )}
                                                    Save
                                                </Button>
                                                <Button
                                                    variant="outline"
                                                    size="sm"
                                                    onClick={cancelEditing}
                                                    disabled={isUpdating}
                                                >
                                                    <XIcon className="h-4 w-4 mr-2"/>
                                                    Cancel
                                                </Button>
                                            </>
                                        ) : (
                                            <>
                                                <Button
                                                    variant="outline"
                                                    size="sm"
                                                    onClick={() => startEditing(family)}
                                                >
                                                    Edit
                                                </Button>
                                                <AddFamilyMemberDialog familyId={family.id}/>
                                            </>
                                        )}
                                    </div>
                                </div>
                            </CardHeader>
                            <CardContent>
                                <Table>
                                    <TableHeader>
                                        <TableRow>
                                            <TableHead>Name</TableHead>
                                            <TableHead>Role</TableHead>
                                            <TableHead className="w-[100px]">Actions</TableHead>
                                        </TableRow>
                                    </TableHeader>
                                    <TableBody>
                                        {family.members.map((member) => (
                                            <TableRow key={member.id}>
                                                <TableCell className="font-medium">
                                                    {editingMemberId === member.id ? (
                                                        <Input
                                                            value={editedMemberName}
                                                            onChange={(e) => setEditedMemberName(e.target.value)}
                                                            className="w-full"
                                                            placeholder="Member name"
                                                        />
                                                    ) : (
                                                        member.isYou ? (
                                                            <span>{member.name} <i>(You)</i></span>
                                                        ) : (
                                                            <span>{member.name}</span>
                                                        )
                                                    )}
                                                </TableCell>
                                                <TableCell>
                                                    {editingMemberId === member.id ? (
                                                        <div className="flex items-center space-x-2">
                                                            <Checkbox
                                                                checked={editedMemberIsKid}
                                                                onCheckedChange={(checked) => setEditedMemberIsKid(checked === true)}
                                                                id={`kid-checkbox-${member.id}`}
                                                            />
                                                            <label htmlFor={`kid-checkbox-${member.id}`}>Kid</label>
                                                        </div>
                                                    ) : (
                                                        member.isKid ? 'Kid' : 'Adult'
                                                    )}
                                                </TableCell>
                                                <TableCell>
                                                    {editingMemberId === member.id ? (
                                                        <div className="flex gap-2">
                                                            <Button
                                                                variant="outline"
                                                                size="sm"
                                                                onClick={() => saveMemberChanges(family.id, member)}
                                                                disabled={isUpdatingMember}
                                                            >
                                                                {isUpdatingMember ? (
                                                                    <Loader2 className="h-4 w-4 animate-spin mr-2"/>
                                                                ) : (
                                                                    <CheckIcon className="h-4 w-4 mr-2"/>
                                                                )}
                                                                Save
                                                            </Button>
                                                            <Button
                                                                variant="outline"
                                                                size="sm"
                                                                onClick={cancelEditingMember}
                                                                disabled={isUpdatingMember}
                                                            >
                                                                <XIcon className="h-4 w-4 mr-2"/>
                                                                Cancel
                                                            </Button>
                                                        </div>
                                                    ) : (
                                                        <div className="flex gap-2">
                                                            <Button
                                                                variant="ghost"
                                                                size="sm"
                                                                onClick={() => startEditingMember(member)}
                                                            >
                                                                Edit
                                                            </Button>
                                                            {/* Don't show remove button for yourself if you're the owner */}
                                                            {!(member.isYou && family.isOwner) && (
                                                                <Button
                                                                    variant="ghost"
                                                                    size="sm"
                                                                    onClick={() => removeMember(family.id, member, family.isOwner)}
                                                                    disabled={isDeletingMember}
                                                                    className="text-red-500 hover:text-red-700 hover:bg-red-50"
                                                                >
                                                                    {isDeletingMember ? (
                                                                        <Loader2 className="h-4 w-4 animate-spin mr-2"/>
                                                                    ) : (
                                                                        <XIcon className="h-4 w-4 mr-2"/>
                                                                    )}
                                                                    Remove
                                                                </Button>
                                                            )}
                                                        </div>
                                                    )}
                                                </TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </CardContent>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    );
};

export default FamiliesPage;
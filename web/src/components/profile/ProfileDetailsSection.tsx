import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

interface ProfileDetailsProps {
  givenName: string;
  familyName: string;
  email: string;
  joinDate: string;
  onSave: (data: { givenName: string; familyName: string; email: string }) => void;
  onCancel?: () => void;
  isLoading?: boolean;
}

/**
 * ProfileDetailsSection component
 * Displays and manages the editing of user profile details
 */
export const ProfileDetailsSection = ({
  givenName,
  familyName,
  email,
  joinDate,
  onSave,
  onCancel,
  isEditing = false,
  isLoading = false,
}: ProfileDetailsProps & { isEditing?: boolean }) => {
  const [formData, setFormData] = useState({ givenName, familyName, email });
  const [originalData] = useState({ givenName, familyName, email });

  // Update form data when props change
  useEffect(() => {
    setFormData({ givenName, familyName, email });
  }, [givenName, familyName, email]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  // Check if form data has changed from original values
  const hasChanges = formData.givenName !== originalData.givenName || 
                     formData.familyName !== originalData.familyName;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (hasChanges) {
      onSave(formData);
    }
  };

  return (
    <>
      {isEditing ? (
        <form onSubmit={handleSubmit} className="space-y-6 bg-muted/20 p-6 rounded-lg">
          <div className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label htmlFor="givenName" className="block text-sm font-medium mb-1">
                  First Name
                </label>
                <Input
                  id="givenName"
                  name="givenName"
                  type="text"
                  value={formData.givenName}
                  onChange={handleInputChange}
                  required
                  className="w-full"
                />
              </div>
              <div>
                <label htmlFor="familyName" className="block text-sm font-medium mb-1">
                  Last Name
                </label>
                <Input
                  id="familyName"
                  name="familyName"
                  type="text"
                  value={formData.familyName}
                  onChange={handleInputChange}
                  required
                  className="w-full"
                />
              </div>
            </div>

            <div>
              <label htmlFor="email" className="block text-sm font-medium mb-1">
                Email
                <span className="ml-2 text-xs text-muted-foreground">(managed by authentication provider)</span>
              </label>
              <Input
                id="email"
                name="email"
                type="email"
                value={formData.email}
                onChange={handleInputChange}
                required
                disabled
                className="w-full opacity-70 cursor-not-allowed"
              />
              <p className="text-xs text-muted-foreground mt-1">Email address cannot be changed through this application</p>
            </div>
          </div>

          <div className="flex justify-end gap-2 mt-6">
            <Button
              type="button"
              variant="outline"
              onClick={() => {
                setFormData({ givenName, familyName, email });
                if (onCancel) {
                  onCancel();
                }
              }}
            >
              Cancel
            </Button>
            <Button 
              type="submit" 
              disabled={isLoading || !hasChanges}
              title={!hasChanges ? "No changes to save" : ""}
            >
              {isLoading ? (
                <>
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-current"></div>
                  Saving...
                </>
              ) : (
                'Save Changes'
              )}
            </Button>
          </div>
        </form>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <h3 className="text-sm font-medium text-muted-foreground">First Name</h3>
            <p className="mt-1 font-medium">{givenName}</p>
          </div>
          <div>
            <h3 className="text-sm font-medium text-muted-foreground">Last Name</h3>
            <p className="mt-1 font-medium">{familyName}</p>
          </div>
          <div>
            <h3 className="text-sm font-medium text-muted-foreground">Email</h3>
            <p className="mt-1 font-medium">{email}</p>
          </div>
          <div>
            <h3 className="text-sm font-medium text-muted-foreground">Member Since</h3>
            <p className="mt-1 font-medium">{joinDate}</p>
          </div>
        </div>
      )}
    </>
  );
};

/**
 * ProfileDetailsContainer component
 * Container component that includes the edit button and the ProfileDetailsSection
 */
export const ProfileDetailsContainer = ({
  givenName,
  familyName,
  email,
  joinDate,
  onSave,
  isLoading = false,
}: ProfileDetailsProps) => {
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState({ givenName, familyName, email });
  const [isSaving, setIsSaving] = useState(false);
  const [originalData] = useState({ givenName, familyName, email });

  // Check if data has changed from original values
  const hasDataChanged = (data: { givenName: string; familyName: string; email: string }) => {
    return data.givenName !== originalData.givenName || data.familyName !== originalData.familyName;
  };

  const handleSave = (data: { givenName: string; familyName: string; email: string }) => {
    setFormData(data);
    
    // Only call onSave and set isSaving if data has actually changed
    if (hasDataChanged(data)) {
      setIsSaving(true);
      onSave(data);
    } else {
      // If no changes, just close the editing mode without calling the backend
      setIsEditing(false);
    }
  };

  // Close editing mode only after backend response is received
  useEffect(() => {
    if (isSaving && !isLoading) {
      // Backend has responded (isLoading changed from true to false)
      setIsSaving(false);
      setIsEditing(false);
    }
  }, [isLoading, isSaving]);

  return (
    <div className="space-y-8">
      <div className="flex justify-end">
        {!isEditing && (
          <Button onClick={() => setIsEditing(true)}>Edit Profile</Button>
        )}
      </div>
      
      <ProfileDetailsSection
        givenName={isEditing ? formData.givenName : givenName}
        familyName={isEditing ? formData.familyName : familyName}
        email={isEditing ? formData.email : email}
        joinDate={joinDate}
        isEditing={isEditing}
        isLoading={isLoading}
        onSave={(data) => {
          handleSave(data);
        }}
        onCancel={() => {
          setIsEditing(false);
        }}
      />
    </div>
  );
};
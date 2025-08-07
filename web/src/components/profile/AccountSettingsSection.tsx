import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";

interface AccountSettingsSectionProps {
  onEmailNotificationsToggle?: () => void;
  onChangePassword?: () => void;
  onDeleteAccount?: () => void;
}

/**
 * AccountSettingsSection component
 * Displays account settings options like email notifications, password change, and account deletion
 */
export const AccountSettingsSection = ({
  onEmailNotificationsToggle,
  onChangePassword,
  onDeleteAccount,
}: AccountSettingsSectionProps) => {
  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-bold">Account Settings</h2>
        <p className="text-muted-foreground">Manage your account settings and preferences</p>
      </div>

      <div className="space-y-6 bg-muted/20 p-6 rounded-lg">
        <div className="flex justify-between items-center">
          <div>
            <h4 className="font-medium">Email Notifications</h4>
            <p className="text-sm text-muted-foreground">
              Receive email notifications for upcoming payments
            </p>
          </div>
          <Button 
            variant="outline" 
            size="sm" 
            onClick={onEmailNotificationsToggle}
          >
            Enabled
          </Button>
        </div>

        <Separator />

        <div className="flex justify-between items-center">
          <div>
            <h4 className="font-medium">Change Password</h4>
            <p className="text-sm text-muted-foreground">Update your account password</p>
          </div>
          <Button 
            variant="outline" 
            size="sm" 
            onClick={onChangePassword}
          >
            Change
          </Button>
        </div>

        <Separator />

        <div className="flex justify-between items-center">
          <div>
            <h4 className="font-medium text-destructive">Delete Account</h4>
            <p className="text-sm text-muted-foreground">
              Permanently delete your account and all data
            </p>
          </div>
          <Button 
            variant="outline" 
            size="sm" 
            className="border-destructive text-destructive"
            onClick={onDeleteAccount}
          >
            Delete
          </Button>
        </div>
      </div>
    </div>
  );
};
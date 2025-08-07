import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

interface UserProfileSectionProps {
  name: string;
  email: string;
  picture?: string;
}

/**
 * UserProfileSection component
 * Displays the user's avatar and basic information
 */
export const UserProfileSection = ({ name, email, picture }: UserProfileSectionProps) => {
  return (
    <div className="flex flex-col items-center space-y-4 py-8 bg-primary/5 rounded-lg">
      <Avatar className="h-24 w-24 border-4 border-background">
        {picture ? (
          <AvatarImage src={picture} alt={name} />
        ) : (
          <AvatarFallback className="text-2xl">
            {name.split(' ').map(n => n[0]).join('').toUpperCase()}
          </AvatarFallback>
        )}
      </Avatar>
      <div className="text-center">
        <h2 className="text-2xl font-bold">{name}</h2>
        <p className="text-muted-foreground">{email}</p>
      </div>
    </div>
  );
};
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import Provider from "@/models/provider";
import { getBadgeText, getBadgeVariant } from "../utils/badgeUtils";
import { Edit } from "lucide-react";
import { useProvidersMutations } from "@/hooks/providers/useProvidersMutations";

interface ProviderCardProps {
  provider: Provider;
  onEdit: (provider: Provider) => void;
}

export const ProviderCard = ({ provider, onEdit }: ProviderCardProps) => {
  const { canModifyProvider } = useProvidersMutations();
  const isEditable = canModifyProvider(provider);

  return (
    <Card key={provider.id} className="overflow-hidden">
      <CardHeader className="pb-2">
        <div className="flex justify-between items-start">
          <CardTitle>{provider.name}</CardTitle>
          <div className="flex items-center space-x-2">
            {isEditable && (
              <Button 
                variant="ghost" 
                size="sm" 
                className="h-8 w-8 p-0" 
                onClick={() => onEdit(provider)}
                title="Edit provider"
              >
                <Edit className="h-4 w-4" />
              </Button>
            )}
            <Badge variant={getBadgeVariant(provider.owner.type)}>
              {getBadgeText(provider.owner.type)}
            </Badge>
          </div>
        </div>
        {provider.url && (
          <CardDescription>
            <a href={provider.url} target="_blank" rel="noopener noreferrer" className="text-blue-500 hover:underline">
              {provider.url}
            </a>
          </CardDescription>
        )}
      </CardHeader>
      <CardContent>
        {provider.description && <p>{provider.description}</p>}
        {provider.labels.length > 0 && (
          <div className="flex flex-wrap gap-2 mt-2">
            {provider.labels.map((label, index) => (
              <Badge key={index} variant="outline">{label}</Badge>
            ))}
          </div>
        )}
      </CardContent>
      <CardFooter className="flex justify-between">
        {provider.pricingPageUrl && (
          <a href={provider.pricingPageUrl} target="_blank" rel="noopener noreferrer">
            <Button variant="outline" size="sm">View Pricing</Button>
          </a>
        )}
        {provider.plans.length > 0 && (
          <Badge variant="secondary">{provider.plans.length} Plans</Badge>
        )}
      </CardFooter>
    </Card>
  );
};
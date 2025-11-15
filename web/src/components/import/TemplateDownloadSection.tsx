import { FileSpreadsheet, FileJson, FileCode, Download } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

interface TemplateDownloadSectionProps {
  entityType: 'labels' | 'providers' | 'subscriptions';
}

interface TemplateLink {
  format: 'csv' | 'json' | 'yaml';
  label: string;
  path: string;
  icon: React.ReactNode;
  description: string;
}

export function TemplateDownloadSection({ entityType }: TemplateDownloadSectionProps) {
  const templates: TemplateLink[] = [
    {
      format: 'csv',
      label: 'CSV Template',
      path: `/templates/${entityType}/${entityType}-template.csv`,
      icon: <FileSpreadsheet className="h-4 w-4" aria-hidden="true" />,
      description: 'Comma-separated values format',
    },
    {
      format: 'json',
      label: 'JSON Template',
      path: `/templates/${entityType}/${entityType}-template.json`,
      icon: <FileJson className="h-4 w-4" aria-hidden="true" />,
      description: 'JavaScript Object Notation format',
    },
    {
      format: 'yaml',
      label: 'YAML Template',
      path: `/templates/${entityType}/${entityType}-template.yaml`,
      icon: <FileCode className="h-4 w-4" aria-hidden="true" />,
      description: 'YAML Ain\'t Markup Language format',
    },
  ];

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Download className="h-5 w-5" aria-hidden="true" />
          Download Template Files
        </CardTitle>
        <CardDescription>
          Choose a template format to get started with your {entityType} import. 
          Each template includes example data and all available fields.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
          {templates.map((template) => (
            <Button
              key={template.format}
              variant="outline"
              className="h-auto flex-col items-start gap-2 p-4 hover:bg-accent hover:border-primary transition-colors"
              asChild
            >
              <a
                href={template.path}
                download={`${entityType}-template.${template.format}`}
                aria-label={`Download ${template.label} for ${entityType}`}
              >
                <div className="flex items-center gap-2 w-full">
                  {template.icon}
                  <span className="font-semibold">{template.label}</span>
                </div>
                <span className="text-xs text-muted-foreground text-left w-full">
                  {template.description}
                </span>
              </a>
            </Button>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

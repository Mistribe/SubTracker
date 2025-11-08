import { Info, Download, FileText, AlertCircle } from 'lucide-react';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

interface ImportHelpProps {
  entityType: 'labels' | 'providers' | 'subscriptions';
}

export function ImportHelp({ entityType }: ImportHelpProps) {
  const getTemplateUrl = () => {
    const templates = {
      labels: '/templates/labels-template.csv',
      providers: '/templates/providers-template.json',
      subscriptions: '/templates/subscriptions-template.yaml',
    };
    return templates[entityType];
  };

  const getTemplateFilename = () => {
    const filenames = {
      labels: 'labels-template.csv',
      providers: 'providers-template.json',
      subscriptions: 'subscriptions-template.yaml',
    };
    return filenames[entityType];
  };

  const handleDownloadTemplate = () => {
    const url = getTemplateUrl();
    const filename = getTemplateFilename();
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  return (
    <Card className="mb-6">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Info className="h-5 w-5" />
          Import Help
        </CardTitle>
        <CardDescription>
          Learn how to import {entityType} from files
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <Alert>
          <FileText className="h-4 w-4" />
          <AlertDescription>
            Download the template file to see the correct format and field names for importing {entityType}.
          </AlertDescription>
        </Alert>

        <Button
          variant="outline"
          onClick={handleDownloadTemplate}
          className="w-full sm:w-auto"
        >
          <Download className="mr-2 h-4 w-4" />
          Download Template
        </Button>

        <Accordion type="single" collapsible className="w-full">
          <AccordionItem value="process">
            <AccordionTrigger>Import Process</AccordionTrigger>
            <AccordionContent className="space-y-2 text-sm">
              <ol className="list-decimal list-inside space-y-2">
                <li>Select or drag-and-drop a file (CSV, JSON, or YAML format)</li>
                <li>The file will be parsed and validated in your browser</li>
                <li>Review the parsed records in the preview table</li>
                <li>Select which records you want to import</li>
                <li>Click "Import Selected" or "Import All" to begin importing</li>
                <li>Records will be imported one at a time to the server</li>
                <li>Review the import results and retry any failed imports if needed</li>
              </ol>
            </AccordionContent>
          </AccordionItem>

          <AccordionItem value="formats">
            <AccordionTrigger>Supported File Formats</AccordionTrigger>
            <AccordionContent className="space-y-2 text-sm">
              <ul className="list-disc list-inside space-y-2">
                <li><strong>CSV (.csv)</strong> - Comma-separated values with headers</li>
                <li><strong>JSON (.json)</strong> - Array of objects</li>
                <li><strong>YAML (.yaml, .yml)</strong> - Array of objects in YAML format</li>
              </ul>
              <p className="text-muted-foreground mt-2">
                Maximum file size: 10MB
              </p>
            </AccordionContent>
          </AccordionItem>

          <AccordionItem value="fields">
            <AccordionTrigger>Field Mapping</AccordionTrigger>
            <AccordionContent className="space-y-3 text-sm">
              {entityType === 'labels' && <LabelsFieldMapping />}
              {entityType === 'providers' && <ProvidersFieldMapping />}
              {entityType === 'subscriptions' && <SubscriptionsFieldMapping />}
            </AccordionContent>
          </AccordionItem>

          <AccordionItem value="validation">
            <AccordionTrigger>Validation Rules</AccordionTrigger>
            <AccordionContent className="space-y-3 text-sm">
              {entityType === 'labels' && <LabelsValidation />}
              {entityType === 'providers' && <ProvidersValidation />}
              {entityType === 'subscriptions' && <SubscriptionsValidation />}
            </AccordionContent>
          </AccordionItem>

          <AccordionItem value="errors">
            <AccordionTrigger>Common Errors</AccordionTrigger>
            <AccordionContent className="space-y-3 text-sm">
              <CommonErrors />
            </AccordionContent>
          </AccordionItem>
        </Accordion>
      </CardContent>
    </Card>
  );
}

function LabelsFieldMapping() {
  return (
    <div className="space-y-3">
      <div>
        <h4 className="font-semibold mb-1">Required Fields</h4>
        <ul className="list-disc list-inside space-y-1 ml-2">
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">name</code> - Label name (string)</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">color</code> - Hex color code (e.g., #FF5733 or #AAFF5733)</li>
        </ul>
      </div>
      <div>
        <h4 className="font-semibold mb-1">Optional Fields</h4>
        <ul className="list-disc list-inside space-y-1 ml-2">
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">ownerType</code> - Owner type: personal, family, or system (defaults to personal)</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">ownerFamilyId</code> - Family ID (required if ownerType is family)</li>
        </ul>
      </div>
    </div>
  );
}

function ProvidersFieldMapping() {
  return (
    <div className="space-y-3">
      <div>
        <h4 className="font-semibold mb-1">Required Fields</h4>
        <ul className="list-disc list-inside space-y-1 ml-2">
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">name</code> - Provider name (string)</li>
        </ul>
      </div>
      <div>
        <h4 className="font-semibold mb-1">Optional Fields</h4>
        <ul className="list-disc list-inside space-y-1 ml-2">
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">description</code> - Provider description (string)</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">url</code> - Provider website URL (valid URL)</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">iconUrl</code> - Provider icon URL (valid URL)</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">pricingPageUrl</code> - Pricing page URL (valid URL)</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">labels</code> - Array of label names or comma-separated string</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">ownerType</code> - Owner type: personal, family, or system</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">ownerFamilyId</code> - Family ID (required if ownerType is family)</li>
        </ul>
      </div>
    </div>
  );
}

function SubscriptionsFieldMapping() {
  return (
    <div className="space-y-3">
      <div>
        <h4 className="font-semibold mb-1">Required Fields</h4>
        <ul className="list-disc list-inside space-y-1 ml-2">
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">providerId</code> - Provider ID (UUID string)</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">startDate</code> - Start date (ISO 8601 format: YYYY-MM-DD)</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">recurrency</code> - Billing cycle: daily, weekly, monthly, quarterly, yearly, or custom</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">ownerType</code> - Owner type: personal, family, or system</li>
        </ul>
      </div>
      <div>
        <h4 className="font-semibold mb-1">Optional Fields</h4>
        <ul className="list-disc list-inside space-y-1 ml-2">
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">friendlyName</code> - Custom subscription name (string)</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">endDate</code> - End date (ISO 8601 format: YYYY-MM-DD)</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">customRecurrency</code> - Custom recurrency period (number, required if recurrency is custom)</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">customPriceAmount</code> - Price amount (number)</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">customPriceCurrency</code> - Currency code (3-letter ISO code, e.g., USD, EUR)</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">labels</code> - Array of label names or comma-separated string</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">ownerFamilyId</code> - Family ID (required if ownerType is family)</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">payerType</code> - Payer type: family or family_member</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">payerFamilyId</code> - Payer family ID</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">payerMemberId</code> - Payer member ID (required if payerType is family_member)</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">freeTrialStartDate</code> - Free trial start date (ISO 8601 format)</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">freeTrialEndDate</code> - Free trial end date (ISO 8601 format)</li>
          <li><code className="text-xs bg-muted px-1 py-0.5 rounded">familyUsers</code> - Array of user IDs or comma-separated string</li>
        </ul>
      </div>
    </div>
  );
}

function LabelsValidation() {
  return (
    <div className="space-y-2">
      <Alert variant="default">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Name:</strong> Cannot be empty
        </AlertDescription>
      </Alert>
      <Alert variant="default">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Color:</strong> Must be a valid hex color code (e.g., #FF5733 or #AAFF5733). The # prefix is optional in your file.
        </AlertDescription>
      </Alert>
      <Alert variant="default">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Owner Type:</strong> Must be one of: personal, family, or system
        </AlertDescription>
      </Alert>
      <Alert variant="default">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Family ID:</strong> Required when owner type is family
        </AlertDescription>
      </Alert>
    </div>
  );
}

function ProvidersValidation() {
  return (
    <div className="space-y-2">
      <Alert variant="default">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Name:</strong> Cannot be empty
        </AlertDescription>
      </Alert>
      <Alert variant="default">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>URLs:</strong> All URL fields (url, iconUrl, pricingPageUrl) must be valid URLs starting with http:// or https://
        </AlertDescription>
      </Alert>
      <Alert variant="default">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Owner Type:</strong> Must be one of: personal, family, or system
        </AlertDescription>
      </Alert>
      <Alert variant="default">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Family ID:</strong> Required when owner type is family
        </AlertDescription>
      </Alert>
    </div>
  );
}

function SubscriptionsValidation() {
  return (
    <div className="space-y-2">
      <Alert variant="default">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Provider ID:</strong> Cannot be empty. Must be a valid UUID of an existing provider.
        </AlertDescription>
      </Alert>
      <Alert variant="default">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Start Date:</strong> Required and must be a valid date in ISO 8601 format (YYYY-MM-DD)
        </AlertDescription>
      </Alert>
      <Alert variant="default">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>End Date:</strong> Must be a valid date and after the start date
        </AlertDescription>
      </Alert>
      <Alert variant="default">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Recurrency:</strong> Must be one of: daily, weekly, monthly, quarterly, yearly, or custom
        </AlertDescription>
      </Alert>
      <Alert variant="default">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Custom Recurrency:</strong> Required when recurrency is custom. Must be a positive number.
        </AlertDescription>
      </Alert>
      <Alert variant="default">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Custom Price:</strong> Both amount (non-negative number) and currency (3-letter ISO code) are required if price is specified
        </AlertDescription>
      </Alert>
      <Alert variant="default">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Owner Type:</strong> Must be one of: personal, family, or system
        </AlertDescription>
      </Alert>
      <Alert variant="default">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Payer Type:</strong> Must be one of: family or family_member
        </AlertDescription>
      </Alert>
      <Alert variant="default">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Free Trial:</strong> Both start and end dates are required. End date must be after start date.
        </AlertDescription>
      </Alert>
    </div>
  );
}

function CommonErrors() {
  return (
    <div className="space-y-2">
      <Alert variant="destructive">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>File format not supported:</strong> Make sure your file has a .csv, .json, .yaml, or .yml extension
        </AlertDescription>
      </Alert>
      <Alert variant="destructive">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Parse error:</strong> Check that your file is properly formatted. For CSV, ensure headers are present. For JSON/YAML, ensure valid syntax.
        </AlertDescription>
      </Alert>
      <Alert variant="destructive">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Required field missing:</strong> Check that all required fields are present in your file and not empty
        </AlertDescription>
      </Alert>
      <Alert variant="destructive">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>Invalid format:</strong> Ensure dates are in YYYY-MM-DD format, colors are hex codes, and URLs start with http:// or https://
        </AlertDescription>
      </Alert>
      <Alert variant="destructive">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>
          <strong>API error:</strong> The record may already exist, or you may have reached your quota limit. Check the error message for details.
        </AlertDescription>
      </Alert>
    </div>
  );
}

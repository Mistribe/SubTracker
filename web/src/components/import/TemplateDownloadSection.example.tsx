import { TemplateDownloadSection } from "./TemplateDownloadSection";

export default function TemplateDownloadSectionExample() {
  return (
    <div className="space-y-6 p-6">
      <div>
        <h2 className="text-2xl font-bold mb-4">Labels Template Downloads</h2>
        <TemplateDownloadSection entityType="labels" />
      </div>

      <div>
        <h2 className="text-2xl font-bold mb-4">Providers Template Downloads</h2>
        <TemplateDownloadSection entityType="providers" />
      </div>

      <div>
        <h2 className="text-2xl font-bold mb-4">Subscriptions Template Downloads</h2>
        <TemplateDownloadSection entityType="subscriptions" />
      </div>
    </div>
  );
}

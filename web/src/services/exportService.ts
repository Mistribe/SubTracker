import { useApiClient } from "@/hooks/use-api-client";

export type ExportFormat = "csv" | "json" | "yaml";
export type EntityType = "labels" | "subscriptions" | "providers";

/**
 * Downloads exported data as a file
 */
function downloadFile(data: Blob, filename: string) {
    const url = window.URL.createObjectURL(data);
    const link = document.createElement("a");
    link.href = url;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
}

/**
 * Hook for exporting entities
 */
export function useExportService() {
    const {apiClient} = useApiClient();

    const exportLabels = async (format: ExportFormat = "csv") => {
        const response = await apiClient.labels.labelsExportGetRaw({ format });
        const blob = await response.raw.blob();
        const filename = `labels-export.${format}`;
        downloadFile(blob, filename);
    };

    const exportSubscriptions = async (format: ExportFormat = "csv") => {
        const response = await apiClient.subscriptions.subscriptionsExportGetRaw({ format });
        const blob = await response.raw.blob();
        const filename = `subscriptions-export.${format}`;
        downloadFile(blob, filename);
    };

    const exportProviders = async (format: ExportFormat = "csv") => {
        const response = await apiClient.providers.providersExportGetRaw({ format });
        const blob = await response.raw.blob();
        const filename = `providers-export.${format}`;
        downloadFile(blob, filename);
    };

    return {
        exportLabels,
        exportSubscriptions,
        exportProviders,
    };
}

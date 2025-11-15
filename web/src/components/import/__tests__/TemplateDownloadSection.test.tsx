import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { TemplateDownloadSection } from "../TemplateDownloadSection";

describe("TemplateDownloadSection", () => {
  describe("Labels entity type", () => {
    it("renders all three format links for labels", () => {
      render(<TemplateDownloadSection entityType="labels" />);

      expect(screen.getByText("CSV Template")).toBeInTheDocument();
      expect(screen.getByText("JSON Template")).toBeInTheDocument();
      expect(screen.getByText("YAML Template")).toBeInTheDocument();
    });

    it("has correct file paths for labels templates", () => {
      render(<TemplateDownloadSection entityType="labels" />);

      const csvLink = screen.getByLabelText("Download CSV Template for labels");
      const jsonLink = screen.getByLabelText("Download JSON Template for labels");
      const yamlLink = screen.getByLabelText("Download YAML Template for labels");

      expect(csvLink).toHaveAttribute("href", "/templates/labels/labels-template.csv");
      expect(jsonLink).toHaveAttribute("href", "/templates/labels/labels-template.json");
      expect(yamlLink).toHaveAttribute("href", "/templates/labels/labels-template.yaml");
    });

    it("has download attribute with correct filename for labels", () => {
      render(<TemplateDownloadSection entityType="labels" />);

      const csvLink = screen.getByLabelText("Download CSV Template for labels");
      const jsonLink = screen.getByLabelText("Download JSON Template for labels");
      const yamlLink = screen.getByLabelText("Download YAML Template for labels");

      expect(csvLink).toHaveAttribute("download", "labels-template.csv");
      expect(jsonLink).toHaveAttribute("download", "labels-template.json");
      expect(yamlLink).toHaveAttribute("download", "labels-template.yaml");
    });
  });

  describe("Providers entity type", () => {
    it("renders all three format links for providers", () => {
      render(<TemplateDownloadSection entityType="providers" />);

      expect(screen.getByText("CSV Template")).toBeInTheDocument();
      expect(screen.getByText("JSON Template")).toBeInTheDocument();
      expect(screen.getByText("YAML Template")).toBeInTheDocument();
    });

    it("has correct file paths for providers templates", () => {
      render(<TemplateDownloadSection entityType="providers" />);

      const csvLink = screen.getByLabelText("Download CSV Template for providers");
      const jsonLink = screen.getByLabelText("Download JSON Template for providers");
      const yamlLink = screen.getByLabelText("Download YAML Template for providers");

      expect(csvLink).toHaveAttribute("href", "/templates/providers/providers-template.csv");
      expect(jsonLink).toHaveAttribute("href", "/templates/providers/providers-template.json");
      expect(yamlLink).toHaveAttribute("href", "/templates/providers/providers-template.yaml");
    });

    it("has download attribute with correct filename for providers", () => {
      render(<TemplateDownloadSection entityType="providers" />);

      const csvLink = screen.getByLabelText("Download CSV Template for providers");
      const jsonLink = screen.getByLabelText("Download JSON Template for providers");
      const yamlLink = screen.getByLabelText("Download YAML Template for providers");

      expect(csvLink).toHaveAttribute("download", "providers-template.csv");
      expect(jsonLink).toHaveAttribute("download", "providers-template.json");
      expect(yamlLink).toHaveAttribute("download", "providers-template.yaml");
    });
  });

  describe("Subscriptions entity type", () => {
    it("renders all three format links for subscriptions", () => {
      render(<TemplateDownloadSection entityType="subscriptions" />);

      expect(screen.getByText("CSV Template")).toBeInTheDocument();
      expect(screen.getByText("JSON Template")).toBeInTheDocument();
      expect(screen.getByText("YAML Template")).toBeInTheDocument();
    });

    it("has correct file paths for subscriptions templates", () => {
      render(<TemplateDownloadSection entityType="subscriptions" />);

      const csvLink = screen.getByLabelText("Download CSV Template for subscriptions");
      const jsonLink = screen.getByLabelText("Download JSON Template for subscriptions");
      const yamlLink = screen.getByLabelText("Download YAML Template for subscriptions");

      expect(csvLink).toHaveAttribute("href", "/templates/subscriptions/subscriptions-template.csv");
      expect(jsonLink).toHaveAttribute("href", "/templates/subscriptions/subscriptions-template.json");
      expect(yamlLink).toHaveAttribute("href", "/templates/subscriptions/subscriptions-template.yaml");
    });

    it("has download attribute with correct filename for subscriptions", () => {
      render(<TemplateDownloadSection entityType="subscriptions" />);

      const csvLink = screen.getByLabelText("Download CSV Template for subscriptions");
      const jsonLink = screen.getByLabelText("Download JSON Template for subscriptions");
      const yamlLink = screen.getByLabelText("Download YAML Template for subscriptions");

      expect(csvLink).toHaveAttribute("download", "subscriptions-template.csv");
      expect(jsonLink).toHaveAttribute("download", "subscriptions-template.json");
      expect(yamlLink).toHaveAttribute("download", "subscriptions-template.yaml");
    });
  });

  describe("Accessibility", () => {
    it("has proper ARIA labels for all download links", () => {
      render(<TemplateDownloadSection entityType="labels" />);

      expect(screen.getByLabelText("Download CSV Template for labels")).toBeInTheDocument();
      expect(screen.getByLabelText("Download JSON Template for labels")).toBeInTheDocument();
      expect(screen.getByLabelText("Download YAML Template for labels")).toBeInTheDocument();
    });

    it("has aria-hidden on decorative icons", () => {
      const { container } = render(<TemplateDownloadSection entityType="labels" />);

      // Check that SVG icons have aria-hidden attribute
      const icons = container.querySelectorAll('svg[aria-hidden="true"]');
      expect(icons.length).toBeGreaterThan(0);
    });

    it("renders descriptive text for each format", () => {
      render(<TemplateDownloadSection entityType="labels" />);

      expect(screen.getByText("Comma-separated values format")).toBeInTheDocument();
      expect(screen.getByText("JavaScript Object Notation format")).toBeInTheDocument();
      expect(screen.getByText("YAML Ain't Markup Language format")).toBeInTheDocument();
    });

    it("renders section title and description", () => {
      render(<TemplateDownloadSection entityType="labels" />);

      expect(screen.getByText("Download Template Files")).toBeInTheDocument();
      expect(screen.getByText(/choose a template format to get started/i)).toBeInTheDocument();
    });
  });

  describe("Responsive layout", () => {
    it("renders with grid layout classes", () => {
      const { container } = render(<TemplateDownloadSection entityType="labels" />);

      const gridContainer = container.querySelector('.grid.grid-cols-1');
      expect(gridContainer).toBeInTheDocument();
      expect(gridContainer).toHaveClass('grid-cols-1', 'sm:grid-cols-3', 'gap-3');
    });
  });
});

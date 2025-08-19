"use client"

import { useMemo, useState } from "react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { Check, ChevronsUpDown, X } from "lucide-react";
import { cn } from "@/lib/utils";

export type MultiSelectOption = {
  value: string;
  label: string;
};

export interface MultiSelectProps {
  options: MultiSelectOption[];
  selectedValues: string[];
  onToggle?: (value: string) => void;
  onChange?: (values: string[]) => void;
  placeholder?: string;
  searchPlaceholder?: string;
  emptyMessage?: string;
  buttonClassName?: string;
  className?: string;
}

/**
 * Generic multi-select component using shadcn/ui primitives + Tailwind.
 * - Shows selected values as removable badges in the trigger button.
 * - Uses Popover + Command for searchable options list with check indicators.
 */
export function MultiSelect({
  options,
  selectedValues,
  onToggle,
  onChange,
  placeholder = "Select options",
  searchPlaceholder = "Search...",
  emptyMessage = "No results.",
  buttonClassName,
  className,
}: MultiSelectProps) {
  const [open, setOpen] = useState(false);

  const selectedMap = useMemo(() => new Set(selectedValues), [selectedValues]);

  const handleToggle = (v: string) => {
    if (onToggle) {
      onToggle(v);
      return;
    }
    if (onChange) {
      const next = new Set(selectedValues);
      if (next.has(v)) next.delete(v);
      else next.add(v);
      onChange(Array.from(next));
    }
  };

  return (
    <div className={cn("w-full", className)}>
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            role="combobox"
            aria-expanded={open}
            className={cn("w-full justify-between flex-wrap gap-2", buttonClassName)}
          >
            {selectedValues.length > 0 ? (
              <div className="flex flex-wrap gap-1">
                {selectedValues
                  .map((v) => options.find((o) => o.value === v))
                  .filter(Boolean)
                  .map((opt) => (
                    <Badge key={opt!.value} variant="outline" className="flex items-center gap-1">
                      {opt!.label}
                      <X
                        className="h-3 w-3 cursor-pointer"
                        onMouseDown={(e) => {
                          e.preventDefault();
                          e.stopPropagation();
                        }}
                        onClick={(e) => {
                          e.preventDefault();
                          e.stopPropagation();
                          handleToggle(opt!.value);
                        }}
                      />
                    </Badge>
                  ))}
              </div>
            ) : (
              <span className="text-muted-foreground">{placeholder}</span>
            )}
            <ChevronsUpDown className="ml-auto h-4 w-4 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-72 p-0">
          <Command>
            <CommandInput placeholder={searchPlaceholder} />
            <CommandList>
              <CommandEmpty>{emptyMessage}</CommandEmpty>
              <CommandGroup>
                {options.map((opt) => (
                  <CommandItem
                    key={opt.value}
                    value={opt.label}
                    onSelect={() => {
                      handleToggle(opt.value);
                    }}
                  >
                    <Check
                      className={cn(
                        "mr-2 h-4 w-4",
                        selectedMap.has(opt.value) ? "opacity-100" : "opacity-0"
                      )}
                    />
                    {opt.label}
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
    </div>
  );
}

export default MultiSelect;

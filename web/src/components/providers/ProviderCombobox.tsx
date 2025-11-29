import { useEffect, useMemo, useRef, useState } from "react";
import { Check, ChevronsUpDown } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import Provider from "@/models/provider";
import { useAllProvidersQuery } from "@/hooks/providers/useAllProvidersQuery";
import { useProviderQuery } from "@/hooks/providers/useProviderQuery";

interface ProviderComboboxProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  emptyMessage?: string;
}

export function ProviderCombobox({
  value,
  onChange,
  placeholder = "Select a provider",
  emptyMessage = "No provider found.",
}: ProviderComboboxProps) {
  const [open, setOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const triggerRef = useRef<HTMLButtonElement>(null);
  const listRef = useRef<HTMLDivElement>(null);

  // Fetch providers from backend with search and infinite scroll support
  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
    isError,
  } = useAllProvidersQuery({ search: searchQuery });

  const fetchedProviders: Provider[] = useMemo(
    () => data?.pages.flatMap((p) => p.providers) ?? [],
    [data]
  );

  // Resolve selected provider: from current pages or fetch by id if missing
  const providerFromPages = useMemo(
    () => fetchedProviders.find((p) => p.id === value),
    [fetchedProviders, value]
  );
  const { data: providerById } = useProviderQuery(providerFromPages ? undefined : value);
  const selectedProvider = providerFromPages ?? providerById;

  // Close the popover when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        triggerRef.current &&
        !triggerRef.current.contains(event.target as Node)
      ) {
        setOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  // Infinite scroll handler for the list
  const handleScroll = (e: React.UIEvent<HTMLDivElement>) => {
    const el = e.currentTarget;
    const nearBottom = el.scrollTop + el.clientHeight >= el.scrollHeight - 16;
    if (nearBottom && hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  };

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          ref={triggerRef}
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className="w-full justify-between"
        >
          {selectedProvider ? selectedProvider.name : placeholder}
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-full p-0">
        <Command shouldFilter={false}>
          <CommandInput
            placeholder="Search provider..."
            value={searchQuery}
            onValueChange={setSearchQuery}
          />
          <CommandList ref={listRef} onScroll={handleScroll}>
            <CommandEmpty>
              {isLoading ? "Loading providers..." : isError ? "Failed to load providers." : emptyMessage}
            </CommandEmpty>
            <CommandGroup>
              {fetchedProviders.map((provider) => (
                <CommandItem
                  key={provider.id}
                  value={provider.id}
                  onSelect={(currentValue) => {
                    onChange(currentValue);
                    setOpen(false);
                    setSearchQuery("");
                  }}
                >
                  <Check
                    className={cn(
                      "mr-2 h-4 w-4",
                      value === provider.id ? "opacity-100" : "opacity-0"
                    )}
                  />
                  {provider.name}
                </CommandItem>
              ))}
            </CommandGroup>
            {hasNextPage && (
              <div className="py-2 text-center text-muted-foreground text-sm">
                {isFetchingNextPage ? "Loading more..." : "Scroll to load more"}
              </div>
            )}
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
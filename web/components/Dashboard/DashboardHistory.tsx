import useScanHistory from "@/hooks/dashboard/useScanHistory";
import { useEffect } from "react";
import DashboardEmptyState from "./DashboardEmptyState";
import { useUser } from "@clerk/nextjs";
import DashboardHistoryCard from "./DashboardHistoryCard";
import {
  Select,
  Button,
  ListBox,
  ListBoxItem,
  Popover,
  SelectValue,
} from "react-aria-components";
import { ChevronDownIcon } from "lucide-react";

const DashboardHistory = () => {
  const { scans, getScanHistory } = useScanHistory();
  const { isLoaded } = useUser();

  useEffect(() => {
    if (isLoaded) {
      getScanHistory();
    }
  }, [getScanHistory, isLoaded]);

  if (scans.length === 0) {
    return <DashboardEmptyState />;
  }

  const filteredScans = scans.filter((scan) => scan.metadata.title !== "");

  return (
    <section className="flex flex-col gap-8 p-8">
      <div className="flex items-center justify-between">
        <div className="space-y-1">
          <h2 className="text-2xl font-medium text-gray-800">
            History of scans
          </h2>
          <p className="text-sm text-gray-500">
            View and manage your previous website scans
          </p>
        </div>
        <div className="flex gap-4">
          <Select defaultSelectedKey="all">
            <Button className="flex items-center justify-between w-[160px] px-4 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 bg-white hover:border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
              <SelectValue />
              <ChevronDownIcon className="w-4 h-4 ml-2" />
            </Button>
            <Popover>
              <ListBox className="w-[160px] bg-white border border-gray-200 rounded-lg shadow-lg mt-1 py-1">
                <ListBoxItem
                  id="all"
                  className="px-4 py-2 outline-none cursor-default hover:bg-gray-50 focus:bg-gray-50"
                >
                  All statuses
                </ListBoxItem>
                <ListBoxItem
                  id="success"
                  className="px-4 py-2 outline-none cursor-default hover:bg-gray-50 focus:bg-gray-50"
                >
                  Success
                </ListBoxItem>
                <ListBoxItem
                  id="needs-improvement"
                  className="px-4 py-2 outline-none cursor-default hover:bg-gray-50 focus:bg-gray-50"
                >
                  Needs improvement
                </ListBoxItem>
                <ListBoxItem
                  id="failed"
                  className="px-4 py-2 outline-none cursor-default hover:bg-gray-50 focus:bg-gray-50"
                >
                  Failed
                </ListBoxItem>
              </ListBox>
            </Popover>
          </Select>
          <Select defaultSelectedKey="recent">
            <Button className="flex items-center justify-between w-[160px] px-4 py-2 border border-gray-200 rounded-lg text-sm text-gray-600 bg-white hover:border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
              <SelectValue />
              <ChevronDownIcon className="w-4 h-4 ml-2" />
            </Button>
            <Popover>
              <ListBox className="w-[160px] bg-white border border-gray-200 rounded-lg shadow-lg mt-1 py-1">
                <ListBoxItem
                  id="recent"
                  className="px-4 py-2 outline-none cursor-default hover:bg-gray-50 focus:bg-gray-50"
                >
                  Most recent
                </ListBoxItem>
                <ListBoxItem
                  id="oldest"
                  className="px-4 py-2 outline-none cursor-default hover:bg-gray-50 focus:bg-gray-50"
                >
                  Oldest
                </ListBoxItem>
              </ListBox>
            </Popover>
          </Select>
        </div>
      </div>
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {filteredScans.map((scan) => (
          <DashboardHistoryCard
            key={scan.id}
            metadata={scan.metadata}
            scan_id={scan.id}
          />
        ))}
      </div>
    </section>
  );
};

export default DashboardHistory;

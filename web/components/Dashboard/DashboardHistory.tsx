import useScanHistory from "@/hooks/dashboard/useScanHistory";
import { useEffect } from "react";
import DashboardEmptyState from "./DashboardEmptyState";
import { useUser } from "@clerk/nextjs";
import DashboardHistoryCard from "./DashboardHistoryCard";

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

  return (
    <div className="flex flex-col gap-4">
      {scans &&
        scans.map((scan) => (
          <DashboardHistoryCard key={scan.id} metadata={scan.metadata} />
        ))}
    </div>
  );
};

export default DashboardHistory;

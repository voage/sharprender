import { fetcher } from "@/lib/fetcher";
import { Scan } from "@/types/scan";
import { useUser } from "@clerk/nextjs";
import { useState, useCallback } from "react";

const useScanHistory = () => {
  const [scans, setScans] = useState<Scan[]>([]);
  const { user } = useUser();

  const getScanHistory = useCallback(async () => {
    const response = await fetcher<{ scans: Scan[] }>(
      `/scan/history?user_id=${user?.id}`
    );
    setScans(response.scans);
  }, [user?.id]);

  return { scans, getScanHistory };
};

export default useScanHistory;

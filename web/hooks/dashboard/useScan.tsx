import { useState } from "react";
import { ScanResult } from "@/types/scan";
import { fetcher } from "@/lib/fetcher";
import { useUser } from "@clerk/nextjs";

const useScan = () => {
  const [scan, setScan] = useState<ScanResult | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { user } = useUser();

  const createScan = async (url: string): Promise<string> => {
    if (!user || !user.id) {
      throw new Error("User is not logged in or user ID is missing.");
    }

    setIsLoading(true);
    try {
      const response = await fetcher<{ scan_id: string }>(`/scan`, {
        method: "POST",
        data: {
          user_id: user.id.toString(),
          url,
        },
      });
      return response.scan_id;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : String(err);
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  const getScanResults = async (scanId: string): Promise<void> => {
    setIsLoading(true);
    try {
      const response = await fetcher<ScanResult>(`/scan/${scanId}`);
      setScan(response);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : String(err);
      setError(errorMessage);
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  return { createScan, getScanResults, scan, isLoading, error };
};

export default useScan;

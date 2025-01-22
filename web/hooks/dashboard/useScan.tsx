import { useState } from "react";
import { Scan } from "@/types/scan";
import { fetcher } from "@/lib/fetcher";

const useScan = () => {
  const [scan, setScan] = useState<Scan | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const createScan = async (url: string): Promise<string> => {
    setIsLoading(true);
    try {
      const response = await fetcher<{ scan_id: string }>(`/scan`, {
        method: "POST",
        data: { url },
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
      const response = await fetcher<Scan>(`/scan/${scanId}`);
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

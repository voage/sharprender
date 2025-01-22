import { DashboardLayout } from "@/components/DashboardLayout";
import { isValidURL, formatURL } from "@/lib/urlUtils";
import { useState } from "react";
import { TextField, Input, Form, Button, Text } from "react-aria-components";
import useScan from "@/hooks/dashboard/useScan";
import toast from "react-hot-toast";
import DashboardDataGrid from "@/components/Dashboard/DashboardDataGrid";
import DashboardLoader from "@/components/Dashboard/DashboardLoader";
import DashboardEmptyState from "@/components/Dashboard/DashboardEmptyState";

export default function Home() {
  const [url, setUrl] = useState("");
  const { scan, createScan, getScanResults, isLoading, error } = useScan();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const formattedUrl = formatURL(url);

    if (!isValidURL(formattedUrl)) {
      toast.error("Please enter a valid URL");
      return;
    }

    try {
      const scanId = await createScan(formattedUrl);
      setUrl("");
      toast.success("Scan started");

      await getScanResults(scanId);
    } catch {
      toast.error("An error occurred. Please try again.");
    }
  };

  return (
    <DashboardLayout className="flex flex-col gap-10">
      <section>
        <Text className="text-gray-700 font-medium" slot="description">
          Enter the URL of the page you want to scan
        </Text>
        <Form
          onSubmit={handleSubmit}
          className="w-full flex flex-row gap-4 mt-2"
        >
          <TextField className="w-full flex-1">
            <Input
              className="w-full rounded-md px-4 py-2 border border-gray-200/50 shadow-sm 
            focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
              placeholder="https://example.com"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
            />
          </TextField>

          <Button
            type="submit"
            className="w-full max-w-fit bg-primary text-white px-8 py-2 rounded-md shadow-sm"
          >
            Scan
          </Button>
        </Form>

        {error && <p className="text-red-500 mt-2">{error}</p>}
      </section>

      {isLoading && <DashboardLoader />}

      {!scan && !isLoading && <DashboardEmptyState />}

      {scan && <DashboardDataGrid data={scan} />}
    </DashboardLayout>
  );
}

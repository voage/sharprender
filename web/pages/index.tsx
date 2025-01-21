import DashboardOverviewCard from "@/components/Dashboard/DashboardOverviewCard";
import DashboardTableOverview from "@/components/Dashboard/DashboardTableOverview";
import { DashboardLayout } from "@/components/DashboardLayout";
import dynamic from "next/dynamic";
import { Bolt, Globe, Paperclip } from "lucide-react";
import { isValidURL, formatURL } from "@/lib/urlUtils";
import { submitScan } from "@/lib/api/dashboard";
import { useState } from "react";
import { TextField, Input, Form, Button, Text } from "react-aria-components";

const DashboardPieChart = dynamic(
  () => import("@/components/Dashboard/DashboardPieChart"),
  { ssr: false }
);
const DashboardScatterPlotChart = dynamic(
  () => import("@/components/Dashboard/DashboardScatterPlotChart"),
  { ssr: false }
);

export default function Home() {
  const [url, setUrl] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    const formattedUrl = formatURL(url);

    if (!isValidURL(formattedUrl)) {
      setError("Please enter a valid URL.");
      return;
    }

    try {
      await submitScan(formattedUrl);
      setUrl(""); // Clear the input on success
      setError(""); // Clear any previous error
    } catch (err: unknown) {
      console.error("Error submitting form:", err);
      setError("An error occurred. Please try again.");
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

      <section className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <DashboardOverviewCard
          metric="Total Pages"
          value="656"
          description="The total number of pages scanned"
          icon={<Globe className="w-4 h-4 text-gray-500" />}
        />

        <DashboardOverviewCard
          metric="Avg. Image Size"
          value="10.5 MB"
          description="The average size of images on the page"
          icon={<Paperclip className="w-4 h-4 text-gray-500" />}
        />

        <DashboardOverviewCard
          metric="Total Load Time"
          value="15.5s"
          description="The total time it takes to load the page"
          icon={<Bolt className="w-4 h-4 text-gray-500" />}
        />
      </section>

      <section className="flex flex-row gap-4">
        <DashboardPieChart />
        <DashboardScatterPlotChart />
      </section>

      <section>
        <DashboardTableOverview />
      </section>
    </DashboardLayout>
  );
}

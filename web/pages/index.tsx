import DashboardOverviewCard from "@/components/Dashboard/DashboardOverviewCard";
import DashboardTableOverview from "@/components/Dashboard/DashboardTableOverview";
import { DashboardLayout } from "@/components/DashboardLayout";
import dynamic from "next/dynamic";
import { Bolt, Globe, Paperclip } from "lucide-react";
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
  return (
    <DashboardLayout className="flex flex-col gap-10">
      <section>
        <Text className="text-gray-700 font-medium" slot="description">
          Enter the URL of the page you want to scan
        </Text>
        <Form onSubmit={() => {}} className="w-full flex flex-row gap-4 mt-2">
          <TextField className="w-full flex-1">
            <Input
              className="w-full rounded-md px-4 py-2 border border-gray-200/50 shadow-sm 
            focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
              placeholder="https://example.com"
            />
          </TextField>

          <Button
            type="submit"
            className="w-full max-w-fit bg-primary text-white px-8 py-2 rounded-md shadow-sm"
          >
            Scan
          </Button>
        </Form>
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

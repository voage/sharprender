import { Scan } from "@/types/scan";
import DashboardOverviewCard from "./DashboardOverviewCard";
import { Globe, Paperclip, Bolt } from "lucide-react";
import DashboardTableOverview from "./DashboardTableOverview";
import dynamic from "next/dynamic";
const DashboardPieChart = dynamic(
  () => import("@/components/Dashboard/DashboardPieChart"),
  { ssr: false }
);
const DashboardScatterPlotChart = dynamic(
  () => import("@/components/Dashboard/DashboardScatterPlotChart"),
  { ssr: false }
);

interface DashboardDataGridProps {
  data: Scan;
}

const DashboardDataGrid = ({ data }: DashboardDataGridProps) => {
  return (
    <>
      <section className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <DashboardOverviewCard
          metric="Total Images"
          value={`${data.aggregations.imageCount}`}
          description="The total number of images on the page"
          icon={<Bolt className="w-4 h-4 text-gray-500" />}
        />
        <DashboardOverviewCard
          metric="Avg. Load Time"
          value={`${data.aggregations.avgLoadTime}s`}
          description="The average time it takes to load the page"
          icon={<Globe className="w-4 h-4 text-gray-500" />}
        />

        <DashboardOverviewCard
          metric="Avg. Image Size"
          value={`${(data.aggregations.avgSize / 1024).toFixed(2)} KB`}
          description="The average size of images on the page"
          icon={<Paperclip className="w-4 h-4 text-gray-500" />}
        />
      </section>

      <section className="flex flex-row gap-4">
        <DashboardPieChart
          formatDistribution={data.aggregations.formatDistribution}
          totalImages={data.aggregations.imageCount}
        />
        <DashboardScatterPlotChart />
      </section>

      <section>
        <DashboardTableOverview images={data.images} />
      </section>
    </>
  );
};

export default DashboardDataGrid;

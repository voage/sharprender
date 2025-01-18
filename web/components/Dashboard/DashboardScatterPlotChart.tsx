import React from "react";
import {
  ScatterChart,
  Scatter,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

const mockImageData = [
  {
    id: 1,
    fileName: "discord_icon.png",
    fileSize: 107,
    loadTime: 52,
    format: "PNG",
  },
  {
    id: 2,
    fileName: "companyBanner.png",
    fileSize: 40,
    loadTime: 130,
    format: "PNG",
  },
  {
    id: 3,
    fileName: "Google_logo.png",
    fileSize: 83,
    loadTime: 168,
    format: "PNG",
  },
  {
    id: 4,
    fileName: "example_image.webp",
    fileSize: 15,
    loadTime: 20,
    format: "WebP",
  },
  { id: 5, fileName: "logo.svg", fileSize: 10, loadTime: 15, format: "SVG" },
];

const DashboardScatterPlotChart = () => {
  const CustomTooltip = ({
    active,
    payload,
  }: {
    active?: boolean;
    payload?: Array<{ name: string; value: number }>;
  }) => {
    if (active && payload && payload.length) {
      const data = payload[0];
      return (
        <div className="bg-white p-2 border border-gray-200 rounded-md shadow-sm">
          <p className="text-sm font-medium mb-1">{data.name}</p>
          <p className="text-sm text-gray-600">{`File Size: ${data.value} KB`}</p>
          <p className="text-sm text-gray-600">{`Load Time: ${data.value} ms`}</p>
        </div>
      );
    }
    return null;
  };

  return (
    <div className="rounded-md shadow-sm shadow-gray-100 border border-gray-100 px-6 py-4 flex-1 w-full">
      <h3 className="text-lg text-gray-700 font-semibold mb-4">
        Load Time vs File Size
      </h3>
      <ResponsiveContainer width="100%" height={350}>
        <ScatterChart
          width={500}
          height={300}
          margin={{ top: 20, right: 20, bottom: 20, left: 20 }}
        >
          <CartesianGrid strokeDasharray="3 3" stroke="#E2E8F0" />
          <XAxis
            type="number"
            dataKey="fileSize"
            name="File Size (KB)"
            unit=" KB"
            tick={{ fill: "#4A5568", fontSize: 12 }}
            stroke="#CBD5E0"
          />
          <YAxis
            type="number"
            dataKey="loadTime"
            name="Load Time (ms)"
            unit=" ms"
            tick={{ fill: "#4A5568", fontSize: 12 }}
            stroke="#CBD5E0"
          />
          <Tooltip content={<CustomTooltip />} />
          <Scatter
            name="Images"
            data={mockImageData}
            fill="#3182CE"
            stroke="#2C5282"
            strokeWidth={1}
            r={6}
          />
        </ScatterChart>
      </ResponsiveContainer>
    </div>
  );
};

export default DashboardScatterPlotChart;

import { ImageScanResult } from "@/types/scan";
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

// Add this interface for our transformed data
interface ScatterDataPoint {
  id: string;
  fileName: string;
  fileSize: number;
  loadTime: number;
  format: string;
}

const DashboardScatterPlotChart = ({
  images,
}: {
  images: ImageScanResult[];
}) => {
  // Transform ImageScanResult[] into scatter plot data
  const scatterData: ScatterDataPoint[] = images.map((image) => ({
    id: image.network.request_id,
    fileName: image.src.split("/").pop() || "unknown",
    // Convert bytes to KB and ensure it's a number
    fileSize: Number((image.size / 1024).toFixed(2)),
    // Use network.load_time and convert to ms
    loadTime: Number((image.network.load_time * 1000).toFixed(2)),
    format: image.format || "unknown",
  }));

  // Add console.log to debug the data
  console.log("Scatter Data:", scatterData);

  const CustomTooltip = ({
    active,
    payload,
  }: {
    active?: boolean;
    payload?: Array<any>;
  }) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div className="bg-white p-2 border border-gray-200 rounded-md shadow-sm">
          <p className="text-sm font-medium mb-1">{data.fileName}</p>
          <p className="text-sm text-gray-600">
            File Size: {data.fileSize.toFixed(2)} KB
          </p>
          <p className="text-sm text-gray-600">
            Load Time: {data.loadTime.toFixed(2)} ms
          </p>
          <p className="text-sm text-gray-600">Format: {data.format}</p>
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
            data={scatterData}
            fill="hsl(var(--primary))"
            stroke="hsl(var(--primary-foreground))"
            strokeWidth={1}
            r={6}
          />
        </ScatterChart>
      </ResponsiveContainer>
    </div>
  );
};

export default DashboardScatterPlotChart;

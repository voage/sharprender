import React from "react";
import {
  PieChart,
  Pie,
  Cell,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";

const formatData = [
  { format: "PNG", count: 3 },
  { format: "WebP", count: 1 },
  { format: "SVG", count: 1 },
];

const COLORS = ["#0EA5E9", "#0284C7", "#0369A1", "#0C4A6E"];

const DashboardPieChart = () => {
  const CustomTooltip = ({
    active,
    payload,
  }: {
    active?: boolean;
    payload?: Array<{ name: string; value: number }>;
  }) => {
    if (active && payload && payload.length) {
      return (
        <div className="bg-white p-2 border border-gray-200 rounded-md shadow-sm">
          <p className="text-sm font-medium">{`${payload[0].name}: ${payload[0].value}`}</p>
        </div>
      );
    }
    return null;
  };

  return (
    <div className="rounded-md shadow-sm shadow-gray-100 border border-gray-100 px-6 py-4 w-1/3">
      <h3 className="text-lg text-gray-700 font-semibold mb-4">
        Format Distribution
      </h3>
      <ResponsiveContainer width="100%" height={300}>
        <PieChart width={400} height={300}>
          <Pie
            data={formatData}
            dataKey="count"
            nameKey="format"
            cx="50%"
            cy="50%"
            outerRadius={100}
            innerRadius={60}
            fill="#8884d8"
            paddingAngle={5}
            label={({ name, percent }) =>
              `${name} (${(percent * 100).toFixed(0)}%)`
            }
            labelLine={false}
          >
            {formatData.map((entry, index) => (
              <Cell
                key={`cell-${index}`}
                fill={COLORS[index % COLORS.length]}
                stroke="white"
                strokeWidth={2}
              />
            ))}
          </Pie>
          <Tooltip content={<CustomTooltip />} />
          <Legend
            verticalAlign="bottom"
            height={36}
            iconType="circle"
            formatter={(value) => (
              <span className="text-sm text-gray-600">{value}</span>
            )}
          />
        </PieChart>
      </ResponsiveContainer>
    </div>
  );
};

export default DashboardPieChart;

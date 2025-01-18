import { Text } from "react-aria-components";

interface DashboardOverviewCardProps {
  metric: string;
  value: string;
  description: string;
  icon: React.ReactNode;
}

const DashboardOverviewCard = ({
  metric,
  value,
  description,
  icon,
}: DashboardOverviewCardProps) => {
  return (
    <div className="border border-gray-100 shadow-sm shadow-gray-100 rounded-lg px-4 py-2.5 flex flex-col gap-2">
      <div className="flex flex-col">
        <div className="flex flex-row justify-between items-center">
          <Text className="text-sm text-gray-500" slot="title">
            {metric}
          </Text>
          {icon}
        </div>
        <Text className="text-xl font-bold" slot="description">
          {value}
        </Text>
      </div>

      <Text className="text-sm text-gray-500" slot="description">
        {description}
      </Text>
    </div>
  );
};

export default DashboardOverviewCard;

import { WebsiteMetadata } from "@/types/scan";
import Link from "next/link";

interface DashboardHistoryCardProps {
  metadata: WebsiteMetadata;
  scan_id: string;
}

const DashboardHistoryCard = ({
  metadata,
  scan_id,
}: DashboardHistoryCardProps) => {
  const randomStatusArray = ["success", "needs_improvement", "failed"];
  const randomStatus =
    randomStatusArray[Math.floor(Math.random() * randomStatusArray.length)];

  return (
    <div className="group flex flex-col justify-between p-6 bg-white rounded-xl border border-gray-100 shadow-sm shadow-gray-100 hover:shadow-gray-200 hover:cursor-pointer transition-all duration-200">
      <div className="space-y-4">
        <div className="flex items-start justify-between gap-4">
          <h3 className="text-base font-medium text-gray-800 flex-1 line-clamp-1">
            {metadata.title}
          </h3>
          <span
            className={`shrink-0 px-3 py-1 text-xs font-medium rounded-full ${
              randomStatus === "success"
                ? "bg-green-50 text-green-700 border border-green-200"
                : randomStatus === "needs_improvement"
                ? "bg-yellow-50 text-yellow-700 border border-yellow-200"
                : "bg-red-50 text-red-700 border border-red-200"
            }`}
          >
            {randomStatus.replace("_", " ")}
          </span>
        </div>
        <p className="text-sm text-gray-600 line-clamp-3">
          {metadata.description}
        </p>
      </div>
      <div className="flex items-center justify-between mt-6">
        <span className="text-sm text-gray-500">Scanned 2 hours ago</span>
        <Link
          href={`/history/${scan_id}`}
          className="inline-flex items-center text-sm text-primary-500 hover:text-primary-700 font-medium group-hover:gap-1 transition-all"
        >
          View details
        </Link>
      </div>
    </div>
  );
};

export default DashboardHistoryCard;

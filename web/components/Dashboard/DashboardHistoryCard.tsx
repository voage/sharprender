import { WebsiteMetadata } from "@/types/scan";
import Image from "next/image";

interface DashboardHistoryCardProps {
  metadata: WebsiteMetadata;
}

const DashboardHistoryCard = ({ metadata }: DashboardHistoryCardProps) => {
  return (
    <div className="flex flex-col p-4 bg-white rounded-lg shadow-sm shadow-gray-200">
      <div className="flex flex-row gap-5 items-center justify-start">
        <Image
          src={
            metadata.og_image ||
            metadata.favicon ||
            "https://placehold.co/600x400?text=Hello+World"
          }
          alt="favicon"
          width={16}
          height={16}
          className="rounded-full"
        />
        <div className="flex flex-col gap-2">
          <h1 className="text-lg font-bold">{metadata.title}</h1>
          <p className="text-sm text-gray-500">{metadata.description}</p>
        </div>
      </div>
    </div>
  );
};

export default DashboardHistoryCard;

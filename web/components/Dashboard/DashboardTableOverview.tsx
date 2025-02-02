import { cn } from "@/lib/utils";
import Image from "next/image";
import { Table } from "react-aria-components";
import { useState } from "react";

import {
  Cell,
  Row,
  TableBody,
  TableHeader,
  DialogTrigger,
  Button,
} from "react-aria-components";
import { Column } from "react-aria-components";
import DashboardImageDetailModal from "./DashboardImageDetailModal";
import { ImageScanResult } from "@/types/scan";

interface DashboardTableOverviewProps {
  images: ImageScanResult[];
}

const TableColumn = ({
  children,
  className,
  isRowHeader = false,
}: {
  children: React.ReactNode;
  className?: string;
  isRowHeader?: boolean;
}) => {
  return (
    <Column
      className={cn(
        "py-3 px-4 font-medium text-sm text-gray-600 text-left",
        className
      )}
      isRowHeader={isRowHeader}
    >
      {children}
    </Column>
  );
};

const TableCell = ({
  children,
  className,
}: {
  children: React.ReactNode;
  className?: string;
}) => {
  return (
    <Cell className={cn("py-4 px-4 text-sm text-gray-700", className)}>
      {children}
    </Cell>
  );
};

const DashboardTableOverview = ({ images }: DashboardTableOverviewProps) => {
  const [isOpen, setOpen] = useState(false);
  const filteredImages = images.filter((image) => image.size > 0);
  return (
    <div className="w-full rounded-lg border border-gray-100 shadow-sm shadow-gray-100">
      <div className="px-6 py-4 border-b border-gray-100">
        <h2 className="text-lg font-semibold text-gray-900">Scan Results</h2>
      </div>
      <Table aria-label="Image Scan Results" className="w-full">
        <TableHeader className="bg-gray-50 border-b border-gray-100">
          <TableColumn isRowHeader>Thumbnail</TableColumn>
          <TableColumn isRowHeader>File Name</TableColumn>
          <TableColumn isRowHeader>Dimensions</TableColumn>
          <TableColumn isRowHeader>File Size</TableColumn>
          <TableColumn isRowHeader>Format</TableColumn>
          <TableColumn isRowHeader>Load Time</TableColumn>
          <TableColumn isRowHeader>Status</TableColumn>
          <TableColumn isRowHeader>Actions</TableColumn>
        </TableHeader>
        <TableBody>
          {filteredImages.map((image, index) => {
            const loadTimeStatus =
              image.network.load_time < 100
                ? "success"
                : image.network.load_time < 200
                ? "warning"
                : "error";

            return (
              <Row
                key={image.network.request_id || index}
                className="hover:bg-gray-50"
              >
                <TableCell>
                  <Image
                    src={image.src}
                    alt={image.alt || "Image thumbnail"}
                    width={40}
                    height={40}
                    className="rounded-lg object-cover"
                  />
                </TableCell>
                <TableCell>{image.alt || "No alt text"}</TableCell>
                <TableCell>{`${image.width}x${image.height} px`}</TableCell>
                <TableCell>{`${(image.size / 1024).toFixed(1)} KB`}</TableCell>
                <TableCell>{image.format.split("/")[1]}</TableCell>
                <TableCell>{`${(image.network.load_time * 1000).toFixed(
                  0
                )} ms`}</TableCell>
                <TableCell>
                  <span
                    className={cn(
                      "inline-flex px-2 py-1 rounded-full text-xs font-medium",
                      {
                        "bg-green-50 text-green-700":
                          loadTimeStatus === "success",
                        "bg-yellow-50 text-yellow-700":
                          loadTimeStatus === "warning",
                        "bg-red-50 text-red-700": loadTimeStatus === "error",
                      }
                    )}
                  >
                    {loadTimeStatus}
                  </span>
                </TableCell>
                <TableCell>
                  <DialogTrigger isOpen={isOpen} onOpenChange={setOpen}>
                    <Button className="text-sm text-gray-700 hover:bg-gray-400 py-2 px-4 rounded-lg">View</Button>
                    <DashboardImageDetailModal
                      image={image}
                      onClose={() => setOpen(false)}
                    />
                  </DialogTrigger>
                </TableCell>
              </Row>
            );
          })}
        </TableBody>
      </Table>
    </div>
  );
};

export default DashboardTableOverview;

import { cn } from "@/lib/utils";
import Image from "next/image";
import { Table } from "react-aria-components";

import { Cell, Row, TableBody, TableHeader } from "react-aria-components";
import { Column } from "react-aria-components";

const mockData = [
  {
    id: 1,
    thumbnail: "https://picsum.photos/seed/1/40",
    fileName: "discord_icon.png",
    dimensions: "16x16 px",
    fileSize: "107 KB",
    format: "PNG",
    loadTime: "52 ms",
    status: "success",
  },
  {
    id: 2,
    thumbnail: "https://picsum.photos/seed/2/40",
    fileName: "companyBanner.png",
    dimensions: "896x45 px",
    fileSize: "40 KB",
    format: "PNG",
    loadTime: "130 ms",
    status: "warning",
  },
  {
    id: 3,
    thumbnail: "https://picsum.photos/seed/3/40",
    fileName: "Google_logo.png",
    dimensions: "2008x2048 px",
    fileSize: "83 KB",
    format: "PNG",
    loadTime: "168 ms",
    status: "error",
  },
  {
    id: 4,
    thumbnail: "https://picsum.photos/seed/4/40",
    fileName: "Google_logo.png",
    dimensions: "2008x2048 px",
    fileSize: "83 KB",
    format: "PNG",
    loadTime: "168 ms",
    status: "warning",
  },
  {
    id: 5,
    thumbnail: "https://picsum.photos/seed/5/40",
    fileName: "Google_logo.png",
    dimensions: "2008x2048 px",
    fileSize: "83 KB",
    format: "PNG",
    loadTime: "168 ms",
    status: "success",
  },
  {
    id: 6,
    thumbnail: "https://picsum.photos/seed/6/40",
    fileName: "Google_logo.png",
    dimensions: "2008x2048 px",
    fileSize: "83 KB",
    format: "PNG",
    loadTime: "168 ms",
    status: "error",
  },
];

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

const DashboardTableOverview = () => {
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
        </TableHeader>
        <TableBody>
          {mockData.map((item) => (
            <Row key={item.id} className="hover:bg-gray-50">
              <TableCell>
                <Image
                  src={item.thumbnail}
                  alt={item.fileName}
                  width={40}
                  height={40}
                  className="rounded-lg object-cover"
                />
              </TableCell>
              <TableCell>{item.fileName}</TableCell>
              <TableCell>{item.dimensions}</TableCell>
              <TableCell>{item.fileSize}</TableCell>
              <TableCell>{item.format}</TableCell>
              <TableCell>{item.loadTime}</TableCell>
              <TableCell>
                <span
                  className={cn(
                    "inline-flex px-2 py-1 rounded-full text-xs font-medium",
                    {
                      "bg-green-50 text-green-700": item.status === "success",
                      "bg-yellow-50 text-yellow-700": item.status === "warning",
                      "bg-red-50 text-red-700": item.status === "error",
                    }
                  )}
                >
                  {item.status}
                </span>
              </TableCell>
            </Row>
          ))}
        </TableBody>
      </Table>
    </div>
  );
};

export default DashboardTableOverview;

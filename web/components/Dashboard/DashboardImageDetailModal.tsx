import { ImageScanResult } from "@/types/scan";
import { CheckCircle2Icon, NetworkIcon, SparkleIcon } from "lucide-react";
import Image from "next/image";
import { Button, Dialog, Modal } from "react-aria-components";

const DashboardImageDetailModal = ({
  image,
  onClose,
}: {
  image: ImageScanResult;
  onClose: () => void;
}) => {
  return (
    <Modal className="fixed inset-0 z-50 overflow-auto bg-black/25 flex min-h-full items-center justify-center p-4">
      <Dialog className="w-full max-w-4xl rounded-lg bg-white p-6 shadow-xl flex flex-col gap-4">
        <section className="grid grid-cols-12 gap-8 border-b border-gray-100 pb-8">
          <div className="relative col-span-4 flex flex-col gap-4">
            <div className="relative w-full max-w-xs h-full">
              <Image
                src={image.src}
                alt={image.alt}
                fill
                className="rounded-xl object-cover"
              />
            </div>
          </div>

          <div className="col-span-4">
            <div className="bg-gray-50 rounded-xl p-5 h-full border border-gray-100">
              <h3 className="text-gray-800 font-semibold text-lg mb-4 flex items-center gap-2">
                <svg
                  className="w-4 h-4 text-primary"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                  />
                </svg>
                Image Details
              </h3>
              <dl className="grid grid-cols-2 gap-4 items-center justify-between">
                {[
                  {
                    label: "Size",
                    value: `${(image.size / 1024).toFixed(1)} KB`,
                  },
                  { label: "Format", value: image.format },
                  {
                    label: "Dimensions",
                    value: `${image.width} × ${image.height}`,
                  },
                  {
                    label: "File Name",
                    value: image.src.split("/").pop() || image.alt,
                  },
                ].map((item) => (
                  <div
                    key={item.label}
                    className="group transition-all duration-200 rounded-lg"
                  >
                    <dt className="text-gray-500 text-xs mb-0.5">
                      {item.label}
                    </dt>
                    <dd className="text-gray-900 font-medium text-sm line-clamp-1">
                      {item.value}
                    </dd>
                  </div>
                ))}
              </dl>
            </div>
          </div>

          <div className="col-span-4">
            <div className="bg-gray-50 rounded-xl p-5 h-full border border-gray-100">
              <h3 className="text-gray-800 font-semibold text-lg mb-4 flex items-center gap-2">
                <NetworkIcon className="w-4 h-4 text-primary" />
                Network Data
              </h3>
              <dl className="grid grid-cols-2 gap-4 items-center justify-between">
                {[
                  {
                    label: "Status",
                    value: image.network.status,
                    badge: image.network.status === 200 ? "success" : undefined,
                  },
                  {
                    label: "Protocol",
                    value: image.network.protocol.toUpperCase(),
                  },
                  {
                    label: "Load Time",
                    value: `${image.network.load_time.toFixed(3)}s`,
                  },
                  {
                    label: "MIME Type",
                    value: image.network.mime_type,
                  },
                  {
                    label: "Cache Status",
                    value: "HIT", // TODO: get from response headers
                  },
                  {
                    label: "Content Length",
                    value: `${(image.size / 1024).toFixed(1)} KB`, // TODO: get from response headers
                  },
                ].map((item) => (
                  <div
                    key={item.label}
                    className="group transition-all duration-200 rounded-lg"
                  >
                    <dt className="text-gray-500 text-xs mb-0.5">
                      {item.label}
                    </dt>
                    <dd className="text-gray-900 font-medium text-sm flex items-center gap-2">
                      {item.value}
                      {item.badge === "success" && (
                        <span className="px-1.5 py-0.5 rounded-full bg-green-100 text-green-700 text-xs font-medium">
                          OK
                        </span>
                      )}
                    </dd>
                  </div>
                ))}
              </dl>
            </div>
          </div>
        </section>

        <section className="grid grid-cols-12 gap-8">
          <div className="col-span-12 p-4 rounded-lg flex flex-col gap-4">
            <h4 className="flex items-center gap-2 text-gray-800 font-semibold tracking-wide text-lg">
              <SparkleIcon color="purple" className="w-4 h-4" /> AI
              Recommendations
            </h4>
            <ul className="flex flex-col gap-2">
              {Object.entries(image.ai_recommendation).map(([key, value]) => (
                <li
                  key={key}
                  className="text-gray-700 flex text-sm tracking-wide items-center gap-2"
                >
                  <CheckCircle2Icon className="w-4 h-4" /> {value}
                </li>
              ))}
            </ul>
          </div>
        </section>

        <section className="flex justify-end">
          <div className="flex gap-4">
            <Button className="bg-primary px-4 py-2 rounded-lg text-white">
              Optimize
            </Button>

            <Button
              onPress={onClose}
              className="bg-gray-200 px-4 py-2 rounded-lg text-gray-800"
            >
              Close
            </Button>
          </div>
        </section>
      </Dialog>
    </Modal>
  );
};

export default DashboardImageDetailModal;

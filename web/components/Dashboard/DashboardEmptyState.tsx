import Image from "next/image";

const DashboardEmptyState = () => {
  return (
    <div className="flex flex-col items-center justify-center max-w-2xl mx-auto text-center gap-6">
      <Image
        src="/images/Waiting.png"
        alt="Scan placeholder"
        width={400}
        height={400}
        className="w-auto h-auto object-contain opacity-90"
      />
      <div className="space-y-3">
        <h2 className="text-2xl font-semibold text-gray-800">
          Ready to analyze your website?
        </h2>
        <p className="text-gray-600 max-w-md mx-auto">
          Enter a URL above to start scanning and get detailed insights about
          your page&apos;s image performance.
        </p>
      </div>
    </div>
  );
};

export default DashboardEmptyState;

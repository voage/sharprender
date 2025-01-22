import Lottie from "lottie-react";
import loadingAnimation from "@/public/animations/loading.json";
const DashboardLoader = () => {
  return (
    <div className="relative min-h-[600px] flex flex-col items-center justify-center">
      <div className="absolute inset-0 bg-gradient-to-b from-white via-white to-transparent z-10" />
      <div className="relative z-20 flex flex-col items-center justify-center max-w-2xl mx-auto text-center">
        <div className="w-full h-48">
          <Lottie
            animationData={loadingAnimation}
            loop={true}
            className="w-full h-full"
          />
        </div>
        <div className="space-y-3">
          <h2 className="text-2xl font-semibold bg-gradient-to-r from-primary/80 to-primary bg-clip-text text-transparent">
            Analyzing Your Website
          </h2>
          <div className="space-y-2">
            <p className="text-gray-600 text-lg">
              We&apos;re scanning your page for optimization opportunities
            </p>
            <p className="text-sm text-gray-500">
              This usually takes about 30 seconds
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DashboardLoader;

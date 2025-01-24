import Image from "next/image";

export default function DecorativeBlock() {
  return (
    <div className="rounded-lg shadow h-[95vh] w-[80vh] bg-primary-700 relative flex items-center justify-center">
      <div className="relative flex justify-center items-center w-[600px] h-[400px] mb-16">
        <div className="absolute top-0 left-12">
          <h1 className="text-2xl font-bold text-white">
            Optimize, Accelerate, Dominate
          </h1>
          <p className="text-sm text-gray-200">
            Take Control of Your Website Performance
          </p>
        </div>

        {/* Bottom Image */}
        <Image
          src="/images/speed.png"
          alt="Bottom Image"
          width={320}
          height={400}
          className="absolute z-40 rounded-lg shadow-xl mx-10 -right-4 -bottom-2 border-2 border-purple-200"
        />
        {/* Top Image */}
        <Image
          src="/images/dashboard.webp"
          alt="Top Image"
          width={460}
          height={400}
          className="absolute z-20 rounded-lg shadow-xl border-2 border-purple-200"
        />
      </div>
    </div>
  );
}

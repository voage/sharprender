import React from "react";
import DecorativeBlock from "./decorative_block";

interface AuthLayoutProps {
  title: string;
  description: string;
  children: React.ReactNode;
}
export default function AuthLayout({
  title,
  description,
  children,
}: AuthLayoutProps) {
  return (
    <div className="flex justify-center align-items h-screen mx-10">
      {/* Left Panel */}
      <div className="flex-1 flex flex-col justify-center px-8 bg-white">
        <div className="max-w-md mx-auto">
          <h1 className="text-3xl font-bold mb-4">{title}</h1>
          <p className="text-gray-600 mb-6">{description}</p>
          {children}
        </div>
      </div>
      {/* Right Panel */}
      <DecorativeBlock />
    </div>
  );
}

import { SignUp } from "@clerk/nextjs";
import AuthLayout from "@/components/Auth/auth_layout";
import { clerkAppearance } from "@/lib/clerk_appearance";

export default function Page() {
  return (
    <AuthLayout title="" description="">
      <div className="w-full max-w-md mx-auto  ">
        <SignUp appearance={clerkAppearance} />
      </div>
    </AuthLayout>
  );
}

import { Theme } from "@clerk/types";

export const clerkAppearance: Theme = {
  layout: {
    socialButtonsPlacement: "bottom",
    socialButtonsVariant: "blockButton",
    unsafe_disableDevelopmentModeWarnings: true,
    
  },
  elements: {
    formButtonPrimary: "primary-500 hover:primary-600 text-white font-semibold",
  },
};

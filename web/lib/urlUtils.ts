export const isValidURL = (url: string): boolean => {
  try {
    const parsedUrl = new URL(url);
    return parsedUrl.protocol === "http:" || parsedUrl.protocol === "https:";
  } catch (error) {
    return false; // If URL parsing fails, return false
    console.log(error);
  }
};

export const formatURL = (input: string): string => {
  if (input.startsWith("http://") || input.startsWith("https://")) {
    return input;
  }

  if (input.startsWith("www.")) {
    return `https://${input}`;
  }

  return `https://${input}`;
};

import axios, { AxiosRequestConfig } from "axios";

const axiosInstance = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080",
  headers: {
    "Content-Type": "application/json",
  },
});

export const fetcher = async <T>(
  url: string,
  config?: AxiosRequestConfig
): Promise<T> => {
  const response = await axiosInstance.get<T>(url, config);
  return response.data;
};

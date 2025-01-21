import axios from "axios";

export const submitScan = async (url: string): Promise<void> => {
  await axios.post(`http://localhost:8888/scan?url=${encodeURIComponent(url)}`);
};

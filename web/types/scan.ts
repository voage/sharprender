interface Aggregations {
  avgLoadTime: number;
  avgSize: number;
  imageCount: number;
  totalSize: number;
  formatDistribution: Record<string, number>;
}

export interface WebsiteMetadata {
  title: string;
  description: string;
  favicon: string;
  og_image: string;
  og_title: string;
  og_description: string;
  language: string;
}

interface AIRecommendation {
  format_recommendations: string;
  resize_recommendations: string;
  compression_recommendations: string;
  caching_recommendations: string;
  other_recommendations: string;
}

interface NetworkInfo {
  request_id: string;
  document_url: string;
  initiator_type: string;
  initiator_url: string;
  initiator_line_no: number;
  initiator_col_no: number;
  method: string;
  status: number;
  mime_type: string;
  protocol: string;
  remote_ip_address: string;
  remote_port: number;
  encoded_data_length: number;
  request_time: number;
  response_time: number;
  load_time: number;
  request_headers: Record<string, string>;
  response_headers: Record<string, string>;
}

interface TimingInfo {
  dns_lookup: number;
  connection_time: number;
  ssl_time: number;
  ttfb: number;
  content_download_time: number;
  transfer_size: number;
  encoded_body_size: number;
  decoded_body_size: number;
}

export interface ImageScanResult {
  src: string;
  alt: string;
  width: number;
  height: number;
  format: string;
  size: number;
  network: NetworkInfo;
  timing: TimingInfo;
  ai_recommendation: AIRecommendation;
}

export interface Scan {
  id: string;
  scan_id: string;
  user_id: string;
  url: string;
  metadata: WebsiteMetadata;
  images: ImageScanResult[];
  created_at: string;
}

export interface ScanResult {
  images: ImageScanResult[];
  aggregations: Aggregations;
}

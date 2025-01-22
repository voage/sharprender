interface Aggregations {
  avgLoadTime: number;
  avgSize: number;
  imageCount: number;
  totalSize: number;
  formatDistribution: Record<string, number>;
}

export interface ImageScanResult {
  src: string;
  alt: string;
  width: number;
  height: number;
  format: string;
  size: number;
  network: {
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
  };
  timing: {
    dns_lookup: number;
    connection_time: number;
    ssl_time: number;
    ttfb: number;
    content_download_time: number;
    transfer_size: number;
    encoded_body_size: number;
    decoded_body_size: number;
  };
  ai_recommendation: {
    format_recommendations: string;
    resize_recommendations: string;
    compression_recommendations: string;
    caching_recommendations: string;
    other_recommendations: string;
  };
}

export interface Scan {
  scan_id: string;
  url: string;
  created_at: string;
  images: ImageScanResult[];
  aggregations: Aggregations;
}

import DashboardDataGrid from "@/components/Dashboard/DashboardDataGrid";
import { DashboardLayout } from "@/components/DashboardLayout";
import { fetcher } from "@/lib/fetcher";
import { ScanResult } from "@/types/scan";
import { GetServerSidePropsContext } from "next";

export default function History({ scan }: { scan: ScanResult }) {
  return (
    <DashboardLayout>
      <DashboardDataGrid data={scan} />
    </DashboardLayout>
  );
}

export async function getServerSideProps(context: GetServerSidePropsContext) {
  const { id } = context.params as { id: string };
  if (!id) {
    return {
      notFound: true,
    };
  }

  try {
    const scan = await fetcher<ScanResult>(`/scan/${id}`);

    if (!scan) {
      return {
        notFound: true,
      };
    }

    return {
      props: { scan },
    };
  } catch (error) {
    console.error("Error fetching scan:", error);
    return {
      notFound: true,
    };
  }
}

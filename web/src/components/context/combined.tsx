import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

const qc = new QueryClient();

export const CombinedContext = ({ children }: React.PropsWithChildren) => {
  return <QueryClientProvider client={qc}>{children}</QueryClientProvider>;
};

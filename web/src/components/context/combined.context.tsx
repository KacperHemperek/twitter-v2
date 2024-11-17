import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { AuthContextProvider } from "./auth.context";

const qc = new QueryClient();

export const CombinedContext = ({ children }: React.PropsWithChildren) => {
  return (
    <QueryClientProvider client={qc}>
      <AuthContextProvider>{children}</AuthContextProvider>
    </QueryClientProvider>
  );
};

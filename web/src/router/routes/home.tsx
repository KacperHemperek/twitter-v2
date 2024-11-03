import { useQuery } from "@tanstack/react-query";

type HealthCheckResponse = {
  message: "string";
};

export const HomePage = () => {
  const { data, isLoading, isError } = useQuery({
    queryKey: ["test"],
    queryFn: async () => {
      const res = await fetch("http://localhost:1337/healthcheck", {
        headers: {
          "Content-Type": "application/json",
        },
      });

      const data = (await res.json()) as HealthCheckResponse;

      if (!res.ok) {
        throw new Error("Could not retrieve healthcheck data");
      }

      return data;
    },
  });

  if (!data && isLoading) {
    return <div>Loading healthcheck</div>;
  }

  if (!data && isError) {
    return <div>API is not responding</div>;
  }

  if (!data) {
    return <div>No data from endpoint</div>;
  }

  return <div>Response from healthcheck endpoint: {data.message}</div>;
};

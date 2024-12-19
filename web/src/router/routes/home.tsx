import { Link } from "@tanstack/react-router";
import { useAuth } from "../../components/context/auth.context";
import { TwButton } from "../../components/common/button";
import * as z from "zod";
import { useForm } from "react-hook-form";
import { useMutation } from "@tanstack/react-query";
import { ApiErrorResponse } from "../../types";

const createTweetSchema = z.object({
  body: z
    .string()
    .min(1, "Tweet must be at least 1 character long")
    .max(180, "Tweet must be at most 180 characters long"),
});

type CreateTweetFormValue = z.infer<typeof createTweetSchema>;

const CreateTweetForm = () => {
  const mutation = useMutation({
    mutationFn: async (data: CreateTweetFormValue) => {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/tweets`, {
        credentials: "include",
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      const body = await res.json();

      if (!res.ok) {
        const error = body as ApiErrorResponse;
        throw new Error(error.message);
      }
    },
    onSuccess: () => {
      form.reset();
    },
  });
  const form = useForm<CreateTweetFormValue>({
    defaultValues: {
      body: "",
    },
  });

  const sendTweet = (data: CreateTweetFormValue) => {
    mutation.mutate(data);
  };

  return (
    <form
      aria-disabled={mutation.isPending}
      className="flex flex-col"
      onSubmit={form.handleSubmit(sendTweet)}
    >
      <textarea
        {...form.register("body")}
        className="w-full p-2 bg-gray-900 border border-gray-800 rounded resize-none mb-4"
        placeholder="What's happening?"
      ></textarea>
      <TwButton
        aria-disabled={mutation.isPending}
        disabled={mutation.isPending}
        className="max-w-fit self-end"
        type="submit"
      >
        Tweet
      </TwButton>
    </form>
  );
};

export const HomePage = () => {
  const { user, logout, loadingUser } = useAuth();

  if (loadingUser) {
    return <div>loading...</div>;
  }

  return (
    <div className="px-12 py-6">
      {user && (
        <div className="flex flex-col gap-4">
          <CreateTweetForm />
          Logged in as {user.name}{" "}
          <TwButton
            className="max-w-fit"
            variant="danger"
            onClick={() => logout()}
          >
            Logout
          </TwButton>
        </div>
      )}
      {!user && (
        <div>
          <h1>Home</h1>
          <div>
            You are not logged in. Please{" "}
            <Link
              href="/login"
              className="px-3 py-1 rounded bg-blue-500 text-white"
            >
              Login
            </Link>
          </div>
        </div>
      )}
    </div>
  );
};

import { Link } from "@tanstack/react-router";
import { useAuth } from "../../components/context/auth.context";

export const HomePage = () => {
  const { user, logout, loadingUser } = useAuth();

  if (loadingUser) {
    return <div>Loading...</div>;
  }

  return (
    <div className="px-12 py-6">
      {user && (
        <div className="flex flex-col gap-4">
          Logged in as {user.name}{" "}
          <button
            onClick={() => logout()}
            className="px-4 py-2 w-min bg-red-500 rounded-full"
          >
            Logout
          </button>
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

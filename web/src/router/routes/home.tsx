import { useAuth } from "../../components/context/auth.context";

export const HomePage = () => {
  const { user, logout } = useAuth();

  return (
    <div>
      {user && (
        <div className="flex flex-col gap-4">
          Logged in as {user.name}{" "}
          <button
            onClick={() => logout()}
            className="px-3 py-1 w-min bg-red-500 rounded"
          >
            Logout
          </button>
        </div>
      )}
      {!user && (
        <a
          href="http://localhost:1337/api/auth/google/login"
          className="px-3 py-1 rounded bg-blue-500 text-white"
        >
          Login
        </a>
      )}
    </div>
  );
};

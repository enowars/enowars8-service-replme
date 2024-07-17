import { navigate } from "@/actions/navigate";
import { useMutation } from "@tanstack/react-query";
import axios from "axios";

export type RegisterPayload = {
  username: string;
  password: string;
};

export type RegisterMutationOptions = {
  onError?: () => void;
};

export function useRegisterMutation(options: RegisterMutationOptions) {
  return useMutation({
    mutationFn: (credentials: RegisterPayload) =>
      axios.post(
        (process.env.NEXT_PUBLIC_API ?? "") + "/api/auth/register",
        credentials,
        {
          withCredentials: true,
        },
      ),
    onSuccess: () => {
      navigate("/login");
    },
    onError: () => {
      if (options.onError) options.onError();
    },
  });
}

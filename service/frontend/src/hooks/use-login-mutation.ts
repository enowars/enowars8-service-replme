"use client";

import { navigate } from "@/actions/navigate";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";

export type LoginPayload = {
  username: string;
  password: string;
};

export type LoginMutationOptions = {
  onError?: () => Promise<unknown> | undefined | void;
};

export function useLoginMutation(options: LoginMutationOptions) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (credentials: LoginPayload) =>
      axios.post(
        (process.env.NEXT_PUBLIC_API ?? "") + "/api/auth/login",
        credentials,
        {
          withCredentials: true,
        },
      ),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user"] });
      navigate("/");
    },
    onError: options.onError,
  });
}

"use client";

import { navigate } from "@/actions/navigate";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import axios from "axios";

export function useLogoutMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () =>
      axios.post(
        (process.env.NEXT_PUBLIC_API ?? "") + "/api/auth/logout",
        undefined,
        {
          withCredentials: true,
        },
      ),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["user"] });
      navigate("/");
    },
  });
}

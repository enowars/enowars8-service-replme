"use client";

import { GetUserResponse } from "@/lib/types";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";

export function useUserQuery() {
  return useQuery({
    queryKey: ["user"],
    queryFn: () =>
      axios.get<GetUserResponse>(
        (process.env.NEXT_PUBLIC_API ?? "") + "/api/auth/user",
        {
          withCredentials: true,
        },
      ),
    staleTime: Infinity,
  });
}

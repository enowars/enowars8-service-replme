"use client";

import { useQuery } from "@tanstack/react-query";
import { useRef } from "react";

export type DevenvStateOptions = {
  uuid: string;
};

export function useDevenvState(options: DevenvStateOptions) {
  const generationRef = useRef<string>("ok");
  const query = useQuery<"ok" | "dirty" | undefined>({
    queryKey: ["devenv", options.uuid, "state"],
    staleTime: Infinity,
    retry: false,
  });
  if (query.isSuccess && query.data !== undefined) {
    generationRef.current = query.data;
    return generationRef.current;
  } else {
    return generationRef.current;
  }
}

"use client";

import { useQuery } from "@tanstack/react-query";
import { useRef } from "react";

export type DevenvGenerationOptions = {
  uuid: string;
};

export function useDevenvGeneration(options: DevenvGenerationOptions) {
  const generationRef = useRef<number>(0);
  const query = useQuery<number>({
    queryKey: ["devenv", options.uuid, "generation"],
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

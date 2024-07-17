import { useMutation, useQueryClient } from "@tanstack/react-query";

export type DevenvGenerationMutationOptions = {
  uuid?: string;
};

export function useDevenvGenerationMutation(
  options: DevenvGenerationMutationOptions,
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (gen?: number | null) => {
      if (gen === null) {
        queryClient.setQueryData(["devenv", options.uuid, "generation"], null);
      } else if (gen !== undefined) {
        queryClient.setQueryData(["devenv", options.uuid, "generation"], gen);
      } else {
        let generation = queryClient.getQueryData<number>([
          "devenv",
          options.uuid,
          "generation",
        ]);
        generation = generation ? generation + 1 : 1;
        queryClient.setQueryData(
          ["devenv", options.uuid, "generation"],
          generation,
        );
      }
    },
  });
}

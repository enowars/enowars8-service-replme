"use client";

import { ExitIcon } from "@radix-ui/react-icons";
import { Button } from "./ui/button";
import { useLogoutMutation } from "@/hooks/use-logout-mutation";

export function LogoutButton() {
  const logoutMutation = useLogoutMutation();

  return (
    <Button
      variant="outline"
      size="icon"
      onClick={() => logoutMutation.mutate()}
    >
      <ExitIcon className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100" />
    </Button>
  );
}

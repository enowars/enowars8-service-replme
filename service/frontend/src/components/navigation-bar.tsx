"use client";

import Image from "next/image";
import Link from "next/link";
import { ModeToggle } from "./mode-toggle";
import ReplMenu from "./repl-menu";
import { LoginButton } from "./login-button";
import DevenvMenu from "./devenv-menu";
import { LogoutButton } from "./logout-button";
import { useUserQuery } from "@/hooks/use-user-query";
import { usePathname } from "next/navigation";
import { Button } from "./ui/button";
import { useDevenvGenerationMutation } from "@/hooks/use-devenv-generation-mutation";
import {
  CheckIcon,
  CircleIcon,
  PlayIcon,
  ReloadIcon,
} from "@radix-ui/react-icons";
import DevenvSettingsMenu from "./devenv-settings-menu";
import CreateDevenvFileMenu from "./create-devenv-file-menu";
import { useMutationState } from "@tanstack/react-query";
import { useDevenvState } from "@/hooks/use-devenv-state";
import { useTheme } from "next-themes";
import { useEffect, useState } from "react";

const Navbar = () => {
  const pathname = usePathname();
  const { resolvedTheme } = useTheme();
  const userQuery = useUserQuery();
  const isAuthenticatedMode = !userQuery.isStale && userQuery.isSuccess;

  const [theme, setTheme] = useState("light");

  const match = pathname.match(
    "(?<=/devenv/)[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}",
  );
  const devenvUuid =
    match !== null && match?.length >= 1 ? match[0] : undefined;

  const devenvGenerationMutation = useDevenvGenerationMutation({
    uuid: devenvUuid,
  });

  const devenvFileContentMutationsStatus = useMutationState({
    filters: { mutationKey: ["devenv", devenvUuid, "files"] },
    select: (mutation) => mutation.state.status,
  });

  const devenvState = useDevenvState({
    uuid: devenvUuid ?? "",
  });

  useEffect(() => {
    setTheme(resolvedTheme ?? "light");
  }, [resolvedTheme]);

  return (
    <nav className="w-full fixed px-10 xs:px-20 py-3 light:bg-white/50 backdrop-blur-lg z-30">
      <div className="flex flex-row justify-between items-center">
        <div className="flex flex-row items-center space-x-10">
          <Link
            href="/"
            className="flex flex-row items-center space-x-5 text-2xl font-bold"
          >
            <Image
              src={
                theme === "dark"
                  ? "/favico-alpha-white.png"
                  : "/favico-alpha-black.png"
              }
              width={25}
              height={25}
              alt=""
            />
            <div>replme</div>
          </Link>
          {devenvUuid && (
            <div className="flex flex-row items-center space-x-3">
              <Button
                className="bg-green-600 dark:bg-green-300 hover:bg-green-800 dark:hover:bg-green-400"
                onClick={() => devenvGenerationMutation.mutate(undefined)}
              >
                <PlayIcon className="mr-2 h-4 w-4" /> Run
              </Button>
              <CreateDevenvFileMenu uuid={devenvUuid} />
              <DevenvSettingsMenu uuid={devenvUuid} />
              {devenvState === "dirty" ? (
                devenvFileContentMutationsStatus.includes("pending") ? (
                  <ReloadIcon className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <CircleIcon className="mr-2 h-4 w-4 text-yellow-500" />
                )
              ) : (
                <div className="flex flex-row items-center space-x-2 text-green-500">
                  <CheckIcon />
                  <div>Saved</div>
                </div>
              )}
            </div>
          )}
        </div>
        <div className="flex flex-row items-center space-x-3">
          {isAuthenticatedMode && <DevenvMenu />}
          <ReplMenu />
          <ModeToggle />
          {isAuthenticatedMode ? <LogoutButton /> : <LoginButton />}
        </div>
      </div>
    </nav>
  );
};

export default Navbar;

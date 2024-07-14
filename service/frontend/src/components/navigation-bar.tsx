"use client";

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
import { PlayIcon } from '@radix-ui/react-icons'
import DevenvSettingsMenu from "./devenv-settings-menu";
import CreateDevenvFileMenu from "./create-devenv-file-menu";

const Navbar = () => {
  const pathname = usePathname();
  const userQuery = useUserQuery()
  const isAuthenticatedMode = !userQuery.isStale && userQuery.isSuccess

  const match = pathname.match("(?<=/devenv/)[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}")
  const devenvUuid = match !== null && match?.length >= 1 ? match[0] : undefined;

  const devenvGenerationMutation = useDevenvGenerationMutation({
    uuid: devenvUuid
  })

  return (
    <nav className="w-full fixed px-10 xs:px-20 py-3 light:bg-white/50 backdrop-blur-lg z-30">
      <div className="flex flex-row justify-between items-center">
        <div className="flex flex-row items-center space-x-10">
          <Link
            href="/"
            className="text-2xl font-bold"
          >
            replme
          </Link>
          {devenvUuid && <div className="flex flex-row items-center space-x-3">
            <Button className="bg-green-600 hover:bg-green-800" onClick={() => devenvGenerationMutation.mutate(undefined)}>
              <PlayIcon className="mr-2 h-4 w-4" /> Run
            </Button>
            <CreateDevenvFileMenu uuid={devenvUuid} />
            <DevenvSettingsMenu uuid={devenvUuid} />
          </div>}
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
}

export default Navbar;


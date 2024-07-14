"use client";

import Link from "next/link";
import { ModeToggle } from "./mode-toggle";
import ReplMenu from "./repl-menu";
import { LoginButton } from "./login-button";
import DevenvMenu from "./devenv-menu";
import { LogoutButton } from "./logout-button";
import { useUserQuery } from "@/hooks/use-user-query";

const Navbar = () => {
  const userQuery = useUserQuery()
  const isAuthenticatedMode = !userQuery.isStale && userQuery.isSuccess

  return (
    <nav className="w-full fixed px-10 xs:px-20 py-3 light:bg-white/50 backdrop-blur-lg z-30">
      <div className="flex flex-row justify-between items-center">
        <Link
          href="/"
          className="text-2xl font-bold"
        >
          replme
        </Link>
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


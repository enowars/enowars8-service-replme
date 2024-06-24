"use client";

import Link from "next/link";
import { ModeToggle } from "./mode-toggle";
import ReplMenu from "./repl-menu";
import { LoginButton } from "./login-button";
import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { GetUserResponse } from "@/lib/types";
import { useEffect } from "react";
import DevenvMenu from "./devenv-menu";
import { LogoutButton } from "./logout-button";

const Navbar = () => {
  const query = useQuery({
    queryKey: ['user'],
    queryFn: () => axios.get<GetUserResponse>(
      process.env.NEXT_PUBLIC_API + "/api/auth/user",
      {
        withCredentials: true
      }
    ),
    staleTime: Infinity
  })

  const isAuthenticatedMode = !query.isStale && query.isSuccess

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


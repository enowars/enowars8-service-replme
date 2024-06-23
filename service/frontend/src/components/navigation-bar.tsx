"use client";

import Link from "next/link";
import { ModeToggle } from "./mode-toggle";
import ReplMenu from "./repl-menu";

function Navbar() {
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
          <ReplMenu />
          <ModeToggle />
        </div>
      </div>
    </nav>
  );
}

export default Navbar;


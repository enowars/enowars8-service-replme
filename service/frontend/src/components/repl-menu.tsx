"use client";

import { Button } from "@/components/ui/button"
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer"
import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { CodeIcon } from "@radix-ui/react-icons";
import CreateReplButton from "./create-repl-button";
import { navigate } from "@/actions/navigate";

const ReplMenu = () => {
  const query = useQuery(
    {
      queryKey: ["repl-sessions"],
      queryFn: () => axios.get<string[]>(process.env.NEXT_PUBLIC_API + '/api/user/sessions', { withCredentials: true }),
    }
  )

  const numSessions = query.data?.data?.length ?? 0;

  return (
    <Drawer>
      <DrawerTrigger asChild>
        <Button className="relative" variant="outline" size="icon">
          <CodeIcon className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all" />
          {Boolean(numSessions) && <div className="absolute w-5 h-5 bg-red-400 -top-2 -right-2 rounded-full text-white">{numSessions}</div>}
        </Button>
      </DrawerTrigger>
      <DrawerContent className="w-full">
        <div className="w-full flex flex-col items-center pb-5">
          <DrawerHeader>
            <DrawerTitle>Your REPLs</DrawerTitle>
            <DrawerDescription>Open a repl by clicking it</DrawerDescription>
          </DrawerHeader>

          <div className="flex flex-row justify-center items-center w-full space-x-5 overflow-auto">
            {query.data?.data?.map((id) => (
              <DrawerClose asChild key={id}>
                <Button variant="outline" onClick={() => navigate("/repl/" + id)}>
                  {id.substring(0, 5)}
                </Button>
              </DrawerClose>
            ))}
            <DrawerClose asChild>
              <CreateReplButton id="close-button" />
            </DrawerClose>
          </div>
        </div>
      </DrawerContent>
    </Drawer>
  )
}

export default ReplMenu;

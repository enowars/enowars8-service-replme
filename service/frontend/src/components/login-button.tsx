import { EnterIcon } from "@radix-ui/react-icons";
import { Button } from "./ui/button";
import { navigate } from "@/actions/navigate";

export function LoginButton() {
  return (
    <Button variant="outline" size="icon" onClick={() => navigate("/login")}>
      <EnterIcon className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100" />
    </Button>
  )
}

